<?php

declare(strict_types=1);

namespace Spiral\RoadRunner;

interface StreamWorkerInterface extends WorkerInterface
{
    public function withStreamMode(): static;
}
