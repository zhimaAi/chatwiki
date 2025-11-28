<?php

declare(strict_types=1);

namespace Spiral\Goridge;

use Spiral\Goridge\Exception\RelayException;

/**
 * This interface describes a Relay that explictily establishes a connection.
 * That connection can also be re-established on the fly (in comparison to StreamRelay, which relies on the existence of the streams).
 * The object is also clonable, i.e. supports cloning without data errors due to shared state.
 */
interface ConnectedRelayInterface extends RelayInterface
{
    /**
     * Returns true if the underlying connection is already established
     */
    public function isConnected(): bool;

    /**
     * Establishes the underlying connection and returns true on success, false on failure, or throws an exception in case of an error.
     *
     * @throws RelayException
     */
    public function connect(): bool;

    /**
     * Closes the underlying connection.
     */
    public function close(): void;

    /**
     * Enforce implementation of __clone magic method
     * @psalm-return void
     */
    public function __clone();
}
