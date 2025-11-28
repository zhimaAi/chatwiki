<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient;

/**
 * ParserInterface represents HTTP response message parser.
 *
 * @author Paul Klimov <klimov.paul@gmail.com>
 * @since 2.0
 */
interface ParserInterface
{
    /**
     * Parses given HTTP response instance.
     * @param Response $response HTTP response instance.
     * @return mixed parsed content data.
     */
    public function parse(Response $response);
}