// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package gemini

import (
	"errors"
	"fmt"
)

type ErrorResponse struct {
	Err            Err `json:"error"`
	HTTPStatusCode int `json:"-"`
}

type Err struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Status  string   `json:"status"`
	Details []Detail `json:"details"`
}

type Detail struct {
	Type     string   `json:"@type"`
	Reason   string   `json:"reason"`
	Domain   string   `json:"domain"`
	Metadata Metadata `json:"metadata"`
}
type Metadata struct {
	Service string `json:"service"`
}

func (e *ErrorResponse) SetHTTPStatusCode(statusCode int) {
	e.HTTPStatusCode = statusCode
}

func (e *ErrorResponse) Error() error {
	if e.HTTPStatusCode > 0 {
		return errors.New(fmt.Sprintf("Gemini request error, status code: %d, status: %s, message: %s", e.HTTPStatusCode, e.Err.Status, e.Err.Message))
	}
	return errors.New(fmt.Sprintf("Gemini request error, message: %s", e.Err.Message))
}
