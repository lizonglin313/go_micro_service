package transport

import (
	"context"
	transport_gRPC "github.com/go-kit/kit/transport/grpc"
	"rpc/support_gRPC_HTTP/endpoint"
	"rpc/support_gRPC_HTTP/proto"
)

// MakeGRPCServer 在使用http里面对应的是MakeHttpHandler.
// 注意 这里处理的是 服务端
func MakeGRPCServer(endpoints endpoint.Endpoints) proto.SumConcatServer {

}


// grpcServer 可以理解成http中mux handler封装的两个路由
// 只不过grpcServer封装了两个grpc Handler
// 注意这里要实现 proto 里面的服务接口.
type grpcServer struct {
	sum    transport_gRPC.Handler
	concat transport_gRPC.Handler
}

func (g *grpcServer) Sum(ctx context.Context, request *proto.SumRequest) (*proto.SumResponse, error) {
	_, resp, err := g.sum.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.SumResponse), nil
}

func (g *grpcServer) Concat(ctx context.Context, request *proto.ConcatRequest) (*proto.ConcatResponse, error) {
	_, resp, err := g.concat.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.ConcatResponse), nil
}

// 进行复杂的编解码工作



