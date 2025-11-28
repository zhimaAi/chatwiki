<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient;

use yii\web\Cookie;
use yii\web\HeaderCollection;

/**
 * Response represents HTTP request response.
 *
 * @property-read bool $isOk Whether response is OK.
 * @property-read string $statusCode Status code.
 *
 * @author Paul Klimov <klimov.paul@gmail.com>
 * @since 2.0
 */
class Response extends Message
{
    /**
     * {@inheritdoc}
     */
    public function getData()
    {
        $data = parent::getData();
        if ($data === null) {
            $content = $this->getContent();
            if (is_string($content) && strlen($content) > 0) {
                $data = $this->getParser()->parse($this);
                $this->setData($data);
            }
        }
        return $data;
    }

    /**
     * {@inheritdoc}
     */
    public function getCookies()
    {
        $cookieCollection = parent::getCookies();
        if ($cookieCollection->getCount() === 0 && $this->getHeaders()->has('set-cookie')) {
            $cookieStrings = $this->getHeaders()->get('set-cookie', [], false);
            foreach ($cookieStrings as $cookieString) {
                $cookieCollection->add($this->parseCookie($cookieString));
            }
        }
        return $cookieCollection;
    }

    /**
     * Returns status code.
     * @throws Exception on failure.
     * @return string status code.
     */
    public function getStatusCode()
    {
        $headers = $this->getHeaders();
        if ($headers->has('http-code')) {
            // take into account possible 'follow location'
            $statusCodeHeaders = $headers->get('http-code', null, false);
            return empty($statusCodeHeaders) ? null : end($statusCodeHeaders);
        }
        throw new Exception('Unable to get status code: referred header information is missing.');
    }

    /**
     * Checks if response status code is OK (status code = 2xx)
     * @return bool whether response is OK.
     * @throws Exception
     */
    public function getIsOk()
    {
        $statusCode = (int)$this->getStatusCode();
        return $statusCode >= 200 && $statusCode < 300;
    }

    /**
     * Returns default format automatically detected from headers and content.
     * @return string|null format name, 'null' - if detection failed.
     */
    protected function defaultFormat()
    {
        $format = $this->detectFormatByHeaders($this->getHeaders());
        if ($format === null) {
            $format = $this->detectFormatByContent($this->getContent());
        }

        return $format;
    }

    /**
     * Detects format from headers.
     * @param HeaderCollection $headers source headers.
     * @return null|string format name, 'null' - if detection failed.
     */
    protected function detectFormatByHeaders(HeaderCollection $headers)
    {
        $contentTypeHeaders = $headers->get('content-type', null, false);

        if (!empty($contentTypeHeaders)) {
            $contentType = end($contentTypeHeaders);
            if (stripos($contentType, 'json') !== false) {
                return Client::FORMAT_JSON;
            }
            if (stripos($contentType, 'urlencoded') !== false) {
                return Client::FORMAT_URLENCODED;
            }
            if (stripos($contentType, 'xml') !== false) {
                return Client::FORMAT_XML;
            }
        }

        return null;
    }

    /**
     * Detects response format from raw content.
     * @param string $content raw response content.
     * @return null|string format name, 'null' - if detection failed.
     */
    protected function detectFormatByContent($content)
    {
        if (preg_match('/^(\\{|\\[\\{).*(\\}|\\}\\])$/is', $content)) {
            return Client::FORMAT_JSON;
        }
        if (preg_match('/^([^=&])+=[^=&]+(&[^=&]+=[^=&]+)*$/', $content)) {
            return Client::FORMAT_URLENCODED;
        }
        if (preg_match('/^<\?xml.*>$/s', $content)) {
            return Client::FORMAT_XML;
        }
        return null;
    }

    /**
     * Parses cookie value string, creating a [[Cookie]] instance.
     * @param string $cookieString cookie header string.
     * @return Cookie cookie object.
     */
    private function parseCookie($cookieString)
    {
        $params = [];
        $pairs = explode(';', $cookieString);
        foreach ($pairs as $number => $pair) {
            $pair = trim($pair);
            if (strpos($pair, '=') === false) {
                $params[$this->normalizeCookieParamName($pair)] = true;
            } else {
                list($name, $value) = explode('=', $pair, 2);
                if ($number === 0) {
                    $params['name'] = $name;
                    $params['value'] = urldecode($value);
                } else {
                    $params[$this->normalizeCookieParamName($name)] = urldecode($value);
                }
            }
        }

        $cookie = new Cookie();
        foreach ($params as $name => $value) {
            if ($cookie->canSetProperty($name)) {
                // Cookie string may contain custom unsupported params
                $cookie->$name = $value;
            }
        }
        return $cookie;
    }

    /**
     * @param string $rawName raw cookie parameter name.
     * @return string name of [[Cookie]] field.
     */
    private function normalizeCookieParamName($rawName)
    {
        static $nameMap = [
            'expires' => 'expire',
            'httponly' => 'httpOnly',
            'max-age' => 'maxAge',
        ];
        $name = strtolower($rawName);
        if (isset($nameMap[$name])) {
            $name = $nameMap[$name];
        }
        return $name;
    }

    /**
     * @return ParserInterface message parser instance.
     * @throws Exception if unable to detect parser.
     * @throws \yii\base\InvalidConfigException
     */
    private function getParser()
    {
        $format = $this->getFormat();
        return $this->client->getParser($format);
    }
}
