<?php

/**
 * This file is part of RoadRunner package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\RoadRunner\Environment;

/**
 * @psalm-type ModeType = Mode::MODE_*
 */
interface Mode
{
    public const MODE_HTTP = 'http';
    public const MODE_TEMPORAL = 'temporal';
    public const MODE_JOBS = 'jobs';
    public const MODE_GRPC = 'grpc';
    public const MODE_TCP = 'tcp';
    public const MODE_CENTRIFUGE = 'centrifuge';
}
