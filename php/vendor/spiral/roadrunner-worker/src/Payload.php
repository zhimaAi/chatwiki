<?php

declare(strict_types=1);

namespace Spiral\RoadRunner;

use JetBrains\PhpStorm\Immutable;

/**
 * @internal
 */
#[Immutable]
class Payload
{
    /**
     * Execution payload (binary).
     *
     * @psalm-readonly
     */
    public readonly string $body;

    /**
     * Execution context (binary).
     *
     * @psalm-readonly
     */
    public readonly string $header;

    public function __construct(
        ?string $body,
        ?string $header = null,

        /**
         * End of stream.
         * The {@see true} value means the Payload block is last in the stream.
         *
         * @psalm-readonly
         */
        public readonly bool $eos = true,
    ) {
        $this->body = $body ?? '';
        $this->header = $header ?? '';
    }
}
