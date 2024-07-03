// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package moonshot

import (
	"errors"
	"fmt"
)

type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
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
		return errors.New(fmt.Sprintf("Moonshot request error, status code: %d, message: %s", e.HTTPStatusCode, e.Err.Message))
	}
	return errors.New(fmt.Sprintf("Moonshot request error, message: %s", e.Err.Message))
}
