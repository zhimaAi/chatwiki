<?php

// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

namespace app\components;

use Spiral\RoadRunner\Environment;
use yii\log\Target;
use yii\log\Logger;
use Spiral\Goridge\RPC\RPC;
use Spiral\Goridge;
use Throwable;

/**
 * GoTarget 通过 RPC 将日志发送到 Golang 服务进行统一处理
 */
class GoTarget extends Target
{
    private RPC $rpc;

    /**
     * @var array syslog levels
     */
    private $logLevels = [
        Logger::LEVEL_TRACE => "Trace",
        Logger::LEVEL_PROFILE_BEGIN => "Debug",
        Logger::LEVEL_PROFILE_END => "Debug",
        Logger::LEVEL_PROFILE => "Debug",
        Logger::LEVEL_INFO => "Info",
        Logger::LEVEL_WARNING => "Warning",
        Logger::LEVEL_ERROR => "Error",
    ];

    public function init()
    {
        parent::init();

        $rrEnv = Environment::fromGlobals();
        $address = $rrEnv->getRPCAddress();
        $this->rpc = new Goridge\RPC\RPC(
            Goridge\Relay::create($address)
        );
    }

    /**
     * 导出日志消息到 Golang 服务
     */
    public function export()
    {
        if (empty($this->messages)) {
            return;
        }

        foreach ($this->messages as $message) {
            list($text, $level, $category, $timestamp) = $message;
            $func = "AppLogger." . $this->logLevels[$level];

            try {
                $this->rpc->call($func, "[{$category}] {$text}");
            } catch (Throwable $e) {
                error_log("GoTarget RPC call failed: " . $e->getMessage());
            }
        }
    }
}
