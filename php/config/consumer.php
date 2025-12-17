<?php

// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

$params = require __DIR__ . '/params.php';

return [
    'id' => 'chatwiki',
    'basePath' => dirname(__DIR__),
    'timeZone' => 'PRC',
    'bootstrap' => ['log'],
    'aliases' => [

    ],
    'components' => [
        'log' => [
            'targets' => [
                [
                    'levels' => YII_DEBUG ? ['error', 'warning', 'info'] : ['error', 'warning'],
                    'class' => 'app\components\GoTarget',
                    'exportInterval' => 1,
                    'logVars' => [],
                ],
            ],
        ],
    ],
    'modules' => [

    ],
    'params' => $params,
];
