<a href="https://roadrunner.dev" target="_blank">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/roadrunner-server/.github/assets/8040338/e6bde856-4ec6-4a52-bd5b-bfe78736c1ff">
    <img align="center" src="https://github.com/roadrunner-server/.github/assets/8040338/040fb694-1dd3-4865-9d29-8e0748c2c8b8">
  </picture>
</a>

# High-performance PHP-to-Golang IPC bridge

[![Latest Stable Version](https://poser.pugx.org/spiral/goridge/v/stable)](https://packagist.org/packages/spiral/goridge)
[![CI](https://github.com/spiral/goridge-php/workflows/CI/badge.svg)](https://github.com/spiral/goridge-php/actions)
[![Codecov](https://codecov.io/gh/roadrunner-php/goridge/branch/master/graph/badge.svg)](https://codecov.io/gh/roadrunner-php/goridge/)
[![Chat](https://img.shields.io/badge/discord-chat-magenta.svg)](https://discord.gg/TFeEmCs)

<img src="https://files.phpclasses.org/graphics/phpclasses/innovation-award-logo.png" height="90px" alt="PHPClasses Innovation Award" align="left"/>

Goridge is high performance PHP-to-Golang codec library which works over native PHP sockets and Golang net/rpc package. The library allows you to call Go service methods from PHP with minimal footprint, structures and `[]byte` support.  
Golang source code can be found in this repository: [goridge](https://github.com/roadrunner-server/goridge)

<br/>
See https://github.com/spiral/roadrunner - High-performance PHP application server, load-balancer and process manager written in Golang
<br/>

## Features

 - no external dependencies or services, drop-in (64bit PHP version required)
 - sockets over TCP or Unix (ext-sockets is required), standard pipes
 - very fast (300k calls per second on Ryzen 1700X over 20 threads)
 - native `net/rpc` integration, ability to connect to existing application(s)
 - standalone protocol usage
 - structured data transfer using json or msgpack
 - `[]byte` transfer, including big payloads
 - service, message and transport level error handling
 - hackable
 - works on Windows
 - unix sockets powered (also on Windows)

## Installation

```
composer require spiral/goridge
```

## Example

```php
<?php

use Spiral\Goridge;
require "vendor/autoload.php";

$rpc = new Goridge\RPC\RPC(
    Goridge\Relay::create('tcp://127.0.0.1:6001')
);

//or, using factory:
$tcpRPC = new Goridge\RPC\RPC(Goridge\Relay::create('tcp://127.0.0.1:6001'));
$unixRPC = new Goridge\RPC\RPC(Goridge\Relay::create('unix:///tmp/rpc.sock'));
$streamRPC = new Goridge\RPC\RPC(Goridge\Relay::create('pipes://stdin:stdout'));

echo $rpc->call("App.Hi", "Antony");
```

> Factory applies the next format: `<protocol>://<arg1>:<arg2>`

More examples can be found in [this directory](./examples).

License
-------
The MIT License (MIT). Please see [`LICENSE`](./LICENSE) for more information.
