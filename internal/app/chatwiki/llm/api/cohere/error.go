// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package cohere

import (
	"errors"
	"fmt"
)

type ErrorResponse struct {
	HTTPStatusCode int    `json:"-"`
	Data           string `json:"data"`
	Message        string `json:"message"`
}

func (e *ErrorResponse) SetHTTPStatusCode(statusCode int) {
	e.HTTPStatusCode = statusCode
}

func (e *ErrorResponse) Error() error {
	if e.HTTPStatusCode > 0 {
		return errors.New(fmt.Sprintf("Cohere request error, status code: %d, message: %s", e.HTTPStatusCode, e.Message))
	}
	return errors.New(fmt.Sprintf("Cohere request error, message: %s", e.Message))
}
