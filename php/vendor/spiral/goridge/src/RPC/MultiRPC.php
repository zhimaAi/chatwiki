<?php

declare(strict_types=1);

namespace Spiral\Goridge\RPC;

use Spiral\Goridge\ConnectedRelayInterface;
use Spiral\Goridge\Exception\RelayException;
use Spiral\Goridge\Exception\TransportException;
use Spiral\Goridge\Frame;
use Spiral\Goridge\MultiRelayHelper;
use Spiral\Goridge\Relay;
use Spiral\Goridge\RelayInterface;
use Spiral\Goridge\RPC\Codec\JsonCodec;
use Spiral\Goridge\RPC\Exception\RPCException;
use Spiral\Goridge\SocketRelay;

class MultiRPC extends AbstractRPC implements AsyncRPCInterface
{
    /**
     * The default is at 10_000 because in tests things like Doctrine hammer this very hard when used in caching.
     * A limit of 1_000 was hit repeatedly. Make it configurable though in case someone wants to change it.
     */
    private const DEFAULT_BUFFER_THRESHOLD = 10_000;

    public const ERR_INVALID_SEQ_NUMBER = 'Invalid sequence number. This may occur if the number was already used, the buffers were flushed due to insufficient getResponse calling, or with a plain incorrect number. Please check your code.';

    /**
     * @var array<int, RelayInterface>
     */
    private static array $freeRelays = [];

    /**
     * Occupied Relays is a map of seq to relay to make removal easier once a response is received.
     * @var array<positive-int, RelayInterface>
     */
    private static array $occupiedRelays = [];

    /**
     * A map of seq to relay to use for decodeResponse().
     * Technically the relay there is only needed in case of an error.
     *
     * @var array<positive-int, RelayInterface>
     */
    private static array $seqToRelayMap = [];

    /**
     * Map of seq to response Frame
     * Should only really need to be used in cases of high amounts of traffic
     *
     * @var array<positive-int, Frame>
     */
    private static array $asyncResponseBuffer = [];

    /**
     * The threshold after which the asyncResponseBuffer is flushed of all entries.
     */
    private int $asyncBufferThreshold = self::DEFAULT_BUFFER_THRESHOLD;

    /**
     * @param array<int, RelayInterface> $relays
     */
    public function __construct(
        array $relays,
        int $asyncBufferThreshold = self::DEFAULT_BUFFER_THRESHOLD,
        CodecInterface $codec = new JsonCodec(),
    ) {
        // Check if we have at least one either existing or new relay here
        if (\count($relays) === 0 && \count(self::$freeRelays) === 0 && \count(self::$occupiedRelays) === 0) {
            throw new RPCException("MultiRPC needs at least one relay. Zero provided.");
        }

        if (\count($relays) > 0) {
            // Check if all new relays are of the same type
            if (\count(\array_unique(\array_map(static fn(RelayInterface $relay) => $relay::class, $relays))) > 1) {
                throw new RPCException("MultiRPC can only be used with all relays of the same type, such as a " . SocketRelay::class);
            }

            // Check if the existing relays (if any) and the new relays are of the same type.
            if (\count(self::$freeRelays) > 0) {
                $existingRelay = self::$freeRelays[0];
            } elseif (\count(self::$occupiedRelays) > 0) {
                $existingRelay = self::$occupiedRelays[\array_key_first(self::$occupiedRelays)];
            } else {
                $existingRelay = null;
            }

            if ($existingRelay !== null && $existingRelay::class !== $relays[0]::class) {
                throw new RPCException("MultiRPC can only be used with all relays of the same type, such as a " . SocketRelay::class);
            }
        }

        // The relays (and related arrays) are static to support cloning this class.
        // Basically the following problem exists:
        // - If we make these arrays instance variables, then we need to recreate the relays on clone, otherwise we'd run into data issues.
        // When we do that, the number of relays in existence can increase quite dramatically, resulting in balooning memory usage for socket buffers.
        // - If we make these arrays static variables, then we need to make certain that they stay the same across all instances
        // of this class. As a result the arrays are basically only appended on, and never deleted or modified.
        // In the end that *can* mean that if someone were to repeatedly call `new MultiRPC([a bunch of relays])` that we'd
        // tack all those relays into this array resulting in the same problem.
        // It also means that different services can cannibalize the number of relays available to them,
        // for example a Metrics service and a KV (Cache) service.
        // IMHO(L3tum) a balooning memory usage that occurs unexpectly is way worse, than any of the other problems. In the end
        // one can work against cannibalized relays by simply upping the number of relays at any point.
        self::$freeRelays = [...self::$freeRelays, ...$relays];
        $this->asyncBufferThreshold = $asyncBufferThreshold;
        parent::__construct($codec);
    }

