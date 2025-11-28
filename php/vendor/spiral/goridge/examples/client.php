<?php

declare(strict_types=1);

use Spiral\Goridge;

require 'vendor/autoload.php';

$rpc = new Goridge\RPC\RPC(
    Goridge\Relay::create('tcp://127.0.0.1:6001')
);

echo $rpc->call('App.Hi', 'Antony');
