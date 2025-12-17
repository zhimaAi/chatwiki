// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

var Config ConfigParam

type ConfigParam struct {
	WebService map[string]string
	RpcService map[string]string
	NumCPU     map[string]string
	Postgres   map[string]string
	Redis      map[string]string
	Xiaokefu   map[string]string
}
