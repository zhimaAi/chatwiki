// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package jina

import (
	"errors"
	"fmt"
)

type ErrorResponse struct {
	HTTPStatusCode int      `json:"-"`
	Detail         []Detail `json:"detail"`
}
type Detail struct {
	Loc  []string `json:"loc"`
	Msg  string   `json:"msg"`
	Type string   `json:"type"`
}

func (e *ErrorResponse) SetHTTPStatusCode(statusCode int) {
	e.HTTPStatusCode = statusCode
}

func (e *ErrorResponse) Error() error {
	if e.HTTPStatusCode > 0 {
		return errors.New(fmt.Sprintf("Cohere request error, status code: %d, detail: %s", e.HTTPStatusCode, e.Detail[0].Msg))
	}
	return errors.New(fmt.Sprintf("Cohere request error, detail: %s", e.Detail[0].Msg))
}
