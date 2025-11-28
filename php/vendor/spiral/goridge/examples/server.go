package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	goridgeRpc "github.com/spiral/goridge/v3/pkg/rpc"
)

// App sample
type App struct{}

// Hi returns greeting message.
func (a *App) Hi(name string, r *string) error {
	*r = fmt.Sprintf("Hello, %s!", name)
	return nil
}

func main() {
	ln, err := net.Listen("tcp", ":6001")
	if err != nil {
		panic(err)
	}

	err = rpc.Register(new(App))
	if err != nil {
		panic(err)
	}
	log.Printf("started")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		log.Printf("new connection %+v", conn)
		go rpc.ServeCodec(goridgeRpc.NewCodec(conn))
	}
}
