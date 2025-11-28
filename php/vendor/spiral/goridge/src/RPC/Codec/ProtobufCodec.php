<?php

declare(strict_types=1);

namespace Spiral\Goridge\RPC\Codec;

use Spiral\Goridge\Frame;
use Google\Protobuf\Internal\Message;
use Spiral\Goridge\RPC\CodecInterface;

final class ProtobufCodec implements CodecInterface
{
    private const ERROR_DEPENDENCY =
        'Could not initialize protobuf codec. ' .
        'Please add "ext-protobuf" PECL extension or ' .
        'install "google/protobuf" Composer dependency.';

    public function __construct()
    {
        $this->assertAvailable();
    }

    public function getIndex(): int
    {
        return Frame::CODEC_PROTO;
    }

    /**
     * @psalm-suppress MixedInferredReturnType
     * @psalm-suppress MixedReturnStatement
     */
    public function encode(mixed $payload): string
    {
        if ($payload instanceof Message) {
            return $payload->serializeToString();
        }

        return $payload;
    }

    public function decode(string $payload, mixed $options = null): mixed
    {
        if (\is_string($options) && \is_subclass_of($options, Message::class, true)) {
            $options = $this->create($options);
        }

        if ($options instanceof Message) {
            $options->mergeFromString($payload);

            return $options;
        }

        return $payload;
    }

    /**
     * @psalm-suppress UnsafeInstantiation
     *
     * @param class-string<Message> $class
     */
    protected function create(string $class): Message
    {
        return new $class();
    }

    private function assertAvailable(): void
    {
        if (!\class_exists(Message::class)) {
            throw new \LogicException(self::ERROR_DEPENDENCY);
        }
    }
}
