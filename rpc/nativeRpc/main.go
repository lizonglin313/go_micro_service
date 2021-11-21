// Package main是服务端.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"rpc/nativeRpc/service"
)

func main()  {
	stringService := new(service.StringService)
	fmt.Println("Register function...")
	rpc.Register(stringService)
	fmt.Println("Function registered, wait request...")
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("listen error: e")
	}
	http.Serve(l, nil)
}