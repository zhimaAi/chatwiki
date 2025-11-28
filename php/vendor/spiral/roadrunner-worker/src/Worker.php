<?php

declare(strict_types=1);

namespace Spiral\RoadRunner;

use Psr\Log\LoggerInterface;
use Spiral\Goridge\BlockingRelayInterface;
use Spiral\Goridge\Exception\GoridgeException;
use Spiral\Goridge\Exception\TransportException;
use Spiral\Goridge\Frame;
use Spiral\Goridge\Relay;
use Spiral\Goridge\RelayInterface;
use Spiral\RoadRunner\Exception\RoadRunnerException;
use Spiral\RoadRunner\Internal\StdoutHandler;
use Spiral\RoadRunner\Message\Command\GetProcessId;
use Spiral\RoadRunner\Message\Command\Pong;
use Spiral\RoadRunner\Message\Command\WorkerStop;
use Spiral\RoadRunner\Message\SkipMessage;

/**
 * Accepts connection from RoadRunner server over given Goridge relay.
 *
 * <code>
 * $worker = Worker::create();
 *
 * while ($receivedPayload = $worker->waitPayload()) {
 *      $worker->respond(new Payload("DONE", json_encode($context)));
 * }
 * </code>
 */
class Worker implements StreamWorkerInterface
{
    private const JSON_ENCODE_FLAGS = \JSON_THROW_ON_ERROR | \JSON_PRESERVE_ZERO_FRACTION;

    /** @var array<int, Payload> */
    private array $payloads = [];

    /** @var int<0, max> Count of frames sent in stream mode */
    private int $framesSent = 0;

    private bool $streamMode = false;
    private bool $shouldPing = false;
    private bool $waitingPong = false;

    public function __construct(
        private readonly RelayInterface $relay,
        bool $interceptSideEffects = true,
        private readonly LoggerInterface $logger = new Logger(),
    ) {
        if ($interceptSideEffects) {
            StdoutHandler::register();
        }
    }

    /**
     * Create a new RoadRunner {@see Worker} using global
     * environment ({@see Environment}) configuration.
     */
    public static function create(bool $interceptSideEffects = true, LoggerInterface $logger = new Logger()): self
    {
        return static::createFromEnvironment(
            env: Environment::fromGlobals(),
            interceptSideEffects: $interceptSideEffects,
            logger: $logger,
        );
    }

    /**
     * Create a new RoadRunner {@see Worker} using passed environment
     * configuration.
     */
    public static function createFromEnvironment(
        EnvironmentInterface $env,
        bool $interceptSideEffects = true,
        LoggerInterface $logger = new Logger(),
    ): self {
        $address = $env->getRelayAddress();
        \assert($address !== '', 'Relay address must be specified in environment');

        return new self(
            relay: Relay::create($address),
            interceptSideEffects: $interceptSideEffects,
            logger: $logger,
        );
    }

    public function getLogger(): LoggerInterface
    {
        return $this->logger;
    }

    public function waitPayload(): ?Payload
    {
        while (true) {
            if ($this->payloads !== []) {
                $payload = \array_shift($this->payloads);
            } else {
                $frame = $this->relay->waitFrame();
                $payload = PayloadFactory::fromFrame($frame);
            }

            switch (true) {
                case $payload::class === Payload::class:
                    return $payload;
                case $payload instanceof WorkerStop:
                    $this->waitingPong = false;
                    return null;
                case $payload::class === GetProcessId::class:
                    $this->sendProcessId();
                    continue 2;
                case $payload instanceof Pong:
                    $this->waitingPong = false;
                    continue 2;
                case $payload instanceof SkipMessage:
                    continue 2;
            }
        }
    }

    public function withStreamMode(): static
    {
        $clone = clone $this;
        $clone->streamMode = true;
        $clone->framesSent = 0;
        $clone->shouldPing = false;
        $clone->waitingPong = false;
        return $clone;
    }

