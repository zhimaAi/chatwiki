// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"bufio"
	"encoding/json"
	"net/http"
)

type StreamReader[T any] struct {
	EmptyMessagesLimit uint
	IsFinished         bool
	Reader             *bufio.Reader
	Response           *http.Response
	ErrorResponse      ErrorResponseInterface
	ErrAccumulator     ErrorAccumulator
	HttpHeader         http.Header
}

func (stream *StreamReader[T]) UnmarshalError() {
	errBytes := stream.ErrAccumulator.Bytes()
	if len(errBytes) == 0 {
		return
	}

	err := json.Unmarshal(errBytes, stream.ErrorResponse)
	if err != nil {
		stream.ErrorResponse = nil
	}
}

func (stream *StreamReader[T]) Close() error {
	return stream.Response.Body.Close()
}
