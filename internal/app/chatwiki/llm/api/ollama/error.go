// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package ollama

import (
	"errors"
	"fmt"
	"net/http"
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
	if e.Err == nil {
		e.Err = &APIError{
			Code:    fmt.Sprintf("%v", e.HTTPStatusCode),
			Message: fmt.Sprintf("%s", http.StatusText(e.HTTPStatusCode)),
		}
	}
	if e.HTTPStatusCode > 0 {
		return errors.New(fmt.Sprintf("Ollama request error, status code: %d, business code: %s, message: %s", e.HTTPStatusCode, e.Err.Code, e.Err.Message))
	}
	return errors.New(fmt.Sprintf("Ollama request error, message: %s", e.Err.Message))
}
