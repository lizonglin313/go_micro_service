package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/transport/grpc"
	"rpc/go_kit_gRPC/proto"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
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
	// TODO: 调试工作
	fmt.Println("Request in transport.go Concat:", request)
	_, resp, err := s.concat.ServeGRPC(ctx, request) // return: ctx, resp, err
	if err != nil {
		return nil, err
	}

	// TODO: 初步判断这里作为 main 的返回点 proto.StringResponse 出现了问题
	protoResp, ok := resp.(*proto.StringResponse)
	if !ok {
		fmt.Println("Interface con fail in transport.go Concat!")
	}
	return protoResp, nil // 使用断言强制转化为返回类型
}

func (s *grpcServer) Diff(ctx context.Context, request *proto.StringRequest) (*proto.StringResponse, error) {
	_, resp, err := s.diff.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	protoResp, ok := resp.(*proto.StringResponse)
	if !ok {
		fmt.Println("Interface con fail in transport.go Concat!")
	}
	return protoResp, nil
}

func NewStringServerGRPC(ctx context.Context, endpoints StringEndpoints) proto.StringServiceServer {
	return &grpcServer{
		concat: grpc.NewServer(
			endpoints.StringEndpoint,
			DecodeConcatStringRequest,
			EncodeStringResponse),
		diff: grpc.NewServer(
			endpoints.StringEndpoint,
			DecodeDiffStringRequest,
			EncodeStringResponse),
	}
}

// DecodeConcatStringRequest : protoReq -> Decoder -> ServiceReq.
func DecodeConcatStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*proto.StringRequest)
	return StringRequest{
		RequestType: "Concat",
		A:           string(req.A),
		B:           string(req.B),
	}, nil
}

func EncodeConcatStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	// req := r.(StringRequest)
	// return &proto.StringRequest{A: req.A, B: req.B}, nil
	return r, nil
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

func EncodeDiffStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(StringRequest)
	return &proto.StringRequest{A: req.A, B: req.B}, nil
}

func DecodeStringResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*proto.StringResponse)
	return StringResponse{
		Result: resp.Ret,
		Error:  nil,
	}, nil
}

// EncodeStringResponse : ServiceResp -> Encoder -> protoResp.
func EncodeStringResponse(_ context.Context, r interface{}) (interface{}, error) {
	fmt.Println("r in transport.go EncodeStringResponse is:", r)
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
