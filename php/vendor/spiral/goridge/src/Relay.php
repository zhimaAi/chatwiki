<?php

declare(strict_types=1);

namespace Spiral\Goridge;

use Spiral\Goridge\Exception\RelayFactoryException;

abstract class Relay implements RelayInterface
{
    public const TCP_SOCKET = 'tcp';
    public const UNIX_SOCKET = 'unix';
    public const PIPES = 'pipes';
    protected const CONNECTION_EXP = '/(?P<protocol>[^:\/]+):\/\/(?P<arg1>[^:]+)(:(?P<arg2>[^:]+))?/';

    /**
     * Create relay using string address.
     *
     * Example:
     *
     * Relay::create("pipes");
     * Relay::create("tcp://localhost:6001");
     *
     * @param non-empty-string $connection
     */
    public static function create(string $connection): RelayInterface
    {
        if ($connection === self::PIPES) {
            return new StreamRelay(STDIN, STDOUT);
        }

        if (!\preg_match(self::CONNECTION_EXP, $connection, $match)) {
            throw new Exception\RelayFactoryException('unsupported connection format');
        }

        /** @var array{protocol: non-empty-string, arg1: non-empty-string, arg2: non-empty-string} $match */
        $protocol = \strtolower($match['protocol']);

        switch ($protocol) {
            case self::TCP_SOCKET:
                //fall through
            case self::UNIX_SOCKET:
                $socketType = $protocol === self::TCP_SOCKET
                    ? SocketType::TCP
                    : SocketType::UNIX;

                $port = isset($match['arg2'])
                    ? (int) $match['arg2']
                    : null;

                /** @psalm-suppress ArgumentTypeCoercion Reason: Checked in the SocketRelay constructor */
                return new SocketRelay($match['arg1'], $port, $socketType);

            case self::PIPES:
                if (!isset($match['arg2'])) {
                    throw new RelayFactoryException('Unsupported stream connection format');
                }

                return new StreamRelay(self::openIn($match['arg1']), self::openOut($match['arg2']));

            default:
                throw new Exception\RelayFactoryException('unknown connection protocol');
        }
    }

    public function hasFrame(): bool
    {
        return false;
    }

    /**
     * @param non-empty-string $input
     * @return resource
     */
    private static function openIn(string $input)
    {
        $resource = @\fopen("php://$input", 'rb');

        if ($resource === false) {
            throw new RelayFactoryException('Could not initiate input stream resource');
        }

        return $resource;
    }

    /**
     * @param non-empty-string $output
     * @return resource
     */
    private static function openOut(string $output)
    {
        $resource = @\fopen("php://$output", 'wb');

        if ($resource === false) {
            throw new RelayFactoryException('could not initiate output stream resource');
        }

        return $resource;
    }
}
