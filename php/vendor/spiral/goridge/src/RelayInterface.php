<?php

declare(strict_types=1);

namespace Spiral\Goridge;

use Spiral\Goridge\Exception\RelayException;

/**
 * Blocking, duplex relay.
 */
interface RelayInterface
{
    /**
     * @throws RelayException
     */
    public function waitFrame(): Frame;

    public function send(Frame $frame): void;

    public function hasFrame(): bool;
}
