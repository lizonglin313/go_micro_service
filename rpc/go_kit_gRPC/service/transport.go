package service

import (
	"context"
	"github.com/go-kit/kit/transport/grpc"
	"rpc/go_kit_gRPC/proto"
)

// grpcServer 构建使用grpc通信的服务.
type grpcServer struct {
	// 使用小写属性名一方面避免包外暴露
	// 另一方面避免与实现的接口方法冲突
	concat grpc.Handler
	diff   grpc.Handler
}

func (s *grpcServer) Concat(ctx context.Context, request *proto.StringRequest) (*proto.StringResponse, error) {
	// 相当于把上下文和请求继续下传一层
	// 事实上 ServerGRPC 帮我们完成了处理工作
	// Request -> ServerGRPC -> Response
	_, resp, err := s.concat.ServeGRPC(ctx, request)	// return: ctx, resp, err
	if err != nil {
		return nil, err
	}
	return resp.(*proto.StringResponse), nil	// 使用断言强制转化为返回类型
}

func (s *grpcServer) Diff(ctx context.Context, request *proto.StringRequest) (*proto.StringResponse, error) {
	_, resp, err := s.diff.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.StringResponse), nil
}

func NewStringServerGRPC(ctx context.Context, endpoints StringEndpoints) proto.StringServiceServer {
	return &grpcServer{
		concat: grpc.NewServer(),
		diff: grpc.NewServer(),
	}
}

