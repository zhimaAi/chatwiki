<?php

declare(strict_types=1);

namespace Spiral\RoadRunner;

interface WorkerAwareInterface
{
    /**
     * Returns underlying binary worker.
     */
    public function getWorker(): WorkerInterface;
}
