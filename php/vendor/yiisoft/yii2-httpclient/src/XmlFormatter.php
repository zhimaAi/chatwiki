<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient;

use DOMDocument;
use DOMElement;
use DOMText;
use SimpleXMLElement;
use yii\base\Arrayable;
use yii\base\BaseObject;
use Yii;
use yii\helpers\StringHelper;

/**
 * XmlFormatter formats HTTP message as XML.
 *
 * @author Paul Klimov <klimov.paul@gmail.com>
 * @since 2.0
 */
class XmlFormatter extends BaseObject implements FormatterInterface
{
    /**
     * @var string the Content-Type header for the response
     */
    public $contentType = 'application/xml';
    /**
     * @var string the XML version
     */
    public $version = '1.0';
    /**
     * @var string the XML encoding. If not set, it will use the value of [[\yii\base\Application::charset]].
     */
    public $encoding;
    /**
     * @var string the name of the root element.
     */
    public $rootTag = 'request';
    /**
     * @var string the name of the elements that represent the array elements with numeric keys.
     * @since 2.0.1
     */
    public $itemTag = 'item';
    /**
     * @var bool whether to interpret objects implementing the [[\Traversable]] interface as arrays.
     * Defaults to `true`.
     * @since 2.0.1
     */
    public $useTraversableAsArray = true;


    /**
     * {@inheritdoc}
     */
    public function format(Request $request)
    {
        $contentType = $this->contentType;
        $charset = $this->encoding === null ? Yii::$app->charset : $this->encoding;
        if (stripos($contentType, 'charset') === false) {
            $contentType .= '; charset=' . $charset;
        }
        $request->getHeaders()->set('Content-Type', $contentType);

        $data = $request->getData();
        if ($data !== null) {
            if ($data instanceof DOMDocument) {
                $content = $data->saveXML();
            } elseif ($data instanceof SimpleXMLElement) {
                $content = $data->saveXML();
            } else {
                $dom = new DOMDocument($this->version, $charset);
                $root = new DOMElement($this->rootTag);
                $dom->appendChild($root);
                $this->buildXml($root, $data);
                $content = $dom->saveXML();
            }
            $request->setContent($content);
        }

        return $request;
    }

    /**
     * @param DOMElement $element
     * @param mixed $data
     */
    protected function buildXml($element, $data)
    {
        if (is_array($data) ||
            ($data instanceof \Traversable && $this->useTraversableAsArray && !$data instanceof Arrayable)
        ) {
            foreach ($data as $name => $value) {
                if (is_int($name) && is_object($value)) {
                    $this->buildXml($element, $value);
                } elseif (is_array($value) || is_object($value)) {
                    $child = new DOMElement(is_int($name) ? $this->itemTag : $name);
                    $element->appendChild($child);
                    $this->buildXml($child, $value);
                } else {
                    $child = new DOMElement(is_int($name) ? $this->itemTag : $name);
                    $element->appendChild($child);
                    $child->appendChild(new DOMText((string) $value));
                }
            }
        } elseif (is_object($data)) {
            $child = new DOMElement(StringHelper::basename(get_class($data)));
            $element->appendChild($child);
            if ($data instanceof Arrayable) {
                $this->buildXml($child, $data->toArray());
            } else {
                $array = [];
                foreach ($data as $name => $value) {
                    $array[$name] = $value;
                }
                $this->buildXml($child, $array);
            }
        } else {
            $element->appendChild(new DOMText((string) $data));
        }
    }
}