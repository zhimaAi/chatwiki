// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package rpc

import (
	"github.com/zhimaAi/go_tools/logs"
)

type AppLogger struct{}

func (s *AppLogger) Info(in string, _ *bool) error {
	logs.Info(in)

	return nil
}

func (s *AppLogger) Debug(in string, _ *bool) error {
	logs.Debug(in)

	return nil
}

func (s *AppLogger) Error(in string, _ *bool) error {
	logs.Error(in)

	return nil
}

func (s *AppLogger) Trace(in string, _ *bool) error {
	logs.Debug(in)

	return nil
}

func (s *AppLogger) Warning(in string, _ *bool) error {
	logs.Warning(in)

	return nil
}
