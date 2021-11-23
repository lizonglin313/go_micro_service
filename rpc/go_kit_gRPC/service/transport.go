package service

import (
	"context"
	"github.com/go-kit/kit/transport/grpc"
	"rpc/go_kit_gRPC/proto"
)

// grpcServer 构建使用grpc通信的服务.
// 所以他要实现proto的concat和diff接口
// 注意是接收 proto格式的Req和Resp
type grpcServer struct {
	// 使用小写属性名一方面避免包外暴露
	// 另一方面避免与实现的接口方法冲突
	concat grpc.Handler
	diff   grpc.Handler
}

func (s *grpcServer) Concat(ctx context.Context, request *proto.StringRequest) (*proto.StringResponse, error) {
	// 相当于把上下文和请求继续下传一层
	// 事实上 ServerGRPC 帮我们完成了处理工作
	// ProtoRequest -> ServerGRPC(->ServiceReq -> Service -> ServiceResp) -> ProtoResponse
	_, resp, err := s.concat.ServeGRPC(ctx, request) // return: ctx, resp, err
	if err != nil {
		return nil, err
	}
	return resp.(*proto.StringResponse), nil // 使用断言强制转化为返回类型
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
		concat: grpc.NewServer(
			endpoints.StringEndpoint,
			DecodeDiffStringRequest,
			EncodeDiffResponse),
		diff:   grpc.NewServer(
			endpoints.StringEndpoint,
			DecodeDiffStringRequest,
			EncodeDiffResponse),
	}
}

// DecodeDiffStringRequest : protoReq -> Decoder -> ServiceReq.
func DecodeDiffStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*proto.StringRequest)
	return StringRequest{
		RequestType: "Diff",
		A:           string(req.A),
		B:           string(req.B),
	}, nil
}

// EncodeDiffResponse : ServiceResp -> Encoder -> protoResp.
func EncodeDiffResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(StringResponse)

	if resp.Error != nil {
		return &proto.StringResponse{
			Ret: resp.Result,
			Err: resp.Error.Error(),
		}, nil
	}

	return &proto.StringResponse{
		Ret: resp.Result,
		Err: "",
	}, nil
}
