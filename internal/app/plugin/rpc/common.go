// Copyright © 2016- 2025 Sesame Network Technology all right reserved

package rpc

import "fmt"

type Common struct{}

func (s *Common) Hello(name string, r *string) error {
	*r = fmt.Sprintf("Hello, %s!", name)
	return nil
}
