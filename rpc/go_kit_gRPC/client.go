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
	//ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	//defer cancel()
	//conn, err := grpc.DialContext(ctx, "127.0.0.1:8123", grpc.WithInsecure())

	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:8123", grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))

	if err != nil {
		fmt.Println("gRPC dial err:", err)
	}
	defer conn.Close()

	// TODO: 排错
	svr := NewStringConcatClient(conn)
	result, err := svr.Concat(ctx, "A", "B")

	if err != nil {
		fmt.Println("Check error:", err.Error())
	}

	fmt.Println("Concat result=", result)

	svr1 := NewStringDiffClient(conn)
	res, err := svr1.Diff(ctx, "ABD", "DBG")
	if err != nil {
		fmt.Println("Check error:", err.Error())
	}
	fmt.Println("Diff result=", res)

}

func NewStringConcatClient(conn *grpc.ClientConn) service.Service {
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

func NewStringDiffClient(conn *grpc.ClientConn) service.Service {
	var ep = grpctransport.NewClient(conn,
		"proto.StringService",
		"Diff",
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
	fmt.Println("request from client.go DecoderStringRequest:", r)
	return r, nil
}

func EncodeStringResponse(_ context.Context, r interface{}) (interface{}, error) {
	fmt.Println("response from client.go EncoderStringResponse:", r)
	return r, nil
}
