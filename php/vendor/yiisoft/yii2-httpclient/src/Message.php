<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient;

use yii\base\Component;
use yii\base\ErrorHandler;
use yii\web\Cookie;
use yii\web\CookieCollection;
use yii\web\HeaderCollection;
use Yii;

/**
 * Message represents a base HTTP message.
 *
 * @property string $content Raw body.
 * @property CookieCollection|Cookie[] $cookies The cookie collection. Note that the type of this property
 * differs in getter and setter. See [[getCookies()]] and [[setCookies()]] for details.
 * @property mixed $data Content data fields.
 * @property string $format Body format name.
 * @property HeaderCollection $headers The header collection. Note that the type of this property differs in
 * getter and setter. See [[getHeaders()]] and [[setHeaders()]] for details.
 *
 * @author Paul Klimov <klimov.paul@gmail.com>
 * @since 2.0
 */
class Message extends Component
{
    /**
     * @var Client owner client instance.
     */
    public $client;

    /**
     * @var HeaderCollection headers.
     */
    private $_headers;
    /**
     * @var CookieCollection cookies.
     */
    private $_cookies;
    /**
     * @var string|null raw content
     */
    private $_content;
    /**
     * @var mixed content data
     */
    private $_data;
    /**
     * @var string content format name
     */
    private $_format;


    /**
     * Sets the HTTP headers associated with HTTP message.
     * @param array|HeaderCollection $headers headers collection or headers list in format: [headerName => headerValue]
     * @return $this self reference.
     */
    public function setHeaders($headers)
    {
        $this->_headers = $headers;
        return $this;
    }

    /**
     * Returns the header collection.
     * The header collection contains the HTTP headers associated with HTTP message.
     * @return HeaderCollection the header collection
     */
    public function getHeaders()
    {
        if (!is_object($this->_headers)) {
            $headerCollection = new HeaderCollection();
            if (is_array($this->_headers)) {
                foreach ($this->_headers as $name => $value) {
                    if (is_int($name)) {
                        // parse raw header :
                        $rawHeader = $value;
                        if (strpos($rawHeader, 'HTTP/') === 0) {
                            $parts = explode(' ', $rawHeader, 3);
                            $headerCollection->add('http-code', $parts[1]);
                        } elseif (($separatorPos = strpos($rawHeader, ':')) !== false) {
                            $name = strtolower(trim(substr($rawHeader, 0, $separatorPos)));
                            $value = trim(substr($rawHeader, $separatorPos + 1));
                            $headerCollection->add($name, $value);
                        } else {
                            $headerCollection->add('raw', $rawHeader);
                        }
                    } else {
                        $headerCollection->set($name, $value);
                    }
                }
            }
            $this->_headers = $headerCollection;
        }
        return $this->_headers;
    }

    /**
     * Adds more headers to the already defined ones.
     * @param array $headers additional headers in format: [headerName => headerValue]
     * @return $this self reference.
     */
    public function addHeaders(array $headers)
    {
        $headerCollection = $this->getHeaders();
        foreach ($headers as $name => $value) {
            $headerCollection->add($name, $value);
        }
        return $this;
    }

    /**
     * Checks of HTTP message contains any header.
     * Using this method you are able to check cookie presence without instantiating [[HeaderCollection]].
     * @return bool whether message contains any header.
     */
    public function hasHeaders()
    {
        if (is_object($this->_headers)) {
            return $this->_headers->getCount() > 0;
        }
        return !empty($this->_headers);
    }

    /**
     * Sets the cookies associated with HTTP message.
     * @param CookieCollection|Cookie[]|array $cookies cookie collection or cookies list.
     * @return $this self reference.
     */
    public function setCookies($cookies)
    {
        $this->_cookies = $cookies;
        return $this;
    }

    /**
     * Returns the cookie collection.
     * The cookie collection contains the cookies associated with HTTP message.
     * @return CookieCollection|Cookie[] the cookie collection.
     */
    public function getCookies()
    {
        if (!is_object($this->_cookies)) {
            $cookieCollection = new CookieCollection();
            if (is_array($this->_cookies)) {
                foreach ($this->_cookies as $cookie) {
                    if (!is_object($cookie)) {
                        $cookie = new Cookie($cookie);
                    }
                    $cookieCollection->add($cookie);
                }
            }
            $this->_cookies = $cookieCollection;
        }
        return $this->_cookies;
    }

