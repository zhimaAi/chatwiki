// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package azure

import (
	"errors"
	"fmt"
)

type APIError struct {
	Code    string  `json:"code,omitempty"`
	Message string  `json:"message"`
	Param   *string `json:"param,omitempty"`
	Type    string  `json:"type"`
}

type ErrorResponse struct {
	Err            *APIError `json:"error,omitempty"`
	HTTPStatusCode int       `json:"-"`
}

func (e *ErrorResponse) SetHTTPStatusCode(statusCode int) {
	e.HTTPStatusCode = statusCode
}

func (e *ErrorResponse) Error() error {
	if e.HTTPStatusCode > 0 {
		return errors.New(fmt.Sprintf("Azure request error, status code: %d, business code: %s, message: %s", e.HTTPStatusCode, e.Err.Code, e.Err.Message))
	}
	return errors.New(fmt.Sprintf("Azure request error, message: %s", e.Err.Message))
}
