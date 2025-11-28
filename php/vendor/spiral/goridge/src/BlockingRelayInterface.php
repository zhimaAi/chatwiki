<?php

declare(strict_types=1);

namespace Spiral\Goridge;

/**
 * Means that relay can't be used for non-blocking flow.
 */
interface BlockingRelayInterface extends RelayInterface {}
