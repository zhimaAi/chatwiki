<?php

// Copyright © 2016- 2025 Sesame Network Technology all right reserved

namespace app\controllers;

abstract class NoticeController extends BaseLambdaController
{
    abstract public function actionSendMessage(array $params);

    abstract public function actionGetSchema();
}