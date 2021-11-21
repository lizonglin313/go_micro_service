// Package string_service 这个包只关注实现服务,也可以说是服务端的内容.
package string_service

import (
	"context"
	"errors"
	"io"
	"log"
	stream_pb "rpc/stream_gRPC/pb"
	"strconv"
	"strings"
)

const (
	StrMaxSize = 1024
)

// Service errors
var (
	ErrMaxSize = errors.New("maximum size of 1024 bytes exceeded")

	ErrStrValue = errors.New("maximum size of 1024 bytes exceeded")
)

type StringService struct{}

// LotsOfServerStream 服务端以流的方式向客户端返回数据.
func (s *StringService) LotsOfServerStream(req *stream_pb.StringRequest, qs stream_pb.StringService_LotsOfServerStreamServer) error {
	for i := 0; i < 10; i++ {
		response := stream_pb.StringResponse{Ret: req.A + req.B + strconv.Itoa(i)}
		// Send 的内容属于消息流的一部分，直到 return
		qs.Send(&response)
	}
	return nil
}

// LotsOfClientStream 客户端以流的方式向服务端发送请求.
func (s *StringService) LotsOfClientStream(qs stream_pb.StringService_LotsOfClientStreamServer) error {
	var params []string
	for {
		// 持续接收客户端的请求流
		in, err := qs.Recv()
		// 直到请求流结束，收到EOF信号
		if err == io.EOF {
			// 再统一进行处理请求进行响应
			qs.SendAndClose(&stream_pb.StringResponse{Ret: strings.Join(params, "-")})
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}
		params = append(params, in.A, in.B)
	}
}

// LotsOfServerAndClientStream 客户端以流的方式请求,服务端以流的方式响应.
func (s *StringService) LotsOfServerAndClientStream(qs stream_pb.StringService_LotsOfServerAndClientStreamServer) error {
	n := 0
	for {	// 在一次for里 即接收 又 响应
		// 持续接收请求流
		in, err := qs.Recv()
		if err == io.EOF {
			// 请求流结束再退出
			return nil
		}
		if err != nil {
			log.Printf("failed to recv %v", err)
			return err
		}
		// 每接收一个请求给出一个流式响应
		qs.Send(&stream_pb.StringResponse{Ret: in.A + in.B + strconv.Itoa(n)})
		n++	// 标记响应次数
	}
	return nil
}

func (s *StringService) Concat(ctx context.Context, req *stream_pb.StringRequest) (*stream_pb.StringResponse, error) {
	if len(req.A)+len(req.B) > StrMaxSize {
		response := stream_pb.StringResponse{Ret: ""}
		return &response, nil
	}
	response := stream_pb.StringResponse{Ret: req.A + req.B}
	return &response, nil
}
