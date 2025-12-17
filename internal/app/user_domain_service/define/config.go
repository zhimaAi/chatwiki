// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

var Config ConfigParam

type ConfigParam struct {
	WebService map[string]string
	NumCPU     map[string]string
	ChatWiki   map[string]string
}

const (
	RobotDomainLabel = 1
)
