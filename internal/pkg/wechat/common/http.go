// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
)

func HttpRead(resp *http.Response) ([]byte, error) {
	if resp == nil || resp.Body == nil {
		return nil, errors.New(`response body empty`)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		return io.ReadAll(reader)
	}
	return io.ReadAll(resp.Body)
}
