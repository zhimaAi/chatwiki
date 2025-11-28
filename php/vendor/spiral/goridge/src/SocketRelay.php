<?php

declare(strict_types=1);

namespace Spiral\Goridge;

use Socket;
use Spiral\Goridge\Exception\HeaderException;
use Spiral\Goridge\Exception\InvalidArgumentException;
use Spiral\Goridge\Exception\RelayException;
use Spiral\Goridge\Exception\TransportException;

/**
 * Communicates with remote server/client over be-directional socket using byte payload:
 *
 * [ prefix       ][ payload                               ]
 * [ 1+8+8 bytes  ][ message length|LE ][message length|BE ]
 *
 * prefix:
 * [ flag       ][ message length, unsigned int 64bits, LittleEndian ]
 *
 * @psalm-type PortType = int<0, max>|null
 *
 * @psalm-suppress DeprecatedInterface
 */
class SocketRelay extends Relay implements \Stringable, ConnectedRelayInterface
{
    final public const RECONNECT_RETRIES = 10;
    final public const RECONNECT_TIMEOUT = 100;

    /**
     * @internal
     * This isn't really ideal but there's no easy way since we need access to the underlying socket
     * to do a socket_select across multiple SocketRelays.
     */
    public ?\Socket $socket = null;

    /**
     * 1) Pathname to "sock" file in case of UNIX socket
     * 2) URI string in case of TCP socket
     */
    private readonly string $address;

    /** @var PortType */
    private readonly ?int $port;

    private readonly SocketType $type;

    /**
     * Example:
     *
     * <code>
     *  $relay = new SocketRelay("localhost", 7000);
     *  $relay = new SocketRelay("/tmp/rpc.sock", null, SocketType::UNIX);
     * </code>
     *
     * @param non-empty-string $address Localhost, ip address or hostname.
     * @param PortType $port Ignored for UNIX sockets.
     *
     * @throws InvalidArgumentException
     */
    public function __construct(
        string $address,
        ?int $port = null,
        SocketType $type = SocketType::TCP,
    ) {
        // Guaranteed at the level of composer's json config
        \assert(\extension_loaded('sockets'));

        switch ($type) {
            case SocketType::TCP:
                // TCP address should always be in lowercase
                $address = \strtolower($address);

                if ($port === null) {
                    throw new InvalidArgumentException(\sprintf("Тo port given for TPC socket on '%s'", $address));
                }

                if ($port < 0 || $port > 65535) {
                    throw new InvalidArgumentException(\sprintf("Invalid port given for TPC socket on '%s'", $address));
                }

                break;

            case SocketType::UNIX:
                $port = null;
                break;
        }

        $this->address = $address;
        $this->port = $port;
        $this->type = $type;
    }

    public function getAddress(): string
    {
        return $this->address;
    }

    /**
     * @return PortType
     */
    public function getPort(): ?int
    {
        return $this->port;
    }

    public function getType(): SocketType
    {
        return $this->type;
    }

    /**
     * @psalm-assert-if-true Socket $this->socket
     * @psalm-assert-if-false null $this->socket
     */
    public function isConnected(): bool
    {
        return $this->socket !== null;
    }

    /**
     * @throws RelayException
     * @psalm-suppress PossiblyNullArgument Reason: Using the "connect()" method guarantees
     *                 the existence of the socket.
     */
    public function waitFrame(): Frame
    {
        $this->connect();

        $header = '';
        $headerLength = \socket_recv($this->socket, $header, 12, \MSG_WAITALL);

        if ($headerLength !== 12) {
            $error = \socket_strerror(\socket_last_error($this->socket));
            throw new HeaderException(\sprintf('Unable to read frame header: %s', $error));
        }

        $parts = Frame::readHeader($header);

        // total payload length
        $payload = '';
        $length = $parts[1] * 4 + $parts[2];

        while ($length > 0) {
            $bufferLength = \socket_recv($this->socket, $buffer, $length, \MSG_WAITALL);

            /**
             * Suppress "buffer === null" assertion, because buffer can contain
             * NULL in case of socket_recv function error.
             *
             * @psalm-suppress TypeDoesNotContainNull
             */
            if ($bufferLength === false || $buffer === null) {
                $message = \socket_strerror(\socket_last_error($this->socket));
                throw new HeaderException(\sprintf('Unable to read payload from socket: %s', $message));
            }

            $payload .= $buffer;
            $length -= $bufferLength;
        }

        return Frame::initFrame($parts, $payload);
    }

    /**
     * @psalm-suppress PossiblyNullArgument Reason: Using the "connect()" method guarantees
     *                 the existence of the socket.
     */
    public function send(Frame $frame): void
    {
        $this->connect();

        $body = Frame::packFrame($frame);

        if (\socket_send($this->socket, $body, \strlen($body), 0) === false) {
            throw new TransportException('Unable to write payload to the stream');
        }
    }

    public function hasFrame(): bool
    {
        if (!$this->isConnected()) {
            return false;
        }

        $read = [$this->socket];
        $write = null;
        $except = null;

        $is = \socket_select($read, $write, $except, 0);

        return $is > 0;
    }

    /**
     * Ensure socket connection. Returns true if socket successfully connected
     * or have already been connected.
     *
     * @param int<0, max> $retries Count of connection tries.
     * @param int<0, max> $timeout Timeout between reconnections in microseconds.
     *
     * @throws RelayException
     */
    public function connect(int $retries = self::RECONNECT_RETRIES, int $timeout = self::RECONNECT_TIMEOUT): bool
    {
        \assert($retries >= 1);
        \assert($timeout > 0);

        if ($this->isConnected()) {
            return true;
        }

        $socket = $this->createSocket();

        if ($socket === false) {
            throw new RelayException("Unable to create socket {$this}");
        }

        try {
            $status = false;

            for ($attempt = 0; $attempt <= $retries; ++$attempt) {
                if ($status = @\socket_connect($socket, $this->address, $this->port ?? 0)) {
                    break;
                }

                \usleep(\max(0, $timeout));
            }

            if ($status === false) {
                throw new RelayException(\socket_strerror(\socket_last_error($socket)));
            }
        } catch (\Throwable $e) {
            throw new RelayException("Unable to establish connection {$this}", 0, $e);
        }

        $this->socket = $socket;

        return true;
    }

    /**
     * Close connection.
     *
     * @throws RelayException
     */
    public function close(): void
    {
        if (!$this->isConnected()) {
            throw new RelayException("Unable to close socket '{$this}', socket already closed");
        }

        \socket_close($this->socket);
        $this->socket = null;
    }

    public function __toString(): string
    {
        if ($this->type === SocketType::TCP) {
            return "tcp://{$this->address}:{$this->port}";
        }

        return "unix://{$this->address}";
    }

    public function __clone()
    {
        // Remove reference to socket on clone
        $this->socket = null;
    }

    /**
     * Destruct connection and disconnect.
     */
    public function __destruct()
    {
        if ($this->isConnected()) {
            $this->close();
        }
    }

    private function createSocket(): \Socket|false
    {
        if ($this->type === SocketType::UNIX) {
            return \socket_create(\AF_UNIX, \SOCK_STREAM, 0);
        }

        return \socket_create(\AF_INET, \SOCK_STREAM, \SOL_TCP);
    }
}
