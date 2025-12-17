// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

var Config ConfigParam

type ConfigParam struct {
	WebService map[string]string
	NumCPU     map[string]string
	Redis      map[string]string
	Postgres   map[string]string
	NsqLookup  map[string]string
	Nsqd       map[string]string
}
