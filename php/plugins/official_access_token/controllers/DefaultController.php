<?php

namespace app\plugins\official_access_token\controllers;

use app\controllers\ExtensionController;

class DefaultController extends ExtensionController
{
    private const string BASE_URL = 'https://api.weixin.qq.com/cgi-bin';

    public function getApiSchema(): array
    {
        return [
            // 获取用户基本信息
            'get_access_token' => [
                'method' => 'GET',
                'path' => '/stable_token',
                'title' => '获取公众号access_token',
                'type' => 'node',
                'common_template' => true,
                'desc' => '本接口用于获取获取全局唯一后台接口调用凭据（Access Token），token 有效期为 7200 秒',
                'params' => [
                    'app_id' => ['required' => true, 'in' => 'query', 'type' => 'string', 'name' => '公众号', 'select_official_component'=>true,],
                    'app_secret' => ['required' => true, 'in' => 'query', 'type' => 'string', 'desc' => '公众号密钥', 'hide_official_component'=>true,],
                ],
                'output' => [
                    'res' => ['type' => 'number','desc' => '错误码，非 0 表示失败'],
                    'msg' => ['type' => 'string', 'desc' => '错误描述'],
                    'data' => [
                        'type' => 'object',
                        'desc' => '响应数据',
                        'properties' => [
                            'access_token' => ['type' => 'string', 'desc' => '获取到的凭证'],
                            'expires_in' => ['type' => 'number', 'desc' => '凭证有效时间，单位：秒。目前是7200秒之内的值。'],
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
        return !isset($result['errcode']) || $result['errcode'] == 0;
    }

    /**
     * 错误信息提取
     */
    private function getBusinessErrorMsg(array $result): string
    {
        return $result['errmsg'] ?? '微信公众号接口未知错误';
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

        $body = [
            'appid' => $arguments['app_id'] ?? '',
            'secret' => $arguments['app_secret'] ?? '',
            'grant_type' => 'client_credential',
            'force_refresh' => false,
        ];
        if (empty($body['appid']) || empty($body['secret'])) {
            return $this->error("缺少参数");
        }
        $headers = [
            'Content-Type' => 'application/json;charset=utf-8'
        ];
        try {
            $response = $this->sendRequest('POST', self::BASE_URL . "/stable_token", [], $body, $headers);
            $result = $response->getData();

            \Yii::info("请求结果: " . json_encode($result));

            // 业务层面的成功判断
            if ($this->checkBusinessSuccess($result)) {
                return $this->success($result);
            }

            return $this->error($this->getBusinessErrorMsg($result));
        } catch (\Throwable $e) {
            return $this->error('执行异常: ' . $e->getMessage());
        }
    }
}