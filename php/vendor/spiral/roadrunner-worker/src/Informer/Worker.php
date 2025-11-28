<?php

declare(strict_types=1);

namespace Spiral\RoadRunner\Informer;

final class Worker
{
    /**
     * @param positive-int $pid process id
     * @param int $statusCode integer status of the worker
     * @param int $executions number of worker executions
     * @param positive-int $createdAt unix nano timestamp of worker creation time
     * @param positive-int $memoryUsage memory usage in bytes. Values might vary for different operating systems and based on RSS
     * @param float $cpuUsage how many percent of the CPU time this process uses
     * @param string $command used in the service plugin and shows a command for the particular service
     */
    public function __construct(
        public int $pid,
        public int $statusCode,
        public int $executions,
        public int $createdAt,
        public int $memoryUsage,
        public float $cpuUsage,
        public string $command,
        public string $status,
    ) {}
}
