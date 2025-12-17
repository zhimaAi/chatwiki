<?php

// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

namespace app\controllers;

abstract class NoticeController extends BaseLambdaController
{
    abstract public function actionSendMessage(array $params);

    abstract public function actionGetSchema();
}