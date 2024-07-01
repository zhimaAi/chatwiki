// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

var Config ConfigParam

type ConfigParam struct {
	NumCPU    map[string]string
	NsqLookup map[string]string
	Nsqd      map[string]string
}
