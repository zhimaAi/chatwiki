// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package llm_runner

const RpcServiceName = `CommandRunner`
const RpcServicePort = `:9000`

type RpcRunRequest struct {
	Command string
}

type RpcRunResponse struct {
	Output   string // Command output content
	ExitCode int    // Command exit code
	IsError  bool
	ErrorMsg string
}
