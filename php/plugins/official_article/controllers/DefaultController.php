<?php

// Copyright © 2016- 2025 Sesame Network Technology all right reserved

namespace app\plugins\official_article\controllers;

use app\controllers\ExtensionController;
use Exception;

class DefaultController extends ExtensionController
{
    public function getApiSchema(): array
    {
        $output = [
            'res' => ['type' => 'number','desc' => '错误码，非 0 表示失败'],
            'msg' => ['type' => 'string', 'desc' => '错误描述'],
        ];

        return [
            'get_login_url' => [
                'method' => 'POST',
                'path' => '/api/get_login_url',
                'title' => '获取爬虫登录地址',
                'type' => 'service',
                'desc' => '获取爬虫登录地址',
                'params' => [],
                'output' => $output + [
                        'data' => [
                            'type' => 'object',
                            'desc' => '爬虫登录地址',
                            'properties' => [
                                'url' => ['type' => 'string', 'desc' => '爬虫登录地址'],
                            ],
                        ],
                    ],
            ],
            'register' => [
                'method' => 'POST',
                'path' => '/api/register',
                'title' => '注册账号',
                'type' => 'service',
                'desc' => '注册账号',
                'params' => [
                    'username' => ['required' => true, 'in' => 'body', 'type' => 'string', 'name' => '用户名', 'desc' => '用户名'],
                    'password' => ['required' => false, 'in' => 'body', 'type' => 'string', 'name' => '密码', 'desc' => '密码'],
                ],
                'output' => $output + [
                    'data' => [
                        'type' => 'object',
                        'desc' => '用户信息',
                        'properties' => [
                            'user' => [
                                'type' => 'object',
                                'desc' => '用户信息',
                                'properties' => [
                                    'is_admin' => ['type' => 'number', 'desc' => '是否是管理员'],
                                    'username' => ['type' => 'string', 'desc' => '用户名'],
                                ]
                            ],
                        ],
                    ],
                ],
            ],

            'wechat_qrcode_login' => [
                'method' => 'POST',
                'path' => '/api/wechat/qrcode',
                'title' => '微信扫码登录',
                'type' => 'service',
                'desc' => '微信扫码登录',
                'params' => [
                    'username' => ['required' => true, 'in' => 'body', 'type' => 'string', 'name' => '用户名', 'desc' => '用户名‘'],
                ],
                'output' => $output + [
                    'data' => [
                        'type' => 'object',
                        'desc' => '二维码信息',
                        'properties' => [
                            'qrcode_base64' => ['type' => 'string', 'desc' => '二维码base64'],
                        ],
                    ]
                ],
            ],

            'get_login_status' => [
                'method' => 'POST',
                'path' => '/api/wechat/status',
                'title' => '获取登录状态',
                'type' => 'service',
                'desc' => '获取登录状态',
                'params' => [
                    'username' => ['required' => true, 'in' => 'body', 'type' => 'string', 'name' => '用户名', 'desc' => '用户名'],
                ],
                'output' => $output + [
                    'data' => [
                        'type' => 'object',
                        'desc' => '登录状态',
                        'properties' => [
                            'online' => ['type' => 'boole', 'desc' => '是否在线'],
                            'nickname' => ['type' => 'string', 'desc' => '昵称'],
                            'headimgurl' => ['type' => 'string', 'desc' => '头像URL'],
                            'login_time' => ['type' => 'integer', 'desc' => '登录时间'],
                            'logout_time' => ['type' => 'register', 'desc' => '退出时间'],
                            'login_duration_seconds' => ['type' => 'integer', 'desc' => '登录时长(秒)'],
                            'login_duration_text' => ['type' => 'string', 'desc' => '登录时长(文本)'],
                        ],
                    ],
                ],
            ],

            // 获取公众号文章列表
            'get_official_article' => [
                'method' => 'POST',
                'path' => '/articles-by-url',
                'title' => '公众号文章',
                'type' => 'node',
                'desc' => '获取任意公众号的最近及历史文章列表，使用任意一篇文章链接即可',
                'params' => [
                    'username' => ['required' => true, 'in' => 'body', 'type' => 'string', 'name' => '用户名', 'desc' => '用户名'],
                    'url' => ['required' => true, 'in' => 'body', 'type' => 'string', 'name' => '公众号文章链接' ,'desc' => '在公众号发布的任何文章的URL，例如: https://mp.weixin.qq.com/s/j4REV58ZPeaLFWSVuiYjSA'],
                    'number' => ['required' => true, 'in' => 'body', 'type' => 'integer', 'name' => '文章数量', 'desc' => '文章列表个数,默认是5', 'default' => 5],
                ],
                'output' => $output + [
                    'data' => [
                        'type' => 'array<object>',
                        'desc' => '文章列表',
                        'items' => [
                            'type' => 'object',
                            'desc' => '文章信息',
                            'properties' => [],
                        ],
                    ],
                ],
            ],
        ];
    }

