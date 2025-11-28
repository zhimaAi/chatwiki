<?php

declare(strict_types=1);

namespace Spiral\RoadRunner\Message;

/**
 * Marked message should be skipped in main worker loop in {@see \Spiral\RoadRunner\Worker::waitPayload()}.
 * For example {@see StreamStop} message makes sense only in stream output. The message can be received
 * after stream end because async and should be skipped in main worker loop.
 * @internal
 */
interface SkipMessage {}
