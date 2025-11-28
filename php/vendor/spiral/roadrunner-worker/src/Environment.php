<?php

declare(strict_types=1);

namespace Spiral\RoadRunner;

use Spiral\RoadRunner\Environment\Mode;

/**
 * @psalm-import-type ModeType from Mode
 * @psalm-type EnvironmentVariables = array{
 *      RR_MODE?:    ModeType|string,
 *      RR_RELAY?:   string,
 *      RR_RPC?:     string,
 *      RR_VERSION?: string,
 * }|array<string, string>
 * @see Mode
 */
class Environment implements EnvironmentInterface
{
    /**
     * @param EnvironmentVariables $env
     */
    public function __construct(
        private array $env = [],
    ) {}

    public static function fromGlobals(): self
    {
        /** @var array<string, string> $env */
        $env = [...$_ENV, ...$_SERVER];

        return new self($env);
    }

    public function getMode(): string
    {
        return $this->get('RR_MODE');
    }

    public function getRelayAddress(): string
    {
        return $this->get('RR_RELAY', 'pipes');
    }

    public function getRPCAddress(): string
    {
        return $this->get('RR_RPC', 'tcp://127.0.0.1:6001');
    }

    public function getVersion(): string
    {
        return $this->get('RR_VERSION');
    }

    /**
     * @template TDefault of string
     *
     * @param non-empty-string $name
     * @param TDefault $default
     * @return string|TDefault
     */
    private function get(string $name, string $default = ''): string
    {
        if (isset($this->env[$name]) || \array_key_exists($name, $this->env)) {
            /** @psalm-suppress RedundantCastGivenDocblockType */
            return (string) $this->env[$name];
        }

        return $default;
    }
}
