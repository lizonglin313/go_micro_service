package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"rpc/gRPC/pb"
)

func main() {
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()

	stringCLient := pb.NewStringServiceClient(conn)
	stringReq := &pb.StringRequest{A: "ADF", B: "BcD"}
	reply, _ := stringCLient.Diff(context.Background(), stringReq)
	fmt.Printf("reply from Concat is: %s\n", reply.Ret)
}
