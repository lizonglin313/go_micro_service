package main

import (
	"context"
	"flag"
	"fmt"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"rpc/go_kit_gRPC/proto"
	"rpc/go_kit_gRPC/service"
	"time"
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "127.0.0.1:8123", grpc.WithInsecure())
	if err != nil {
		fmt.Println("gRPC dial err:", err)
	}
	defer conn.Close()

	svr := NewStringClient(conn)
	result, err := svr.Concat(ctx, "A", "B")
	if err != nil {
		fmt.Println("Check error:", err.Error())
	}

	fmt.Println("result=", result)
}

func NewStringClient(conn *grpc.ClientConn) service.Service {
	var ep = grpctransport.NewClient(conn,
		"proto.StringService",
		"Concat",
		DecodeStringRequest,
		EncodeStringResponse,
		proto.StringResponse{},
	).Endpoint()

	userEp := service.StringEndpoints{
		StringEndpoint: ep,
	}
	return userEp
}

func DecodeStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return r, nil
}

func EncodeStringResponse(ctx context.Context, r interface{}) (interface{}, error) {
	return r, nil
}

//func main() {
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
//	defer cancel()
//	select {
//	case <-time.After(1 * time.Second):
//		fmt.Println("overslept")
//	case <-ctx.Done():
//		fmt.Println(ctx.Err()) // prints "context deadline exceeded"
//	}
//}
