<?php

declare(strict_types=1);

namespace Spiral\Goridge\RPC;

use Spiral\Goridge\Frame;
use Spiral\Goridge\RelayInterface;
use Spiral\Goridge\RPC\Exception\ServiceException;

abstract class AbstractRPC implements RPCInterface
{
    /**
     * RPC calls service prefix.
     *
     * @var non-empty-string|null
     */
    protected ?string $service = null;

    /**
     * @var positive-int
     */
    protected static int $seq = 1;

    public function __construct(
        protected CodecInterface $codec,
    ) {}

    /**
     * @psalm-pure
     */
    public function withServicePrefix(string $service): self
    {
        /** @psalm-suppress ImpureVariable */
        $rpc = clone $this;
        $rpc->service = $service;

        return $rpc;
    }

    /**
     * @psalm-pure
     */
    public function withCodec(CodecInterface $codec): self
    {
        /** @psalm-suppress ImpureVariable */
        $rpc = clone $this;
        $rpc->codec = $codec;

        return $rpc;
    }

    /**
     * @throws Exception\ServiceException
     */
    protected function decodeResponse(Frame $frame, RelayInterface $relay, mixed $options = null): mixed
    {
        // exclude method name
        $body = \substr((string) $frame->payload, $frame->options[1]);

        if ($frame->hasFlag(Frame::ERROR)) {
            $name = $relay instanceof \Stringable
                ? (string) $relay
                : $relay::class;

            throw new ServiceException(\sprintf("Error '%s' on %s", $body, $name));
        }

        return $this->codec->decode($body, $options);
    }

    /**
     * @param non-empty-string $method
     */
    protected function packFrame(string $method, mixed $payload): Frame
    {
        if ($this->service !== null) {
            $method = $this->service . '.' . \ucfirst($method);
        }

        $body = $method . $this->codec->encode($payload);
        return new Frame($body, [self::$seq, \strlen($method)], $this->codec->getIndex());
    }
}