    /**
     * 响应检查逻辑
     */
    private function checkBusinessSuccess(array $result): bool
    {
        return $result['success'] ?? false;
    }

    /**
     * 错误信息提取
     */
    private function getBusinessErrorMsg(array $result): string
    {
        return $result['message'] ?? '获取失败';
    }

    /**
     * 参数构建器
     * @throws Exception
     */
    private function buildHttpParts(array $config, array $args): array
    {
        $path = $config['path'];
        $query = [];
        $body = [];
        $headers = [
            'Content-Type' => 'application/json;charset=utf-8'
        ];

        foreach ($config['params'] as $field => $rules) {
            $val = $args[$field] ?? null;
            $required = $rules['required'] ?? false;
            $in = $rules['in'] ?? 'body'; // 默认为 body
            $type = $rules['type'] ?? 'string';

            // 必填校验
            if ($required && ($val === null || $val === '')) {
                throw new Exception("缺少必填参数: {$field}");
            }

            // 选填且为空，跳过
            if ($val === null) {
                continue;
            }

            // 类型简单校验
            if ($type === 'integer' && !is_numeric($val)) {
                throw new Exception("参数 {$field} 必须是数字");
            }
            if (($type === 'array' || $type === 'object') && !is_array($val)) {
                throw new Exception("参数 {$field} 必须是数组/对象");
            }

            // 核心分拣逻辑
            switch ($in) {
                case 'path':
                    // 替换 URL 中的 {xxx}
                    $path = str_replace('{' . $field . '}', (string)$val, $path);
                    break;
                case 'query':
                    $query[$field] = $val;
                    break;
                case 'header':
                    $headers[$field] = $val;
                    break;
                case 'body':
                default:
                    $body[$field] = $val;
                    break;
            }
        }

        // 安全检查：防止 Path 变量没被替换干净
        if (preg_match('/\{.*?}/', $path)) {
            throw new Exception("URL 路径变量未完全替换，请检查参数: {$path}");
        }

        return [
            'path' => $path,
            'query' => $query,
            'body' => $body,
            'headers' => $headers
        ];
    }

    /**
     * 统一执行入口
     */
    public function actionExec(array $params)
    {
        $business = $params['business'] ?? '';
        $arguments = $params['arguments'] ?? [];
        if (empty($business)) {
            return $this->error("业务标识不能为空");
        }

        if ($business == "get_login_url") {
            return $this->success(['url' => getenv('WECHAT_ARTICLE_CRAWLER_HOST')]);
        }

        try {
            // 加载配置
            $schemaMap = $this->getApiSchema();
            if (!isset($schemaMap[$business])) {
                return $this->error("未定义的业务接口: [{$business}]");
            }
            $config = $schemaMap[$business];

            // 参数分拣与校验
            // 这一步把扁平的 arguments 变成了 HTTP 请求所需的 parts
            $httpParts = $this->buildHttpParts($config, $arguments);

            // 发送请求
            $crawlerUrl = getenv('WECHAT_ARTICLE_CRAWLER_HOST');
            $fullUrl = $crawlerUrl . $httpParts['path'];
            $httpParts['body']['api_token'] = getenv('WECHAT_ARTICLE_CRAWLER_API_TOKEN');
            \Yii::info("请求体: " . json_encode($httpParts['body']));
            $response = $this->sendRequest(
                $config['method'],
                $fullUrl,
                $httpParts['query'],
                $httpParts['body'],
                $httpParts['headers']
            );

            // 响应处理
//            if (!$response->getIsOk()) {
//                 return $this->error('HTTP请求失败: ' . $response->getStatusCode());
//            }

            $result = $response->getData();

            \Yii::info("请求结果: " . json_encode($result));

            // 业务层面的成功判断
            if ($this->checkBusinessSuccess($result)) {
                if ($business == 'get_official_article') {
                    return $this->success($result['data']['articles'] ?? []);
                } else {
                    return $this->success($result['data'] ?? []);
                }
            }

            return $this->error($this->getBusinessErrorMsg($result));

        } catch (\Throwable $e) {
            return $this->error('执行异常: ' . $e->getMessage());
        }
    }
}