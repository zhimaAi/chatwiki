<?php

// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

namespace app\controllers;

class SiteController extends \yii\web\Controller
{
    public function actionIndex()
    {
        return $this->asJson(['message' => 'hello world']);
    }
}