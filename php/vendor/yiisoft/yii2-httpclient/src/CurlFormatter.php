<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient;

use Yii;
use yii\base\BaseObject;

/**
 * CURLFormatter is used with CurlTransport to format the content of the request as an array
 * with the field name as key and field data as value
 * @see CURLOPT_POSTFIELDS
 * @since 2.0.9
 */
class CurlFormatter extends BaseObject implements FormatterInterface
{
    /**
     * @var int URL encoding type.
     * Possible values are:
     *  - PHP_QUERY_RFC1738 - encoding is performed per 'RFC 1738' and the 'application/x-www-form-urlencoded' media type,
     *    which implies that spaces are encoded as plus (+) signs. This is most common encoding type used by most web
     *    applications.
     *  - PHP_QUERY_RFC3986 - then encoding is performed according to 'RFC 3986', and spaces will be percent encoded (%20).
     *    This encoding type is required by OpenID and OAuth protocols.
     */
    public $encodingType = PHP_QUERY_RFC1738;


    /**
     * {@inheritdoc}
     */
    public function format(Request $request)
    {
        $data = $request->getData();

        if (strcasecmp('GET', $request->getMethod()) === 0) {
            if ($data !== null) {
                $content = http_build_query((array)$data, '', '&', $this->encodingType);
                $request->setFullUrl(null);
                $url = $request->getFullUrl();
                $url .= (strpos($url, '?') === false) ? '?' : '&';
                $url .= $content;
                $request->setFullUrl($url);
            }
            return $request;
        }

        if ($data !== null) {
            $request->setContent($data);
        } else {
            $request->getHeaders()->set('Content-Length', '0');
        }

        return $request;
    }
}
