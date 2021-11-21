package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	stream_pb "rpc/stream_gRPC/pb"
	"strconv"
)

func main() {
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()

	// 生成客户端
	stringClient := stream_pb.NewStringServiceClient(conn)
	// 构造请求
	stringReq := &stream_pb.StringRequest{A: "A", B: "B"}
	// 向服务端请求，服务端以流的形式响应
	stream, _ := stringClient.LotsOfServerStream(context.Background(), stringReq)
	for {
		// 持续接收
		item, stream_error := stream.Recv()
		// 至流结束
		if stream_error == io.EOF {
			break
		}
		if stream_error != nil {
			log.Printf("failed to recv: %v", stream_error)
		}
		fmt.Printf("StringService Concat : %s concat %s = %s\n", stringReq.A, stringReq.B, item.GetRet())
	}

	// 再以流的形式向服务端发送请求
	sendClientStreamRequest(stringClient)

	// 双向流的形式
	sendClientAndServerStreamRequest(stringClient)
}

func sendClientStreamRequest(client stream_pb.StringServiceClient) {
	fmt.Printf("test sendClientStreamRequest \n")
	// 建立流连接
	stream, err := client.LotsOfClientStream(context.Background())
	for i := 0; i < 10; i++ {
		if err != nil {
			log.Printf("failed to call: %v", err)
			break
		}
		stream.Send(&stream_pb.StringRequest{A: strconv.Itoa(i), B: strconv.Itoa(i + 1)})
	}
	// 流请求结束并接收响应
	reply, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("failed to recv: %v", err)
	}
	log.Printf("sendClientStreamRequest ret is : %s", reply.Ret)
}

func sendClientAndServerStreamRequest(client stream_pb.StringServiceClient) {
	fmt.Printf("test sendClientAndServerStreamRequest \n")
	var err error
	// 建立流连接
	stream, err := client.LotsOfServerAndClientStream(context.Background())
	if err != nil {
		log.Printf("failed to call: %v", err)
		return
	}
	var i int
	for { // 在一次for里 又发送 又接收
		if i > 10 {
			stream.CloseSend()	// 在这里通知服务器(EOF)我请求流结束
			break	// 避免无限循环
		}
		err1 := stream.Send(&stream_pb.StringRequest{A: strconv.Itoa(i), B: strconv.Itoa(i + 1)})
		if err1 != nil {
			log.Printf("failed to send: %v", err)
			break
		}
		reply, err2 := stream.Recv()
		if err2 != nil {
			log.Printf("failed to recv: %v", err)
			break
		}
		log.Printf("sendClientAndServerStreamRequest Ret is : %s", reply.Ret)
		i++
	}
}
