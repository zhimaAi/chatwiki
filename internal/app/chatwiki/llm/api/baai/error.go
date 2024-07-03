// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package baai

import (
	"errors"
	"fmt"
)

type ErrorResponse struct {
	Type           string      `json:"type"`
	Err            ErrorDetail `json:"error"`
	HTTPStatusCode int         `json:"-"`
}
type ErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e *ErrorResponse) SetHTTPStatusCode(statusCode int) {
	e.HTTPStatusCode = statusCode
}

func (e *ErrorResponse) Error() error {
	if e.HTTPStatusCode > 0 {
		return errors.New(fmt.Sprintf("Bge-service request error, status code: %d, type: %s, message: %s", e.HTTPStatusCode, e.Err.Type, e.Err.Message))
	}
	return errors.New(fmt.Sprintf("Bge-service request error, message: %s", e.Err.Message))
}