    /**
     * @param int|null $codec The codec used for encoding the payload header.
     *        Can be {@see Frame::CODEC_PROTO} for Protocol Buffers or {@see Frame::CODEC_JSON} for JSON.
     *        This parameter will be removed in v4.0 and {@see Frame::CODEC_PROTO} will be used by default.
     */
    public function respond(Payload $payload, ?int $codec = null): void
    {
        $this->streamMode and ++$this->framesSent;
        $this->send($payload->body, $payload->header, $payload->eos, $codec);
    }

    public function error(string $error): void
    {
        $frame = new Frame($error, [], Frame::ERROR);

        $this->sendFrame($frame);
    }

    public function stop(): void
    {
        $this->send('', $this->encode(['stop' => true]));
    }

    public function hasPayload(?string $class = null): bool
    {
        return $this->findPayload($class) !== null;
    }

    public function getPayload(?string $class = null): ?Payload
    {
        $pos = $this->findPayload($class);
        if ($pos === null) {
            return null;
        }
        $result = $this->payloads[$pos];
        unset($this->payloads[$pos]);

        return $result;
    }

    /**
     * @param class-string<Payload>|null $class
     *
     * @return null|int Index in {@see $this->payloads} or null if not found
     */
    private function findPayload(?string $class = null): ?int
    {
        // Find in existing payloads
        if ($this->payloads !== []) {
            if ($class === null) {
                return \array_key_first($this->payloads);
            }

            foreach ($this->payloads as $pos => $payload) {
                if ($payload::class === $class) {
                    return $pos;
                }
            }
        }

        do {
            if ($class === null && $this->payloads !== []) {
                return \array_key_first($this->payloads);
            }

            $payload = $this->pullPayload();
            if ($payload === null || $payload instanceof Pong) {
                break;
            }

            $this->payloads[] = $payload;
            if ($class !== null && $payload::class === $class) {
                return \array_key_last($this->payloads);
            }
        } while (true);

        return null;
    }

    /**
     * Pull {@see Payload} if it is available without blocking.
     */
    private function pullPayload(): ?Payload
    {
        if (!$this->waitingPong && $this->relay instanceof BlockingRelayInterface) {
            if (!$this->streamMode) {
                return null;
            }

            $this->haveToPing();
            return null;
        }

        if (!$this->relay->hasFrame()) {
            return null;
        }

        $frame = $this->relay->waitFrame();
        $payload = PayloadFactory::fromFrame($frame);

        if ($payload instanceof Pong) {
            $this->waitingPong = false;
            return null;
        }

        return $payload;
    }

    private function send(string $body = '', string $header = '', bool $eos = true, ?int $codec = null): void
    {
        $frame = new Frame($header . $body, [\strlen($header)]);

        if (!$eos) {
            $frame->byte10 |= Frame::BYTE10_STREAM;
        }

        if ($this->shouldPing) {
            $frame->byte10 |= Frame::BYTE10_PING;
        }

        if ($codec !== null) {
            $frame->setFlag($codec);
        }

        $this->sendFrame($frame);
    }

    private function sendFrame(Frame $frame): void
    {
        try {
            if ($this->streamMode && ($frame->byte10 & Frame::BYTE10_STREAM) && $this->shouldPing) {
                $frame->byte10 |= Frame::BYTE10_PING;
                $this->shouldPing = false;
                $this->waitingPong = true;
            }

            $this->relay->send($frame);
        } catch (GoridgeException $e) {
            throw new TransportException($e->getMessage(), $e->getCode(), $e);
        } catch (\Throwable $e) {
            throw new RoadRunnerException($e->getMessage(), (int) $e->getCode(), $e);
        }
    }

    private function encode(array $payload): string
    {
        return \json_encode($payload, self::JSON_ENCODE_FLAGS);
    }

    private function sendProcessId(): void
    {
        $frame = new Frame($this->encode(['pid' => \getmypid()]), [], Frame::CONTROL);
        $this->sendFrame($frame);
    }

    private function haveToPing(): void
    {
        if ($this->waitingPong || $this->framesSent === 0) {
            return;
        }

        if ($this->framesSent % 5 === 0) {
            $this->shouldPing = true;
        }
    }
}
