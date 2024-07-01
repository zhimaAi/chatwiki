// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package baidu

import (
	"errors"
	"fmt"
)

type ErrorResponse struct {
	ErrorCode      int    `json:"error_code"`
	ErrorMessage   string `json:"error_msg"`
	HTTPStatusCode int    `json:"-"`
}

func (e *ErrorResponse) SetHTTPStatusCode(statusCode int) {
	e.HTTPStatusCode = statusCode
}

func (e *ErrorResponse) Error() error {
	if e.HTTPStatusCode > 0 {
		return errors.New(fmt.Sprintf("Baidu request error, status code: %d, message: %s", e.HTTPStatusCode, e.ErrorMessage))
	}
	return errors.New(fmt.Sprintf("Baidu request error, message: %s", e.ErrorMessage))
}
