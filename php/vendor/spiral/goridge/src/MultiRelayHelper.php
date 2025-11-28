<?php

declare(strict_types=1);

namespace Spiral\Goridge;

use Spiral\Goridge\RPC\Exception\RPCException;

class MultiRelayHelper
{
    /**
     * @param array<array-key, RelayInterface> $relays
     * @return array-key[]|false
     * @internal
     * Returns either
     *  - an array of array keys, even if only one
     *  - or false if none
     */
    public static function findRelayWithMessage(array $relays, int $timeoutInMicroseconds = 0): array|false
    {
        if (\count($relays) === 0) {
            return false;
        }

        if ($relays[\array_key_first($relays)] instanceof SocketRelay) {
            $sockets = [];
            $socketIdToRelayIndexMap = [];
            foreach ($relays as $relayIndex => $relay) {
                \assert($relay instanceof SocketRelay);

                // Enforce connection
                if ($relay->socket === null) {
                    // Important: Do not force reconnect here as it would otherwise completely ruin further handling
                    continue;
                }

                $sockets[] = $relay->socket;
                $socketIdToRelayIndexMap[\spl_object_id($relay->socket)] = $relayIndex;
            }

            if (\count($sockets) === 0) {
                return false;
            }

            $writes = null;
            $except = null;
            $changes = \socket_select($sockets, $writes, $except, 0, $timeoutInMicroseconds);

            if ($changes > 0) {
                $indexes = [];
                foreach ($sockets as $socket) {
                    $indexes[] = $socketIdToRelayIndexMap[\spl_object_id($socket)] ?? throw new RPCException("Invalid socket??");
                }

                return $indexes;
            }
            return false;

        }

        if ($relays[\array_key_first($relays)] instanceof StreamRelay) {
            $streams = [];
            $streamNameToRelayIndexMap = [];
            foreach ($relays as $relayIndex => $relay) {
                \assert($relay instanceof StreamRelay);

                $streams[] = $relay->in;
                $streamNameToRelayIndexMap[(string) $relay->in] = $relayIndex;
            }

            $writes = null;
            $except = null;
            $changes = \stream_select($streams, $writes, $except, 0, $timeoutInMicroseconds);

            if ($changes > 0) {
                $indexes = [];
                foreach ($streams as $stream) {
                    $indexes[] = $streamNameToRelayIndexMap[(string) $stream] ?? throw new RPCException("Invalid stream??");
                }

                return $indexes;
            }
            return false;

        }

        return false;
    }

    /**
     * @param array<array-key, RelayInterface> $relays
     * @return array-key[]|false
     * @internal
     * Returns either
     *  - an array of array keys, even if only one
     *  - or false if none
     */
    public static function checkConnected(array $relays): array|false
    {
        if (\count($relays) === 0) {
            return false;
        }

        $keysNotConnected = [];
        foreach ($relays as $key => $relay) {
            if ($relay instanceof ConnectedRelayInterface && !$relay->isConnected()) {
                $relay->connect();
                $keysNotConnected[] = $key;
            }
        }

        return $keysNotConnected;
    }
}
