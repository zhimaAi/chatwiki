<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient;

use yii\base\BaseObject;
use yii\helpers\Json;

/**
 * JsonFormatter formats HTTP message as JSON.
 *
 * @author Paul Klimov <klimov.paul@gmail.com>
 * @since 2.0
 */
class JsonFormatter extends BaseObject implements FormatterInterface
{
    /**
     * @var int the encoding options. For more details please refer to
     * <https://www.php.net/manual/en/function.json-encode.php>.
     */
    public $encodeOptions = 0;


    /**
     * {@inheritdoc}
     */
    public function format(Request $request)
    {
        $request->getHeaders()->set('Content-Type', 'application/json; charset=UTF-8');
        if (($data = $request->getData()) !== null) {
            $request->setContent(Json::encode($request->getData(), $this->encodeOptions));
        }
        return $request;
    }
}
