package main

import (
	"context"
	"flag"
	"fmt"
	transport_grpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"log"
	"rpc/support_gRPC_HTTP/endpoint"
	"rpc/support_gRPC_HTTP/proto"
	"rpc/support_gRPC_HTTP/service"
	"rpc/support_gRPC_HTTP/transport"
	"time"
)

// NewTransportGRPCClient 从这里我们也知道 Endpoints 也实现 SumConcatService.
func NewTransportGRPCClient(conn *grpc.ClientConn) service.SumConcatService {
	sumEp := transport_grpc.NewClient(
		conn,
		"proto.SumConcat",
		"Sum",
		transport.EncodeGRPCSumRequest,  // 将客户端的 EP 请求编码成 PB 形式发送
		transport.DecodeGRPCSumResponse, // 将服务端的 PB 响应编码成 EP 形式接收
		proto.SumResponse{},
	).Endpoint()

	concatEp := transport_grpc.NewClient(
		conn,
		"proto.SumConcat",
		"Concat",
		transport.EncodeGRPCConcatRequest,
		transport.DecodeGRPCConcatResponse,
		proto.ConcatResponse{},
	).Endpoint()

	return endpoint.Endpoints{
		SumEndpoint:    sumEp,
		ConcatEndpoint: concatEp,
	}
}

func main() {

	gRPCAddr := flag.String("gRPC", ":9123", "gRPC client")
	flag.Parse()

	conn, err := grpc.Dial(
		*gRPCAddr,
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Second),
	)
	if err != nil {
		log.Fatalln("gRPC dial error:", err)
	}
	defer conn.Close()

	SumConcatService := NewTransportGRPCClient(conn)
	sum(context.Background(), SumConcatService, 111, 234)
	concat(context.Background(), SumConcatService, "111", "234")

}

func sum(ctx context.Context, svc service.SumConcatService, a, b int) {
	reply := svc.Sum(ctx, a, b)
	fmt.Println("Result form sum service is:", reply)
}

func concat(ctx context.Context, svc service.SumConcatService, a, b string) {
	reply := svc.Concat(ctx, a, b)
	fmt.Println("Result form concat service is:", reply)
}
