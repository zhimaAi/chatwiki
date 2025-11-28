<?php

declare(strict_types=1);

namespace Spiral\RoadRunner;

use Spiral\Goridge\RPC\Codec\JsonCodec;
use Spiral\Goridge\RPC\RPCInterface;
use Spiral\RoadRunner\Informer\Workers;
use Spiral\RoadRunner\Informer\Worker as InformerWorker;

/**
 * @psalm-type TInformerWorker = array{
 *     pid: positive-int,
 *     status: int,
 *     numExecs: int,
 *     created: positive-int,
 *     memoryUsage: positive-int,
 *     CPUPercent: float,
 *     command: string,
 *     statusStr: string,
 * }
 */
final class WorkerPool
{
    private readonly RPCInterface $rpc;

    public function __construct(
        RPCInterface $rpc,
    ) {
        $this->rpc = $rpc->withCodec(new JsonCodec());
    }

    /**
     * Add worker to the pool.
     *
     * @param non-empty-string $plugin
     */
    public function addWorker(string $plugin): void
    {
        $this->rpc->call('informer.AddWorker', $plugin);
    }

    /**
     * Get the number of workers for the pool.
     *
     * @param non-empty-string $plugin
     */
    public function countWorkers(string $plugin): int
    {
        return \count($this->getWorkers($plugin));
    }

    /**
     * Get the info about running workers for the pool.
     *
     * @param non-empty-string $plugin
     */
    public function getWorkers(string $plugin): Workers
    {
        /**
         * @var array{workers: list<TInformerWorker>} $data
         */
        $data = $this->rpc->call('informer.Workers', $plugin);

        return new Workers(\array_map(static function (array $worker): InformerWorker {
            return new InformerWorker(
                pid: $worker['pid'],
                statusCode: $worker['status'],
                executions: $worker['numExecs'],
                createdAt: $worker['created'],
                memoryUsage: $worker['memoryUsage'],
                cpuUsage: $worker['CPUPercent'],
                command: $worker['command'],
                status: $worker['statusStr'],
            );
        }, $data['workers']));
    }

    /**
     * Remove worker from the pool.
     *
     * @param non-empty-string $plugin
     */
    public function removeWorker(string $plugin): void
    {
        $this->rpc->call('informer.RemoveWorker', $plugin);
    }
}
