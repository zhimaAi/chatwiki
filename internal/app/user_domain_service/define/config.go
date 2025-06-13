// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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