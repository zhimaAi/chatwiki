<?php

// Copyright © 2016- 2025 Sesame Network Technology all right reserved

namespace app\controllers;

use LogicException;
use yii\httpclient\Client;
use yii\httpclient\Response;

abstract class ExtensionController extends BaseLambdaController
{
    abstract public function getApiSchema();
    abstract public function actionExec(array $params);

    /**
     * 暴露给前端的 Schema 接口
     */
    public function actionGetSchema(array $params)
    {
        return $this->success($this->getApiSchema());
    }

    /**
     * 参数校验接口
     * 校验传入的参数是否符合 API schema 的要求
     */
    public function actionCheckSchema(array $params)
    {
        $business = $params['business'] ?? '';
        $arguments = $params['arguments'] ?? [];

        if (empty($business)) {
            return $this->error("业务标识不能为空");
        }
        if (empty($arguments)) {
            return $this->error("参数不能为空");
        }

        try {
            // 加载配置
            $schemaMap = $this->getApiSchema();
            if (!isset($schemaMap[$business])) {
                throw new LogicException("未定义的业务接口: [{$business}]");
            }
            $config = $schemaMap[$business];

            // 遍历参数定义，逐一校验
            foreach ($config['params'] as $field => $rules) {
                $val = $arguments[$field] ?? null;
                $required = $rules['required'] ?? false;

                // 必填校验
                if ($required && ($val === null || $val === '')) {
                    throw new LogicException("缺少必填参数: {$field}");
                }

                // 选填且为空，跳过
                if ($val === null) {
                    continue;
                }

                // 类型校验
                $type = $rules['type'] ?? 'string';
                $this->validateValue($val, $type, $field, $rules);
            }

        } catch (\Throwable $e) {
            return $this->error('校验异常: ' . $e->getMessage());
        }
        return $this->success([]);
    }

    /**
     * 值校验器
     * @throws Exception
     */
    public function validateValue($val, string $type, string $field, array $rules): void
    {
        // 基础类型校验
        switch ($type) {
            case 'integer':
                if (!is_numeric($val)) {
                    throw new LogicException("参数 {$field} 必须是数字");
                }
                $numVal = (int)$val;
                if (isset($rules['min']) && $numVal < $rules['min']) {
                    throw new LogicException("参数 {$field} 不能小于 {$rules['min']}");
                }
                if (isset($rules['max']) && $numVal > $rules['max']) {
                    throw new LogicException("参数 {$field} 不能大于 {$rules['max']}");
                }
                break;

            case 'array':
                if (!is_array($val)) {
                    throw new LogicException("参数 {$field} 必须是数组");
                }
                // 数组长度校验
                $count = count($val);
                if (isset($rules['min']) && $count < $rules['min']) {
                    throw new LogicException("参数 {$field} 数组长度不能小于 {$rules['min']}");
                }
                if (isset($rules['max']) && $count > $rules['max']) {
                    throw new LogicException("参数 {$field} 数组长度不能大于 {$rules['max']}");
                }
                // 校验数组元素
                if (isset($rules['items']) && !empty($val)) {
                    foreach ($val as $idx => $item) {
                        $itemRules = $rules['items'];
                        $itemType = $itemRules['type'] ?? 'string';
                        $this->validateValue($item, $itemType, "{$field}[{$idx}]", $itemRules);
                    }
                }
                break;

            case 'object':
                if (!is_array($val)) {
                    throw new LogicException("参数 {$field} 必须是对象");
                }
                break;

            case 'boolean':
                if (!is_bool($val)) {
                    throw new LogicException("参数 {$field} 必须是布尔值");
                }
                break;

            case 'json':
                // 允许 JSON 字符串或数组/对象
                if (is_string($val)) {
                    if (!json_validate($val)) {
                        throw new LogicException("参数 {$field} 必须是有效的 JSON 字符串");
                    }
                } elseif (!is_array($val) && !is_object($val)) {
                    throw new LogicException("参数 {$field} 必须是 JSON 字符串或对象/数组");
                }
                break;

            case 'enum':
                if (!isset($rules['enum']) || !is_array($rules['enum'])) {
                    throw new LogicException("参数 {$field} 的枚举值未定义");
                }
                if (!in_array($val, $rules['enum'], true)) {
                    $allowed = implode(', ', $rules['enum']);
                    throw new LogicException("参数 {$field} 的值必须是: {$allowed} 中的一个");
                }
                break;

            case 'string':
            default:
                if (!is_string($val)) {
                    throw new LogicException("参数 {$field} 必须是字符串");
                }
                // 字符串长度校验
                $len = strlen($val);
                if (isset($rules['min']) && $len < $rules['min']) {
                    throw new LogicException("参数 {$field} 字符串长度不能小于 {$rules['min']}");
                }
                if (isset($rules['max']) && $len > $rules['max']) {
                    throw new LogicException("参数 {$field} 字符串长度不能大于 {$rules['max']}");
                }
                break;
        }
    }

    /**
     * HTTP 封装
     */
    public function sendRequest(string $method, string $url, array $query, array $body, array $headers): Response
    {
        // 拼接 Query 参数
        if (!empty($query)) {
            $url .= (str_contains($url, '?') ? '&' : '?') . http_build_query($query);
        }

        $client = new Client();
        $request = $client->createRequest()
            ->setMethod($method)
            ->setUrl($url)
            ->addHeaders($headers)
            ->setOptions([
                'timeout' => 15, // 默认超时
            ]);

        // GET 请求不带 Body
        if (!empty($body) && strtoupper($method) !== 'GET') {
            $request->setContent(json_encode($body, JSON_UNESCAPED_UNICODE));
        }

        return $request->send();
    }
}
