<?php

declare(strict_types=1);

namespace Spiral\RoadRunner\Message\Command;

use Spiral\RoadRunner\Message\ControlMessage;
use Spiral\RoadRunner\Payload;

/**
 * @psalm-immutable
 */
final class WorkerStop extends Payload implements ControlMessage {}
