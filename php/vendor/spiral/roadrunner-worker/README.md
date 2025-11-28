<a href="https://roadrunner.dev" target="_blank">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/roadrunner-server/.github/assets/8040338/e6bde856-4ec6-4a52-bd5b-bfe78736c1ff">
    <img align="center" src="https://github.com/roadrunner-server/.github/assets/8040338/040fb694-1dd3-4865-9d29-8e0748c2c8b8">
  </picture>
</a>
<p align="center">
 <a href="https://packagist.org/packages/spiral/roadrunner"><img src="https://poser.pugx.org/spiral/roadrunner/version"></a>
	<a href="https://pkg.go.dev/github.com/spiral/roadrunner?tab=doc"><img src="https://godoc.org/github.com/spiral/roadrunner?status.svg"></a>
	<a href="https://github.com/spiral/roadrunner/actions"><img src="https://github.com/spiral/roadrunner/workflows/CI/badge.svg" alt=""></a>
	<a href="https://goreportcard.com/report/github.com/spiral/roadrunner"><img src="https://goreportcard.com/badge/github.com/spiral/roadrunner"></a>
	<a href="https://scrutinizer-ci.com/g/spiral/roadrunner/?branch=master"><img src="https://scrutinizer-ci.com/g/spiral/roadrunner/badges/quality-score.png"></a>
	<a href="https://discord.gg/TFeEmCs"><img src="https://img.shields.io/badge/discord-chat-magenta.svg"></a>
	<a href="https://packagist.org/packages/spiral/roadrunner"><img src="https://img.shields.io/packagist/dd/spiral/roadrunner?style=flat-square"></a>
</p>

RoadRunner is an open-source (MIT licensed) high-performance PHP application server, load balancer, and process manager.
It supports running as a service with the ability to extend its functionality on a per-project basis.

RoadRunner includes PSR-7/PSR-17 compatible HTTP and HTTP/2 server and can be used to replace classic Nginx+FPM setup with much greater performance and flexibility.

<p align="center">
	<a href="https://roadrunner.dev/"><b>Official Website</b></a> | 
	<a href="https://docs.roadrunner.dev"><b>Documentation</b></a>
</p>

Repository:
--------

This repository contains the common codebase for all binary roadrunner workers. 
Check [spiral/roadrunner](https://github.com/spiral/roadrunner) to access application
server and [spiral/roadrunner-http](https://github.com/spiral/roadrunner-http) for PSR-7 compatible worker.

You can use the convenient installer to download the latest available compatible version of RoadRunner assembly:

```bash
$ composer require spiral/roadrunner-cli --dev
```

To download latest version of application server:

```bash
$ vendor/bin/rr get
```

Example:
-------

To init abstract RoadRunner worker:

```php
<?php

require __DIR__ . '/vendor/autoload.php';

// Create a new Worker from global environment
$worker = \Spiral\RoadRunner\Worker::create();

while ($data = $worker->waitPayload()) {
    // Received Payload
    var_dump($data);

    // Respond Answer
    $worker->respond(new \Spiral\RoadRunner\Payload('DONE'));
}
```

<a href="https://spiral.dev/">
<img src="https://user-images.githubusercontent.com/773481/220979012-e67b74b5-3db1-41b7-bdb0-8a042587dedc.jpg" alt="try Spiral Framework" />
</a>

Testing:
--------

This codebase is automatically tested via host repository - [spiral/roadrunner](https://github.com/spiral/roadrunner).

License:
--------

The MIT License (MIT). Please see [`LICENSE`](./LICENSE) for more information. Maintained by [Spiral Scout](https://spiralscout.com).
