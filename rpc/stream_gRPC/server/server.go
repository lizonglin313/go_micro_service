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
	grpcServer := grpc.NewServer()
	stringService := new(string_service.StringService)
	pb.RegisterStringServiceServer(grpcServer, stringService)
	grpcServer.Serve(lis)
}
