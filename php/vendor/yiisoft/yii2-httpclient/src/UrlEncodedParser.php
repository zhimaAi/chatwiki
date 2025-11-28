<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient;

use yii\base\BaseObject;

/**
 * UrlEncodedParser parses HTTP message content as 'application/x-www-form-urlencoded'.
 *
 * @author Paul Klimov <klimov.paul@gmail.com>
 * @since 2.0
 */
class UrlEncodedParser extends BaseObject implements ParserInterface
{
    /**
     * {@inheritdoc}
     */
    public function parse(Response $response)
    {
        $data = [];
        parse_str($response->getContent(), $data);
        return $data;
    }
}