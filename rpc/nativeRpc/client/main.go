package main

import (
	"fmt"
	"log"
	"net/rpc"
	"rpc/nativeRpc/service"
	"time"
)

func main() {
	// DialHTTP connects to an HTTP RPC server
	// at the specified network address listening on the default HTTP RPC path
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	stringReq := &service.StringRequest{"A", "B"}
	var reply string
	// 同步调用拼接字符串
	err = client.Call("StringService.Concat", stringReq, &reply)
	if err != nil {
		log.Fatal("Concat error:", err)
	}
	fmt.Printf("Reply form StringService is: %s\n", reply)

	stringReq = &service.StringRequest{"ACD", "BDF"}
	// 异步调用重复子串
	call := client.Go("StringService.Diff", stringReq, &reply, nil)
	callDone := false
	for !callDone {
		select {
		case <-call.Done:
			fmt.Println("异步执行完成！")
			callDone = true
		default:
			fmt.Println("Do other things...")
			time.Sleep(1 * time.Second)
		}
	}

	if err != nil {
		log.Fatal("Concat error:", err)
	}

	fmt.Printf("Reply form  is: %s\n", reply)
}
