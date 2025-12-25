<?php

// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

$params = require __DIR__ . '/params.php';

return [
    'id' => 'chatwiki',
    'basePath' => dirname(__DIR__),
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
        'cache' => [
            'class' => 'yii\caching\FileCache',
            'cachePath' => '@runtime/cache',   // 缓存目录，默认就是这个
            'keyPrefix' => 'php_plugin_',             // 避免多项目冲突
        ],
    ],
    'modules' => [
        
    ],
    'params' => $params,
];