    /**
     * Adds more cookies to the already defined ones.
     * @param Cookie[]|array $cookies additional cookies.
     * @return $this self reference.
     */
    public function addCookies(array $cookies)
    {
        $cookieCollection = $this->getCookies();
        foreach ($cookies as $cookie) {
            if (!is_object($cookie)) {
                $cookie = new Cookie($cookie);
            }
            $cookieCollection->add($cookie);
        }
        return $this;
    }

    /**
     * Checks of HTTP message contains any cookie.
     * Using this method you are able to check cookie presence without instantiating [[CookieCollection]].
     * @return bool whether message contains any cookie.
     */
    public function hasCookies()
    {
        if (is_object($this->_cookies)) {
            return $this->_cookies->getCount() > 0;
        }
        return !empty($this->_cookies);
    }

    /**
     * Sets the HTTP message raw content.
     * @param string $content raw content.
     * @return $this self reference.
     */
    public function setContent($content)
    {
        $this->_content = $content;
        return $this;
    }

    /**
     * Returns HTTP message raw content.
     * @return string raw body.
     */
    public function getContent()
    {
        return $this->_content;
    }

    /**
     * Checks if content with provided name exists
     * @param $key string Name of the content parameter
     * @return bool
     * @since 2.0.10
     */
    public function hasContent($key)
    {
        $content = $this->getContent();
        return is_array($content) && isset($content[$key]);
    }



    /**
     * Sets the data fields, which composes message content.
     * @param mixed $data content data fields.
     * @return $this self reference.
     */
    public function setData($data)
    {
        $this->_data = $data;
        return $this;
    }

    /**
     * Returns the data fields, parsed from raw content.
     * @return mixed content data fields.
     */
    public function getData()
    {
        return $this->_data;
    }

    /**
     * Adds data fields to the existing ones.
     * @param array $data additional content data fields.
     * @return $this self reference.
     * @since 2.0.1
     */
    public function addData($data)
    {
        if (empty($this->_data)) {
            $this->_data = $data;
        } else {
            if (!is_array($this->_data)) {
                throw new \yii\base\Exception('Unable to merge existing data with new data. Existing data is not an array.');
            }
            $this->_data = array_merge($this->_data, $data);
        }
        return $this;
    }

    /**
     * Sets body format.
     * @param string $format body format name.
     * @return $this self reference.
     */
    public function setFormat($format)
    {
        $this->_format = $format;
        return $this;
    }

    /**
     * Returns body format.
     * @return string body format name.
     */
    public function getFormat()
    {
        if ($this->_format === null) {
            $this->_format = $this->defaultFormat();
        }
        return $this->_format;
    }

    /**
     * Returns default format name.
     * @return string default format name.
     */
    protected function defaultFormat()
    {
        return Client::FORMAT_URLENCODED;
    }

    /**
     * Composes raw header lines from [[headers]].
     * Each line will be a string in format: 'header-name: value'.
     * @return array raw header lines.
     */
    public function composeHeaderLines()
    {
        if (!$this->hasHeaders()) {
            return [];
        }
        $headers = [];
        foreach ($this->getHeaders() as $name => $values) {
            $name = str_replace(' ', '-', ucwords(str_replace('-', ' ', $name)));
            foreach ($values as $value) {
                $headers[] = "$name: $value";
            }
        }
        return $headers;
    }

    /**
     * Returns string representation of this HTTP message.
     * @return string the string representation of this HTTP message.
     */
    public function toString()
    {
        $result = '';
        if ($this->hasHeaders()) {
            $headers = $this->composeHeaderLines();
            $result .= implode("\n", $headers);
        }

        $content = $this->getContent();
        if ($content !== null) {
            $result .= "\n\n" . $content;
        }

        return $result;
    }

    /**
     * PHP magic method that returns the string representation of this object.
     * @return string the string representation of this object.
     */
    public function __toString()
    {
        // __toString cannot throw exception
        // use trigger_error to bypass this limitation
        try {
            return $this->toString();
        } catch (\Exception $e) {
            ErrorHandler::convertExceptionToError($e);
            return '';
        }
    }
}