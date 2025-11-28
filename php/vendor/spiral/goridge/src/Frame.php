<?php

declare(strict_types=1);

namespace Spiral\Goridge;

use Spiral\Goridge\Exception\InvalidArgumentException;

/**
 * @psalm-type FrameVersion = Frame::VERSION
 * @psalm-type FrameType = Frame::CONTROL | Frame::ERROR
 * @psalm-type FrameCodec = Frame::CODEC_*
 * @psalm-type FrameCodecValue = int-mask-of<FrameCodec>
 * @psalm-type FrameByte10 = Frame::BYTE10_*
 * @psalm-type FrameByte10Value = int-mask-of<FrameByte10>
 * @psalm-type HeaderList = array{0: int, 1: int, 2: int<0, max>, 3: FrameByte10Value, 4: int}
 */
final class Frame
{
    /**
     * Current protocol version.
     */
    public const VERSION = 0x01;

    /**
     * Control frame type
     */
    public const CONTROL = 0x01;

    /**
     * Error frame type
     */
    public const ERROR = 0x40;

    /*#@+
     * BYTE flags, it means, that we can set multiply flags from this group
     * using bitwise OR.
     */
    public const CODEC_RAW     = 0x04;
    public const CODEC_JSON    = 0x08;
    public const CODEC_MSGPACK = 0x10;
    public const CODEC_GOB     = 0x20;
    public const CODEC_PROTO   = 0x80;
    /*#@-*/

    /*#@+
     * BYTE10 flags, it means, that we can set multiply flags from this group
     * using bitwise OR.
     *
     * @var positive-int Flags for {@see $byte10}
     */
    public const BYTE10_STREAM = 0x01;
    public const BYTE10_STOP = 0x02;
    public const BYTE10_PING = 0x04;
    public const BYTE10_PONG = 0x08;

    /**
     * @psalm-var FrameByte10Value
     */
    public int $byte10 = 0;

    public int $byte11 = 0;

    /**
     * @param array<int> $options
     */
    public function __construct(
        public ?string $payload,
        public array $options = [],
        public int $flags = 0,
    ) {}

    /**
     * @return non-empty-string
     * @internal
     */
    public static function packFrame(Frame $frame): string
    {
        /** @var non-empty-string $header */
        $header = \pack(
            'CCL',
            self::VERSION << 4 | (\count($frame->options) + 3),
            $frame->flags,
            \strlen((string) $frame->payload),
        );

        if ($frame->options !== []) {
            $header .= \pack('LCCL*', \crc32($header), $frame->byte10, $frame->byte11, ...$frame->options);
        } else {
            $header .= \pack('LCC', \crc32($header), $frame->byte10, $frame->byte11);
        }

        return $header . (string) $frame->payload;
    }

    /**
     * Parse header and return [flags, num options, payload length, stream byte10, byte11].
     *
     * @param string $header 8 bytes.
     * @return HeaderList
     * @internal
     *
     * @psalm-suppress LessSpecificReturnStatement, MoreSpecificReturnType
     */
    public static function readHeader(string $header): array
    {
        return [
            \ord($header[1]),
            (\ord($header[0]) & 0x0F) - 3,
            \ord($header[2]) | \ord($header[3]) << 8 | \ord($header[4]) << 16 | \ord($header[5]) << 24,
            \ord($header[10]),
            \ord($header[11]),
        ];
    }

    /**
     * @param HeaderList $header
     *
     * @see self::readHeader()
     * @internal
     */
    public static function initFrame(array $header, string $body): Frame
    {
        /**
         * optimize?
         * @var array<int> $options
         */
        $options = \array_values(\unpack('L*', \substr($body, 0, $header[1] * 4)));

        $frame = new self(\substr($body, $header[1] * 4), $options, $header[0]);
        $frame->byte10 = $header[3] ?? 0;
        $frame->byte11 = $header[4] ?? 0;

        return $frame;
    }

    public function setFlag(int ...$flag): void
    {
        foreach ($flag as $f) {
            if ($f > 255) {
                throw new InvalidArgumentException('Flags can be byte only');
            }

            $this->flags |= $f;
        }
    }

    public function hasFlag(int $flag): bool
    {
        if ($flag > 255) {
            throw new InvalidArgumentException('Flags can be byte only');
        }

        return ($this->flags & $flag) !== 0;
    }

    public function setOptions(int ...$options): void
    {
        $this->options = $options;
    }
}