    /**
     * @param non-empty-string $connection
     * @param positive-int $count
     */
    public static function create(
        string $connection,
        int $count = 50,
        int $asyncBufferThreshold = self::DEFAULT_BUFFER_THRESHOLD,
        CodecInterface $codec = new JsonCodec(),
    ): self {
        \assert($count > 0);
        $count = $count - \count(self::$freeRelays) - \count(self::$occupiedRelays);
        $relays = [];

        for ($i = 0; $i < $count; $i++) {
            $relay = Relay::create($connection);
            $relays[] = $relay;
        }

        return new self($relays, $asyncBufferThreshold, $codec);
    }

    /**
     * Force-connects all relays.
     * @throws RelayException
     */
    public function preConnectRelays(): void
    {
        foreach (self::$freeRelays as $relay) {
            if ($relay instanceof ConnectedRelayInterface) {
                // Force connect
                $relay->connect();
            }
        }
    }

    public function call(string $method, mixed $payload, mixed $options = null): mixed
    {
        $relayIndex = $this->ensureFreeRelayAvailable();
        $relay = self::$freeRelays[$relayIndex];

        $relay->send($this->packFrame($method, $payload));

        // wait for the frame confirmation
        $frame = $this->getResponseFromRelay($relay, self::$seq, true);

        self::$seq++;

        return $this->decodeResponse($frame, $relay, $options);
    }

    public function callIgnoreResponse(string $method, mixed $payload): void
    {
        $relayIndex = $this->ensureFreeRelayAvailable();
        $relay = self::$freeRelays[$relayIndex];

        $relay->send($this->packFrame($method, $payload));

        $seq = self::$seq;
        self::$seq++;
        self::$occupiedRelays[$seq] = $relay;
        // Last index so no need for array_pop or stuff
        unset(self::$freeRelays[$relayIndex]);
    }

    public function callAsync(string $method, mixed $payload): int
    {
        // Flush buffer if someone doesn't call getResponse
        if (\count(self::$asyncResponseBuffer) > $this->asyncBufferThreshold) {
            // We don't need to clean up occupiedRelays here since the buffer is solely for responses already
            // fetched from relays, and those relays are put back to freeRelays in getNextFreeRelay()
            self::$seqToRelayMap = \array_diff_key(self::$seqToRelayMap, self::$asyncResponseBuffer);
            self::$asyncResponseBuffer = [];
        }

        $relayIndex = $this->ensureFreeRelayAvailable();
        $relay = self::$freeRelays[$relayIndex];

        $relay->send($this->packFrame($method, $payload));

        $seq = self::$seq;
        self::$seq++;
        self::$occupiedRelays[$seq] = $relay;
        self::$seqToRelayMap[$seq] = $relay;
        // Last index so no need for array_pop or stuff
        unset(self::$freeRelays[$relayIndex]);

        return $seq;
    }

    public function hasResponse(int $seq): bool
    {
        // Check if we have the response buffered previously due to congestion
        if (isset(self::$asyncResponseBuffer[$seq])) {
            return true;
        }

        // Else check if the relay has the response in its buffer
        if (self::$seqToRelayMap[$seq]->hasFrame()) {
            return true;
        }

        return false;
    }

