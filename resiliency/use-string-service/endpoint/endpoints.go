package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"resiliency/use-string-service/service"
)

type UseStringEndpoints struct {
	UseStringEndpoint   endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

type UseStringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

type UseStringResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

//func MakeUseStringEndpoint(svc service.Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		req := request.(UseStringRequest)
//		var (
//			res, a, b, opErrorString string
//			opError                  error
//		)
//
//		a = req.A
//		b = req.B
//
//		res, opError = svc.UseStringService(req.RequestType, a, b)
//		if opError != nil {
//			opErrorString = opError.Error()
//		}
//		return UseStringResponse{Result: res, Error: opErrorString}, nil
//	}
//}

// MakeUseStringEndpoint 这个版本用于go-kit的hystrix.
func MakeUseStringEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UseStringRequest)

		var (
			res, a, b string
			opError   error
		)

		a = req.A
		b = req.B

		// 这里直接将操作的错误返回给 transport 处理
		res, opError = svc.UseStringService(req.RequestType, a, b)
		return UseStringResponse{Result: res}, opError
	}
}

type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{status}, nil
	}
}
