<?php

// Copyright © 2016- 2025 Sesame Network Technology all right reserved

use Spiral\RoadRunner\Worker;
use Spiral\RoadRunner\Payload;
use Spiral\Roadrunner\Environment;
use yii\console\Application as ConsoleApplication;
use yii\web\Application as WebApplication;
use Spiral\Goridge;

ini_set('display_errors', 'stderr');
defined('YII_DEBUG') or define('YII_DEBUG', true);
defined('YII_ENV') or define('YII_ENV', 'dev');

require(__DIR__ . '/vendor/autoload.php');
require(__DIR__ . '/vendor/yiisoft/yii2/Yii.php');

$worker = Worker::create();
$rrEnv = Environment::fromGlobals();

$address = $rrEnv->getRPCAddress();
$rpc = new Goridge\RPC\RPC(
    Goridge\Relay::create($address)
);

if ($rrEnv->getMode() == "lambda") {
    $lambdaConfig = require __DIR__ . '/config/lambda.php';
    $lambdaApp = (new ConsoleApplication($lambdaConfig));
    $lambdaApp->set('rpc', $rpc);

    while ($payload = $worker->waitPayload()) {
        try {
            $body = json_decode($payload->body, true);
            if (empty($body) || !is_array($body)) {
                throw new Exception("data 字段为空或不是数组");
            }

            $plugin = $body['plugin'] ?? null;
            $action = $body['action'] ?? null;
            $params = $body['params'] ?? [];
            if (!$plugin || !$action) {
                throw new Exception("不合法的协议");
            }
            $lambdaApp->setModule($plugin, "\\app\\plugins\\{$plugin}\\Module");
            $module = $lambdaApp->getModule($plugin);
            if (!$module) {
                throw new \Exception("plugin not found: $plugin");
            }
            $result = $module->runAction($action, $params);
            $worker->respond(new Payload(json_encode($result)));
        } catch (Throwable $e) {
            $worker->error($e);
        }
    }
} elseif ($rrEnv->getMode() == "consumer") {
    $consumerConfig = require __DIR__ . '/config/consumer.php';
    $consumerApp = (new ConsoleApplication($consumerConfig));
    $consumerApp->set('rpc', $rpc);

    throw new \Exception("暂不支持");
} elseif ($rrEnv->getMode() == "web") {
    $webConfig = require __DIR__ . '/config/web.php';
    $webApplication = new WebApplication($webConfig);
    $webApplication->set('rpc', $rpc);
    throw new \Exception("暂不支持");
} else {
    throw new \Exception("不支持的模式");
}

