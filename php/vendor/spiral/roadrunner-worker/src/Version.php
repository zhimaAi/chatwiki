<?php

declare(strict_types=1);

namespace Spiral\RoadRunner;

use Composer\InstalledVersions;

final class Version
{
    public const PACKAGE_NAMES = [
        'spiral/roadrunner',
        'spiral/roadrunner-worker',
    ];
    public const VERSION_FALLBACK = 'dev-master';

    public static function current(): string
    {
        foreach (self::PACKAGE_NAMES as $name) {
            if (InstalledVersions::isInstalled($name)) {
                return \ltrim((string) InstalledVersions::getPrettyVersion($name), 'v');
            }
        }

        return self::VERSION_FALLBACK;
    }

    public static function constraint(): string
    {
        $current = self::current();

        if (\str_contains($current, '.')) {
            [$major] = \explode('.', $current);

            return \is_numeric($major) ? "$major.*" : '*';
        }

        return '*';
    }
}
