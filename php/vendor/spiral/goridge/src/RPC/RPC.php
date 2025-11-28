<?php

declare(strict_types=1);

namespace Spiral\Goridge\RPC;

use Spiral\Goridge\Relay;
use Spiral\Goridge\RelayInterface;
use Spiral\Goridge\RPC\Codec\JsonCodec;
use Spiral\Goridge\RPC\Exception\RPCException;

class RPC extends AbstractRPC
{
    public function __construct(
        private readonly RelayInterface $relay,
        CodecInterface          $codec = new JsonCodec(),
    ) {
        parent::__construct($codec);
    }

    /**
     * @param non-empty-string $connection
     */
    public static function create(string $connection, CodecInterface $codec = new JsonCodec()): RPCInterface
    {
        $relay = Relay::create($connection);

        return new self($relay, $codec);
    }

    public function call(string $method, mixed $payload, mixed $options = null): mixed
    {
        $this->relay->send($this->packFrame($method, $payload));

        // wait for the frame confirmation
        $frame = $this->relay->waitFrame();

        if (\count($frame->options) !== 2) {
            throw new RPCException('Invalid RPC frame, options missing');
        }

        if ($frame->options[0] !== self::$seq) {
            throw new RPCException('Invalid RPC frame, sequence mismatch');
        }

        self::$seq++;

        return $this->decodeResponse($frame, $this->relay, $options);
    }
}
