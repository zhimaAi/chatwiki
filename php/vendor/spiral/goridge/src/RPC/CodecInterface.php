<?php

declare(strict_types=1);

namespace Spiral\Goridge\RPC;

use Spiral\Goridge\RPC\Exception\CodecException;

/**
 * Serializes incoming and deserializes received messages.
 */
interface CodecInterface
{
    /**
     * Coded index, uniquely identified by remote server.
     */
    public function getIndex(): int;

    /**
     * @throws CodecException
     */
    public function encode(mixed $payload): string;

    /**
     * @throws CodecException
     */
    public function decode(string $payload, mixed $options = null): mixed;
}
