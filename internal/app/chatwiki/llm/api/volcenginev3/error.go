// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package volcenginev3

import (
	"errors"
	"fmt"
)

type ErrorResponse struct {
	HTTPStatusCode int    `json:"-"`
	Code           string `json:"code"`
	ReqID          string `json:"req_id"`
}

func (e *ErrorResponse) SetHTTPStatusCode(statusCode int) {
	e.HTTPStatusCode = statusCode
}

func (e *ErrorResponse) Error() error {
	if e.HTTPStatusCode > 0 {
		return errors.New(fmt.Sprintf("Volcengine request error, status code: %d, business code: %s, req_id: %s", e.HTTPStatusCode, e.Code, e.ReqID))
	}
	return errors.New(fmt.Sprintf("Volcengine request error, business code: %s, req_id: %s", e.Code, e.ReqID))
}
