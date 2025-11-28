<?php

declare(strict_types=1);

namespace Spiral\Goridge\RPC\Codec;

use Spiral\Goridge\Frame;
use Spiral\Goridge\RPC\CodecInterface;
use Spiral\Goridge\RPC\Exception\CodecException;

final class JsonCodec implements CodecInterface
{
    public function getIndex(): int
    {
        return Frame::CODEC_JSON;
    }

    public function encode(mixed $payload): string
    {
        try {
            $result = \json_encode($payload, \JSON_THROW_ON_ERROR);
        } catch (\JsonException $e) {
            throw new CodecException(\sprintf('Json encode: %s', $e->getMessage()), $e->getCode(), $e);
        }

        return $result;
    }

    public function decode(string $payload, mixed $options = null): mixed
    {
        try {
            return \json_decode($payload, true, 512, \JSON_THROW_ON_ERROR | (\is_int($options) ? $options : 0));
        } catch (\JsonException $e) {
            throw new CodecException(\sprintf('Json decode: %s', $e->getMessage()), $e->getCode(), $e);
        }
    }
}
