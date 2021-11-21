package main

import (
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "rpc/stream_gRPC/pb"
	"rpc/stream_gRPC/string_service"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// gRPC 可以很方便的注册函数
	grpcServer := grpc.NewServer()
	// 因为 string_service.StringService 实现了 pb.go 里面的 StringServiceServer 接口
	// 所以是符合 gRPC 规范的
	stringService := new(string_service.StringService)
	pb.RegisterStringServiceServer(grpcServer, stringService)
	grpcServer.Serve(lis)
}
