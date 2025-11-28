<?php

// Copyright © 2016- 2025 Sesame Network Technology all right reserved

namespace app\controllers;

class SiteController extends \yii\web\Controller
{
    public function actionIndex()
    {
        return $this->asJson(['message' => 'hello world']);
    }
}