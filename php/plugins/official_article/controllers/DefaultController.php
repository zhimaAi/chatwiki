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
            // 获取公众号文章列表
            'get_official_article' => [
                'method' => 'POST',
                'path' => '/articles-by-url',
                'title' => '公众号文章',
                'type' => 'node',
                'desc' => '获取任意公众号的最近及历史文章列表，使用任意一篇文章链接即可',
                'params' => [
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
                            'properties' => [
                                // 'msgid' => ['type' => 'number', 'desc' => '消息ID'],
                                // 'aid' => ['type' => 'string', 'desc' => '文章ID'],
                                // 'title' => ['type' => 'string', 'desc' => '文章标题'],
                                // 'link' => ['type' => 'string', 'desc' => '文章链接'],
                                // 'cover' => ['type' => 'string', 'desc' => '文章封面'],
                                // 'digest' => ['type' => 'string', 'desc' => '简介'],
                                // 'create_time' => ['type' => 'number', 'desc' => '创建时间(时间戳)'],
                                // 'update_time' => ['type' => 'number', 'desc' => '更新时间(时间戳)'],
                            ],
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
                return $this->success($result['data']['articles'] ?? []);
            }

            return $this->error($this->getBusinessErrorMsg($result));

        } catch (\Throwable $e) {
            return $this->error('执行异常: ' . $e->getMessage());
        }
    }

}