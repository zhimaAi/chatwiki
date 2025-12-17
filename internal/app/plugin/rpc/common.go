// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package rpc

import "fmt"

type Common struct{}

func (s *Common) Hello(name string, r *string) error {
	*r = fmt.Sprintf("Hello, %s!", name)
	return nil
}
