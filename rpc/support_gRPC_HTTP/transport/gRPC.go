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
	return &grpcServer{
		sum: transport_gRPC.NewServer(
			endpoints.SumEndpoint,
			DecodeGRPCSumRequest,
			EncodeGRPCSumResponse,
		),
		concat: transport_gRPC.NewServer(
			endpoints.ConcatEndpoint,
			DecodeGRPCConcatRequest,
			EncodeGRPCConcatResponse,
		),
	}
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
// 由于 gRPC的transport也要使用Endpoint进行后续的处理
// 所以需要将gRPC的pb格式的请求和响应
// 转换成Endpoint中定义的请求和响应结构

// DecodeGRPCSumRequest : PB Sum Req -> EP Sum Req.
func DecodeGRPCSumRequest(_ context.Context, request interface{}) (interface{}, error) {
	pbReq := request.(*proto.SumRequest)
	return endpoint.SumRequest{A: int(pbReq.A), B: int(pbReq.B)}, nil
}

// EncodeGRPCSumRequest : EP Sum Req -> PB Sum Req.
func EncodeGRPCSumRequest(_ context.Context, request interface{}) (interface{}, error) {
	epReq := request.(endpoint.SumRequest)
	return &proto.SumRequest{A: int64(epReq.A), B: int64(epReq.B)}, nil
}

// DecodeGRPCSumResponse : PB Sum Resp -> EP Sum Resp.
func DecodeGRPCSumResponse(_ context.Context, response interface{}) (interface{}, error) {
	pbResp := response.(*proto.SumResponse)
	return endpoint.SumResponse{Reply: int(pbResp.Reply)}, nil
}

// EncodeGRPCSumResponse : EP Sum Resp -> PB Sum Resp.
func EncodeGRPCSumResponse(_ context.Context, response interface{}) (interface{}, error) {
	epResp := response.(endpoint.SumResponse)
	return &proto.SumResponse{Reply: int64(epResp.Reply)}, nil
}

// DecodeGRPCConcatRequest : PB Concat Req -> EP Concat Req.
func DecodeGRPCConcatRequest(_ context.Context, request interface{}) (interface{}, error) {
	pbReq := request.(*proto.ConcatRequest)
	return endpoint.ConcatRequest{A: pbReq.A, B: pbReq.B}, nil
}

// EncodeGRPCConcatRequest : EP Concat Req -> PB Concat Req.
func EncodeGRPCConcatRequest(_ context.Context, request interface{}) (interface{}, error) {
	epReq := request.(endpoint.ConcatRequest)
	return &proto.ConcatRequest{A: epReq.A, B: epReq.B}, nil
}

// DecodeGRPCConcatResponse : PB Concat Resp -> EP Concat Resp.
func DecodeGRPCConcatResponse(_ context.Context, response interface{}) (interface{}, error) {
	pbResp := response.(*proto.ConcatResponse)
	return endpoint.ConcatResponse{Reply: pbResp.Reply}, nil
}

// EncodeGRPCConcatResponse : EP Concat Resp -> PB Concat Resp.
func EncodeGRPCConcatResponse(_ context.Context, response interface{}) (interface{}, error) {
	epResp := response.(endpoint.ConcatResponse)
	return &proto.ConcatResponse{Reply: epResp.Reply}, nil
}
