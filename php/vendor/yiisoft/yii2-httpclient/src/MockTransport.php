<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient;

use Yii;

final class MockTransport extends Transport
{
    /**
     * @var Request[]
     */
    private $requests = [];
    /**
     * @var Response[]
     */
    private $responses = [];


    /**
     * @param Response $response
     */
    public function appendResponse(Response $response)
    {
        $this->responses[] = $response;
    }

    /**
     * @return Request[]
     */
    public function flushRequests()
    {
        $requests = $this->requests;
        $this->requests = [];

        return $requests;
    }

    /**
     * {@inheritdoc}
     */
    public function send($request)
    {
        if (empty($this->responses)) {
            throw new Exception('No Response available');
        }

        $nextResponse = array_shift($this->responses);
        if (null === $nextResponse->client) {
            $nextResponse->client = $request->client;
        }

        $this->requests[] = $request;

        return $nextResponse;
    }
}
