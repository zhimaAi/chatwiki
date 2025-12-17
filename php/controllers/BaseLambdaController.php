<?php

// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

namespace app\controllers;

use Spiral\Goridge\RPC\RPC;
use Yii;

class BaseLambdaController extends \yii\base\Controller
{
    protected function getRpc(): ?RPC
    {
        return Yii::$app->has('rpc') ? Yii::$app->get('rpc') : null;
    }

    public function bindActionParams($action, $params)
    {
        return [$params];
    }

    public function success(array $result)
    {
        return ['res' => 0, 'msg' => 'ok', 'data' => $result];
    }

    public function error(string $msg)
    {
        return ['res' => 1, 'msg' => $msg, 'data' => []];
    }
}