    public function hasResponses(array $seqs): array
    {
        $relays = [];
        /** @var array<int, positive-int> $relayIndexToSeq */
        $relayIndexToSeq = [];
        $seqsWithResponse = [];

        // The behaviour is essentially the same as self::hasResponse, just mapped to multiple $seqs aka $relays.
        // In order to use MultiRelayHelper we create a map of index => seq to map it back after checking for messages.
        foreach ($seqs as $seq) {
            if (isset(self::$asyncResponseBuffer[$seq])) {
                $seqsWithResponse[] = $seq;
            } elseif (isset(self::$seqToRelayMap[$seq])) {
                $relayIndexToSeq[\count($relays)] = $seq;
                $relays[] = self::$seqToRelayMap[$seq];
            }
        }

        /** @var int[]|false $index */
        $index = MultiRelayHelper::findRelayWithMessage($relays);

        if ($index === false) {
            return $seqsWithResponse;
        }

        foreach ($index as $relayIndex) {
            $seqsWithResponse[] = $relayIndexToSeq[$relayIndex];
        }

        return $seqsWithResponse;
    }

    public function getResponse(int $seq, mixed $options = null): mixed
    {
        $relay = self::$seqToRelayMap[$seq] ?? throw new RPCException(self::ERR_INVALID_SEQ_NUMBER);
        unset(self::$seqToRelayMap[$seq]);

        if (($frame = $this->getResponseFromBuffer($seq)) !== null) {
            /**
             * We can assume through @see MultiRPC::ensureFreeRelayAvailable() that a relay whose response is already
             * in this buffer has also been added to freeRelays (or is otherwise occupied).
             * Thus we only re-add (and do so without searching for it first) if we don't have the response yet.
             */
        } else {
            self::$freeRelays[] = self::$occupiedRelays[$seq];
            unset(self::$occupiedRelays[$seq]);

            $frame = $this->getResponseFromRelay($relay, $seq, true);
        }

        return $this->decodeResponse($frame, $relay, $options);
    }

    public function getResponses(array $seqs, mixed $options = null): iterable
    {
        // Quick return
        if (\count($seqs) === 0) {
            return;
        }

        // Flip the array to use the $seqs for key indexing
        $seqsKeyed = [];

        foreach ($seqs as $seq) {
            if (isset(self::$asyncResponseBuffer[$seq])) {
                // We can use getResponse() here since it's doing basically what we want to do here anyway
                yield $seq => $this->getResponse($seq, $options);
            } else {
                $seqsKeyed[$seq] = true;
            }
        }

        // Fetch all relays that are still occupied and which we need responses from
        $seqsToRelays = \array_intersect_key(self::$occupiedRelays, $seqsKeyed);

        // Make sure we have relays for all $seqs, otherwise something went wrong
        if (\count($seqsToRelays) !== \count($seqsKeyed)) {
            throw new RPCException(self::ERR_INVALID_SEQ_NUMBER);
        }

        $timeoutInMicroseconds = 0;
        while (\count($seqsToRelays) > 0) {
            // Do a first pass without a timeout. Maybe there's already most responses which would make a timeout unnecessary.
            /** @var positive-int[]|false $seqsReceivedResponse */
            $seqsReceivedResponse = MultiRelayHelper::findRelayWithMessage($seqsToRelays, $timeoutInMicroseconds);
            $timeoutInMicroseconds = 500;

            if ($seqsReceivedResponse === false) {
                if ($this->checkAllOccupiedRelaysStillConnected()) {
                    // Check if we've lost a relay we were waiting on, if so we need to quit since something is wrong.
                    if (\count(\array_diff_key($seqsToRelays, self::$occupiedRelays)) > 0) {
                        throw new RPCException(self::ERR_INVALID_SEQ_NUMBER);
                    }
                }
                continue;
            }

            foreach ($seqsReceivedResponse as $seq) {
                // Add the previously occupied relay to freeRelays here so that we don't lose it in case of an error
                $relay = $seqsToRelays[$seq];
                self::$freeRelays[] = $relay;
                unset(self::$occupiedRelays[$seq]);

                // Yield the response
                $frame = $this->getResponseFromRelay($relay, $seq, true);
                yield $seq => $this->decodeResponse($frame, $relay, $options);

                // Unset tracking map
                unset($seqsToRelays[$seq], self::$seqToRelayMap[$seq]);
            }
        }
    }

