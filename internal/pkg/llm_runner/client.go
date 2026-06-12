// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package llm_runner

import (
	"net/rpc"
	"strings"
)

func RpcExecuteRun(address, command string) (resp RpcRunResponse) {
	if !strings.Contains(address, `:`) {
		address = address + RpcServicePort
	}
	client, err := rpc.Dial(`tcp`, address)
	if err != nil {
		resp.IsError = true
		resp.ErrorMsg = err.Error()
		return
	}
	defer func(client *rpc.Client) {
		_ = client.Close()
	}(client)
	err = client.Call(RpcServiceName+`.Run`, RpcRunRequest{Command: command}, &resp)
	if err != nil {
		resp.IsError = true
		resp.ErrorMsg = err.Error()
		return
	}
	return
}
