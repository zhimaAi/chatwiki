<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient;

use yii\base\Component;

/**
 * Transport performs actual HTTP request sending.
 *
 * @author Paul Klimov <klimov.paul@gmail.com>
 * @since 2.0
 */
abstract class Transport extends Component
{
    /**
     * Performs given request.
     * @param Request $request request to be sent.
     * @return Response response instance.
     * @throws Exception on failure.
     */
    abstract public function send($request);

    /**
     * Performs multiple HTTP requests.
     * Particular transport may benefit from this method, allowing sending requests in parallel.
     * This method accepts an array of the [[Request]] objects and returns an array of the  [[Response]] objects.
     * Keys of the response array correspond the ones from request array.
     * @param Request[] $requests requests to perform.
     * @return Response[] responses list.
     * @throws Exception
     */
    public function batchSend(array $requests)
    {
        $responses = [];
        foreach ($requests as $key => $request) {
            $responses[$key] = $this->send($request);
        }
        return $responses;
    }
}