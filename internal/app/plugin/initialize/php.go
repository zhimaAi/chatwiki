// Copyright © 2016- 2025 Sesame Network Technology all right reserved

package initialize

import (
	"chatwiki/internal/app/plugin/define"
	"chatwiki/internal/app/plugin/php"
	"chatwiki/internal/app/plugin/rpc"
	"fmt"
	"net"
	netRpc "net/rpc"
	"sync/atomic"

	goridgeRpc "github.com/roadrunner-server/goridge/v3/pkg/rpc"
	"github.com/zhimaAi/go_tools/logs"
)

var listener net.Listener
var closed uint32

func StartPhpGoridge() {
	var err error
	atomic.StoreUint32(&closed, 0)

	port := define.Config.RpcService[`port`]
	listener, err = net.Listen("tcp", ":"+port)
	if err != nil {
		panic(fmt.Sprintf("failed to start php_goridge RPC server: %s", err.Error()))
	}
	server := netRpc.NewServer()
	if err := server.Register(new(rpc.Common)); err != nil {
		panic(fmt.Sprintf("failed to register RPC service: %s", err.Error()))
	}
	if err := server.Register(new(rpc.AppLogger)); err != nil {
		panic(fmt.Sprintf("failed to register RPC service: %s", err.Error()))
	}
	logs.Info("php_goridge RPC server initialized, port = %s", port)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if atomic.LoadUint32(&closed) == 1 {
					return
				}
				logs.Error("failed to accept connection: %v", err.Error())
				continue
			}
			logs.Info("new connection on address %s", conn.RemoteAddr().String())
			go server.ServeCodec(goridgeRpc.NewCodec(conn))
		}
	}()
}

func StopPhpGoridge() {
	atomic.StoreUint32(&closed, 1)
	if err := listener.Close(); err != nil {
		logs.Error("failed to close listener: %v", err.Error())
	}
	logs.Info("php_goridge RPC server stopped")
}

func initPhpPluginPool() {
	define.PhpPlugin = php.NewPhpPluginPool("php/worker.php")
}
