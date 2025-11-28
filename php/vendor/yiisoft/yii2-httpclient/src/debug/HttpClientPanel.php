<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient\debug;

use yii\debug\Panel;
use yii\di\Instance;
use yii\httpclient\Client;
use yii\log\Logger;
use Yii;

/**
 * Debugger panel that collects and displays HTTP requests performed.
 *
 * @property \yii\httpclient\Client $httpClient Note that the type of this property differs in getter and
 * setter. See [[getHttpClient()]] and [[setHttpClient()]] for details.
 * @property array $methods This property is read-only.
 * @property array $types This property is read-only.
 *
 * @author Paul Klimov <klimov.paul@gmail.com>
 * @since 2.0
 */
class HttpClientPanel extends Panel
{
    /**
     * @var array current HTTP request timings
     */
    private $_timings;
    /**
     * @var array HTTP requests info extracted to array as models, to use with data provider.
     */
    private $_models;
    /**
     * @var \yii\httpclient\Client|array|string
     */
    private $_httpClient = 'yii\httpclient\Client';


    /**
     * @param array $httpClient
     */
    public function setHttpClient($httpClient)
    {
        $this->_httpClient = $httpClient;
    }

    /**
     * @return \yii\httpclient\Client
     * @throws \yii\base\InvalidConfigException
     */
    public function getHttpClient()
    {
        if (!is_object($this->_httpClient)) {
            $this->_httpClient = Instance::ensure($this->_httpClient, Client::className());
        }
        return $this->_httpClient;
    }

    /**
     * {@inheritdoc}
     */
    public function init()
    {
        $this->actions['request-execute'] = [
            'class' => 'yii\httpclient\debug\RequestExecuteAction',
            'panel' => $this,
        ];
    }

    /**
     * {@inheritdoc}
     */
    public function getName()
    {
        return 'HTTP Client';
    }

    /**
     * {@inheritdoc}
     */
    public function getSummary()
    {
        $timings = $this->calculateTimings();
        $queryCount = count($timings);
        if ($queryCount === 0) {
            return '';
        }

        $queryTime = number_format($this->getTotalRequestTime($timings) * 1000) . ' ms';

        return Yii::$app->view->render('@yii/httpclient/debug/views/summary', [
            'timings' => $this->calculateTimings(),
            'panel' => $this,
            'queryCount' => $queryCount,
            'queryTime' => $queryTime,
        ]);
    }

    /**
     * {@inheritdoc}
     */
    public function getDetail()
    {
        $searchModel = new SearchModel();
        $dataProvider = $searchModel->search(Yii::$app->request->getQueryParams(), $this->getModels());

        return Yii::$app->view->render('@yii/httpclient/debug/views/detail', [
            'panel' => $this,
            'dataProvider' => $dataProvider,
            'searchModel' => $searchModel,
        ]);
    }

    /**
     * Calculates given request profile timings.
     *
     * @return array timings [token, category, timestamp, traces, nesting level, elapsed time]
     */
    public function calculateTimings()
    {
        if ($this->_timings === null) {
            $this->_timings = Yii::getLogger()->calculateTimings(isset($this->data['messages']) ? $this->data['messages'] : []);
        }

        return $this->_timings;
    }

    /**
     * {@inheritdoc}
     */
    public function save()
    {
        $target = $this->module->logTarget;
        $messages = $target->filterMessages($target->messages, Logger::LEVEL_PROFILE, [
            'yii\httpclient\Transport::*',
            'yii\httpclient\CurlTransport::*',
            'yii\httpclient\StreamTransport::*',
        ]);
        return ['messages' => $messages];
    }

    /**
     * Returns an  array of models that represents logs of the current request.
     * Can be used with data providers such as \yii\data\ArrayDataProvider.
     * @return array models
     */
    protected function getModels()
    {
        if ($this->_models === null) {
            $this->_models = [];
            $timings = $this->calculateTimings();

            foreach ($timings as $seq => $dbTiming) {
                $this->_models[] = [
                    'method' => $this->getRequestMethod($dbTiming['info']),
                    'type' => $this->getRequestType($dbTiming['category']),
                    'request' => $dbTiming['info'],
                    'duration' => ($dbTiming['duration'] * 1000), // in milliseconds
                    'trace' => $dbTiming['trace'],
                    'timestamp' => ($dbTiming['timestamp'] * 1000), // in milliseconds
                    'seq' => $seq,
                ];
            }
        }

        return $this->_models;
    }

    /**
     * Returns HTTP request method.
     *
     * @param string $timing timing procedure string
     * @return string request method such as GET, POST, PUT, etc.
     */
    protected function getRequestMethod($timing)
    {
        $timing = ltrim($timing);
        preg_match('/^([a-zA-z]*)/', $timing, $matches);

        return count($matches) ? $matches[0] : '';
    }

    /**
     * Returns request type.
     *
     * @param string $category
     * @return string request type such as 'normal', 'batch'
     */
    protected function getRequestType($category)
    {
        return (stripos($category, '::batchSend') === false) ? 'normal' : 'batch';
    }

    /**
     * Returns total request time.
     *
     * @param array $timings
     * @return int total time
     */
    protected function getTotalRequestTime($timings)
    {
        $queryTime = 0;

        foreach ($timings as $timing) {
            $queryTime += $timing['duration'];
        }

        return $queryTime;
    }

    /**
     * Returns array request methods
     *
     * @return array
     */
    public function getMethods()
    {
        return array_reduce(
            $this->_models,
            function ($result, $item) {
                $result[$item['method']] = $item['method'];
                return $result;
            },
            []
        );
    }

    /**
     * Returns array request types
     *
     * @return array
     */
    public function getTypes()
    {
        return [
            'normal' => 'Normal',
            'batch' => 'Batch',
        ];
    }
}