<?php

declare(strict_types=1);

namespace Spiral\Goridge;

enum SocketType: string
{
    case TCP = 'tcp';
    case UNIX = 'unix';
}
