// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Header struct {
	Key   string
	Value string
}

type Param struct {
	Key   string
	Value string
}

func HttpPost(url string, headers []Header, params []Param, body any) (response *http.Response, err error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	q := req.URL.Query()
	for _, param := range params {
		q.Add(param.Key, param.Value)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")
	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	response = resp
	return
}

func HttpStreamPost(url string, headers []Header, params []Param, request any) (response *http.Response, err error) {
	newHeaders := append(headers,
		Header{Key: "Accept", Value: "text/event-stream"},
		Header{Key: "Cache-Control", Value: "no-cache"},
		Header{Key: "Connection", Value: "keep-alive"},
	)

	return HttpPost(url, newHeaders, params, request)
}

type ErrorResponseInterface interface {
	SetHTTPStatusCode(statusCode int)
	Error() error
}

func HttpCheckError(resp *http.Response, errorResp ErrorResponseInterface) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		err := json.NewDecoder(resp.Body).Decode(&errorResp)
		if err != nil {
			parseError := &ParseError{
				HTTPStatusCode: resp.StatusCode,
				Err:            err,
			}
			return parseError
		}
		errorResp.SetHTTPStatusCode(resp.StatusCode)
		return errorResp.Error()
	}
	return nil
}

func HttpCheckErrors[T ErrorResponseInterface](resp *http.Response) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errorResp T
		var errorRespArray []T
		err := json.NewDecoder(resp.Body).Decode(&errorRespArray)
		if err != nil {
			parseError := &ParseError{
				HTTPStatusCode: resp.StatusCode,
				Err:            err,
			}
			return parseError
		}
		if len(errorRespArray) > 0 {
			errorResp = errorRespArray[0]
			errorResp.SetHTTPStatusCode(resp.StatusCode)
			return errorResp.Error()
		}
		return errors.New(`parse error`)
	}
	return nil
}

func HttpDecodeResponse(resp *http.Response, v any) (err error) {

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, v)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return
}
