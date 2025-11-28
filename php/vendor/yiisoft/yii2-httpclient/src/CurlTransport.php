<?php
/**
 * @link https://www.yiiframework.com/
 * @copyright Copyright (c) 2008 Yii Software LLC
 * @license https://www.yiiframework.com/license/
 */

namespace yii\httpclient;

use Yii;

/**
 * CurlTransport sends HTTP messages using [Client URL Library (cURL)](https://php.net/manual/en/book.curl.php)
 *
 * Note: this transport requires PHP 'curl' extension installed.
 *
 * For this transport, you may setup request options as [cURL Options](https://php.net/manual/en/function.curl-setopt.php)
 *
 * @author Paul Klimov <klimov.paul@gmail.com>
 * @since 2.0
 */
class CurlTransport extends Transport
{
    /**
     * {@inheritdoc}
     */
    public function send($request)
    {
        $request->beforeSend();

        $curlOptions = $this->prepare($request);
        $curlResource = $this->initCurl($curlOptions);

        $responseHeaders = [];
        $this->setHeaderOutput($curlResource, $responseHeaders);

        $token = $request->client->createRequestLogToken($request->getMethod(), $curlOptions[CURLOPT_URL], $curlOptions[CURLOPT_HTTPHEADER], print_r($request->getContent(), true));
        Yii::info($token, __METHOD__);
        Yii::beginProfile($token, __METHOD__);
        $responseContent = curl_exec($curlResource);
        Yii::endProfile($token, __METHOD__);

        // check cURL error
        $errorNumber = curl_errno($curlResource);
        $errorMessage = curl_error($curlResource);

        curl_close($curlResource);

        if ($errorNumber > 0) {
            throw new Exception('Curl error: #' . $errorNumber . ' - ' . $errorMessage, $errorNumber);
        }

        $response = $request->client->createResponse($responseContent, $responseHeaders);

        $request->afterSend($response);

        return $response;
    }

    /**
     * {@inheritdoc}
     */
    public function batchSend(array $requests)
    {
        $curlBatchResource = curl_multi_init();

        $token = '';
        $curlResources = [];
        $responseHeaders = [];
        foreach ($requests as $key => $request) {
            /* @var $request Request */
            $request->beforeSend();

            $curlOptions = $this->prepare($request);
            $curlResource = $this->initCurl($curlOptions);

            $token .= $request->client->createRequestLogToken($request->getMethod(), $curlOptions[CURLOPT_URL], $curlOptions[CURLOPT_HTTPHEADER], $request->getContent()) . "\n\n";

            $responseHeaders[$key] = [];
            $this->setHeaderOutput($curlResource, $responseHeaders[$key]);
            $curlResources[$key] = $curlResource;
            curl_multi_add_handle($curlBatchResource, $curlResource);
        }

        Yii::info($token, __METHOD__);
        Yii::beginProfile($token, __METHOD__);

        try {
            $isRunning = null;
            do {
                // See https://bugs.php.net/bug.php?id=61141
                if (curl_multi_select($curlBatchResource) === -1) {
                    usleep(100);
                }
                do {
                    $curlExecCode = curl_multi_exec($curlBatchResource, $isRunning);
                } while ($curlExecCode === CURLM_CALL_MULTI_PERFORM);
            } while ($isRunning > 0 && $curlExecCode === CURLM_OK);
        } catch (\Exception $e) {
            Yii::endProfile($token, __METHOD__);
            throw new Exception($e->getMessage(), $e->getCode(), $e);
        }

        Yii::endProfile($token, __METHOD__);

        $responseContents = [];
        foreach ($curlResources as $key => $curlResource) {
            $responseContents[$key] = curl_multi_getcontent($curlResource);
            curl_multi_remove_handle($curlBatchResource, $curlResource);
        }

        curl_multi_close($curlBatchResource);

        $responses = [];
        foreach ($requests as $key => $request) {
            $response = $request->client->createResponse($responseContents[$key], $responseHeaders[$key]);
            $request->afterSend($response);
            $responses[$key] = $response;
        }
        return $responses;
    }

    /**
     * Prepare request for execution, creating cURL resource for it.
     * @param Request $request request instance.
     * @return array cURL options.
     */
    private function prepare($request)
    {
        $request->prepare();

        $curlOptions = $this->composeCurlOptions($request->getOptions());

        $method = strtoupper($request->getMethod());
        switch ($method) {
            case 'POST':
                $curlOptions[CURLOPT_POST] = true;
                break;
            default:
                $curlOptions[CURLOPT_CUSTOMREQUEST] = $method;
        }

        $content = $request->getContent();

        if ($method === 'HEAD') {
            $curlOptions[CURLOPT_NOBODY] = true;
        }

        if ($content !== null) {
            $curlOptions[CURLOPT_POSTFIELDS] = $content;
        }

        $curlOptions[CURLOPT_RETURNTRANSFER] = true;
        $curlOptions[CURLOPT_URL] = $request->getFullUrl();
        $curlOptions[CURLOPT_HTTPHEADER] = $request->composeHeaderLines();

        if ($request->getOutputFile()) {
            $curlOptions[CURLOPT_FILE] = $request->getOutputFile();
        }

        return $curlOptions;
    }

    /**
     * Initializes cURL resource.
     * @param array $curlOptions cURL options.
     * @return resource prepared cURL resource.
     */
    private function initCurl(array $curlOptions)
    {
        $curlResource = curl_init();
        foreach ($curlOptions as $option => $value) {
            curl_setopt($curlResource, $option, $value);
        }

        return $curlResource;
    }

    /**
     * Composes cURL options from raw request options.
     * @param array $options raw request options.
     * @return array cURL options, in format: [curl_constant => value].
     */
    private function composeCurlOptions(array $options)
    {
        static $optionMap = [
            'protocolVersion' => CURLOPT_HTTP_VERSION,
            'maxRedirects' => CURLOPT_MAXREDIRS,
            'sslCapath' => CURLOPT_CAPATH,
            'sslCafile' => CURLOPT_CAINFO,
            'sslLocalCert' => CURLOPT_SSLCERT,
            'sslLocalPk' => CURLOPT_SSLKEY,
            'sslPassphrase' => CURLOPT_SSLCERTPASSWD,
        ];

        $curlOptions = [];
        foreach ($options as $key => $value) {
            if (is_int($key)) {
                $curlOptions[$key] = $value;
            } else {
                if (isset($optionMap[$key])) {
                    $curlOptions[$optionMap[$key]] = $value;
                } else {
                    $key = strtoupper($key);
                    if (strpos($key, 'SSL') === 0) {
                        $key = substr($key, 3);
                        $constantName = 'CURLOPT_SSL_' . $key;
                        if (!defined($constantName)) {
                            $constantName = 'CURLOPT_SSL' . $key;
                        }
                    } else {
                        $constantName = 'CURLOPT_' . strtoupper($key);
                    }
                    $curlOptions[constant($constantName)] = $value;
                }
            }
        }
        return $curlOptions;
    }

    /**
     * Setup a variable, which should collect the cURL response headers.
     * @param resource $curlResource cURL resource.
     * @param array $output variable, which should collection headers.
     */
    private function setHeaderOutput($curlResource, array &$output)
    {
        curl_setopt($curlResource, CURLOPT_HEADERFUNCTION, function($resource, $headerString) use (&$output) {
            $header = trim($headerString, "\n\r");
            if (strlen($header) > 0) {
                $output[] = $header;
            }
            return mb_strlen($headerString, '8bit');
        });
    }
}
