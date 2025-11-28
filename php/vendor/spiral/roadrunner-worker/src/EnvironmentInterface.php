<?php

declare(strict_types=1);

namespace Spiral\RoadRunner;

use JetBrains\PhpStorm\ExpectedValues;
use Spiral\RoadRunner\Environment\Mode;

/**
 * Provides base values to configure roadrunner worker.
 *
 * @psalm-import-type ModeType from Mode
 * @see Mode
 * @method string getVersion()
 */
interface EnvironmentInterface
{
    /**
     * Returns worker mode assigned to the PHP process.
     *
     * @return ModeType|string
     */
    #[ExpectedValues(valuesFromClass: Mode::class)]
    public function getMode(): string;

    /**
     * Address worker should be connected to (or pipes).
     */
    public function getRelayAddress(): string;

    /**
     * RPC address.
     */
    public function getRPCAddress(): string;
}
