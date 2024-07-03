// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package ali

import (
	"errors"
	"fmt"
)

type ErrorResponse struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	RequestId      string `json:"requestId"`
	HTTPStatusCode int    `json:"-"`
}

func (e *ErrorResponse) SetHTTPStatusCode(statusCode int) {
	e.HTTPStatusCode = statusCode
}

func (e *ErrorResponse) Error() error {
	if e.HTTPStatusCode > 0 {
		return errors.New(fmt.Sprintf("ALI request error, status code: %d, business code: %s, request_id: %s, message: %s", e.HTTPStatusCode, e.Code, e.RequestId, e.Message))
	}
	return errors.New(fmt.Sprintf("ALI request error, message: %s", e.Message))
}
