<?php

declare(strict_types=1);

namespace Spiral\RoadRunner\Message\Command;

use Spiral\RoadRunner\Message\SkipMessage;
use Spiral\RoadRunner\Payload;

/**
 * @psalm-immutable
 */
final class Pong extends Payload implements SkipMessage {}
