<?php

declare(strict_types=1);

namespace Spiral\RoadRunner;

use Psr\Log\LoggerInterface;
use Psr\Log\LoggerTrait;

class Logger implements LoggerInterface
{
    use LoggerTrait;

    /**
     * @psalm-suppress RedundantConditionGivenDocblockType
     */
    public function log(mixed $level, string|\Stringable $message, array $context = []): void
    {
        \assert(\is_scalar($level), 'Invalid log level type');
        \assert(\is_string($message), 'Invalid log message type');

        $this->write($this->format((string) $level, $message, $context));
    }

    protected function write(string $message): void
    {
        \file_put_contents('php://stderr', $message);
    }

    protected function format(string $level, string $message, array $context = []): string
    {
        return \sprintf('[php %s] %s %s', $level, $message, $this->formatContext($context));
    }

    protected function formatContext(array $context): string
    {
        try {
            return \json_encode($context, \JSON_THROW_ON_ERROR);
        } catch (\JsonException) {
            return \print_r($context, true);
        }
    }
}
