<?php

declare(strict_types=1);

namespace Spiral\RoadRunner\Internal;

use Symfony\Component\VarDumper\Dumper\AbstractDumper;
use Symfony\Component\VarDumper\Dumper\CliDumper;
use Symfony\Component\VarDumper\Dumper\HtmlDumper;

/**
 * @internal StdoutHandler is an internal library class, please do not use it in your code.
 * @psalm-internal Spiral\RoadRunner
 */
final class StdoutHandler
{
    /**
     * @var string
     */
    private const PIPE_OUT = 'php://stderr';

    /**
     * @var string
     */
    private const ERROR_WRITING_HEADER =
        'Could not explicitly send a headers "%s" using PHP header() function. ' .
        'Please use RoadRunner response object instead';

    /**
     * @var positive-int|0
     */
    private const OB_CHUNK_SIZE = 1;

    /**
     * @param positive-int|0 $chunkSize
     */
    public static function register(int $chunkSize = self::OB_CHUNK_SIZE): void
    {
        \assert($chunkSize >= 0, 'Invalid chunk size argument value');

        self::restreamOutputBuffer($chunkSize);
        self::restreamHeaders();

        // Vendor packages
        self::restreamSymfonyDumper();
    }

    /**
     * Intercept all output headers writing.
     */
    private static function restreamHeaders(): void
    {
        \header_register_callback(static function (): void {
            $headers = \headers_list();

            if ($headers !== []) {
                \file_put_contents(self::PIPE_OUT, self::ERROR_WRITING_HEADER);
            }
        });
    }

    /**
     * Intercept all output buffer write.
     *
     * @param int<0, max> $chunkSize
     */
    private static function restreamOutputBuffer(int $chunkSize): void
    {
        \ob_start(static function (string $chunk, int $phase): void {
            $isWrite = ($phase & \PHP_OUTPUT_HANDLER_WRITE) === \PHP_OUTPUT_HANDLER_WRITE;

            if ($isWrite && $chunk !== '') {
                \file_put_contents(self::PIPE_OUT, $chunk);
            }
        }, $chunkSize);
    }

    private static function restreamSymfonyDumper(): void
    {
        if (\class_exists(AbstractDumper::class)) {
            AbstractDumper::$defaultOutput = self::PIPE_OUT;
            CliDumper::$defaultOutput = self::PIPE_OUT;
            HtmlDumper::$defaultOutput = self::PIPE_OUT;
        }
    }
}
