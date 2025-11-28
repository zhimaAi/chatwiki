<?php

declare(strict_types=1);

namespace Spiral\RoadRunner;

use Spiral\Goridge\Frame;
use Spiral\RoadRunner\Exception\RoadRunnerException;
use Spiral\RoadRunner\Message\Command\GetProcessId;
use Spiral\RoadRunner\Message\Command\Pong;
use Spiral\RoadRunner\Message\Command\StreamStop;
use Spiral\RoadRunner\Message\Command\WorkerStop;

final class PayloadFactory
{
    private const JSON_DECODE_FLAGS = \JSON_THROW_ON_ERROR;

    public static function fromFrame(Frame $frame): Payload
    {
        $payload = $frame->payload ?? '';

        if ($frame->hasFlag(Frame::CONTROL)) {
            return self::makeControl($payload);
        }

        if (($frame->byte10 & Frame::BYTE10_STOP) !== 0) {
            return new StreamStop($payload);
        }

        if (($frame->byte10 & Frame::BYTE10_PONG) !== 0) {
            return new Pong($payload);
        }

        return new Payload(
            \substr($payload, $frame->options[0]),
            \substr($payload, 0, $frame->options[0]),
        );
    }

    private static function makeControl(string $header): Payload
    {
        try {
            $command = self::decode($header);
        } catch (\JsonException $e) {
            throw new RoadRunnerException('Invalid task header, JSON payload is expected: ' . $e->getMessage());
        }

        if (!empty($command['stop'])) {
            return new WorkerStop(null, $header);
        }

        if (!empty($command['pid'])) {
            return new GetProcessId(null, $header);
        }

        throw new RoadRunnerException('Invalid task header, undefined control package');
    }

    /**
     * @throws \JsonException
     * @psalm-assert non-empty-string $json
     */
    private static function decode(string $json): array
    {
        $result = \json_decode($json, true, 512, self::JSON_DECODE_FLAGS);

        if (! \is_array($result)) {
            throw new \JsonException('Json message must be an array or object');
        }

        return $result;
    }
}
