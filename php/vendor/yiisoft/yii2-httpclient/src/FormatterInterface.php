<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient;

/**
 * FormatterInterface represents HTTP request message formatter.
 *
 * @author Paul Klimov <klimov.paul@gmail.com>
 * @since 2.0
 */
interface FormatterInterface
{
    /**
     * Formats given HTTP request message.
     * @param Request $request HTTP request instance.
     * @return Request formatted request.
     */
    public function format(Request $request);
}