    /**
     * Returns array-key of free relay
     * @throws RPCException
     */
    private function ensureFreeRelayAvailable(): int
    {
        if (\count(self::$freeRelays) > 0) {
            // Return the last entry on self::$freeRelays so that further code can use unset() instead of array_splice (index handling)
            /** @psalm-return int */
            return \array_key_last(self::$freeRelays);
        }

        if (\count(self::$occupiedRelays) === 0) {
            // If we have neither freeRelays nor occupiedRelays then someone either initialized this with 0 relays
            // or something went terribly wrong. Either way we need to quit.
            throw new RPCException("No relays available at all");
        }

        while (\count(self::$freeRelays) === 0) {
            /** @var positive-int[]|false $index */
            $index = MultiRelayHelper::findRelayWithMessage(self::$occupiedRelays);

            if ($index === false) {
                // Check if all currently occupied relays are even still connected. Do another loop if they aren't.
                if ($this->checkAllOccupiedRelaysStillConnected()) {
                    continue;
                }
                // Just choose the first occupiedRelay to wait on since instead we may busyloop here
                // checking relay status and not giving RR the chance to actually answer (in a single core env for example).
                $index = [\array_key_first(self::$occupiedRelays)];

            }

            // Flush as many relays as we can up until a limit (arbitrarily 10?)
            for ($i = 0, $max = \min(10, \count($index)); $i < $max; $i++) {
                /** @var positive-int $seq */
                $seq = $index[$i];
                // Move relay from occupiedRelays into freeRelays before trying to get the response from it
                // in case something happens, so we don't lose it.
                $relay = self::$occupiedRelays[$seq];
                self::$freeRelays[] = $relay;
                unset(self::$occupiedRelays[$seq]);
                // Save response if in seqToRelayMap (aka a response is expected)
                // only save response in case of mismatched seq = response not in seqToRelayMap
                try {
                    $this->getResponseFromRelay($relay, $seq, !isset(self::$seqToRelayMap[$seq]));
                } catch (RelayException|RPCException) {
                    // Intentionally left blank
                }
            }
        }

        // Sometimes check if all occupied relays are even still connected
        $this->checkAllOccupiedRelaysStillConnected();

        // Return the last entry on self::$freeRelays so that further code can use unset() instead of array_splice (index handling)
        return \array_key_last(self::$freeRelays);
    }

    /**
     * Gets a response from the relay, blocking for it if necessary, with some error handling in regards to mismatched seq
     *
     * @param positive-int $expectedSeq
     */
    private function getResponseFromRelay(RelayInterface $relay, int $expectedSeq, bool $onlySaveResponseInCaseOfMismatchedSeq = false): Frame
    {
        if ($relay instanceof ConnectedRelayInterface && !$relay->isConnected()) {
            throw new TransportException("Unable to read payload from the stream");
        }

        $frame = $relay->waitFrame();

        if (\count($frame->options) !== 2) {
            // Expect at least a few options
            throw new RPCException('Invalid RPC frame, options missing');
        }

        if ($frame->options[0] !== $expectedSeq) {
            // Save response since $seq was invalid but the response may not
            /** @var positive-int $responseSeq */
            $responseSeq = $frame->options[0];
            self::$asyncResponseBuffer[$responseSeq] = $frame;

            throw new RPCException('Invalid RPC frame, sequence mismatch');
        }

        if (!$onlySaveResponseInCaseOfMismatchedSeq) {
            // If we want to save the response, regardless of whether the $seq was a match or not,
            // we'll need to add it to the buffer.
            // This is used in e.g. flushing a relay in ensureFreeRelay()
            // so that we can at least *try* to get the resonse back to the user.
            self::$asyncResponseBuffer[$expectedSeq] = $frame;
        }

        return $frame;
    }

    /**
     * Tries to get a response (Frame) from the buffer and unsets the entry if it finds the response.
     *
     * @param positive-int $seq
     */
    private function getResponseFromBuffer(int $seq): ?Frame
    {
        if (($frame = self::$asyncResponseBuffer[$seq] ?? null) !== null) {
            unset(self::$asyncResponseBuffer[$seq]);
        }

        return $frame;
    }

    private function checkAllOccupiedRelaysStillConnected(): bool
    {
        if (($relaysNotConnected = MultiRelayHelper::checkConnected(self::$occupiedRelays)) !== false) {
            /** @var positive-int $seq */
            foreach ($relaysNotConnected as $seq) {
                self::$freeRelays[] = self::$occupiedRelays[$seq];
                unset(self::$seqToRelayMap[$seq], self::$occupiedRelays[$seq]);
            }

            return true;
        }

        return false;
    }
}
