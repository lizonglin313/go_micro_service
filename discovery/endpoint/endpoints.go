// Package endpoint 作为服务的实现逻辑和http交互的中间件
package endpoint

import (
	"context"
	"discovery/service"
	"github.com/go-kit/kit/endpoint"
)

// DiscoveryEndpoints 对应我们提供的三个服务.
type DiscoveryEndpoints struct {
	// 每个endpoint 接收请求并返回响应
	SayHelloEndpoint    endpoint.Endpoint
	DiscoveryEndpoint   endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

type SayHelloRequest struct {
}

type SayHelloResponse struct {
	Message string `json:"message"`
}

func MakeSayHelloEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		message := svc.SayHello()
		return SayHelloResponse{
			Message: message,
		}, nil
	}
}

type DiscoveryRequest struct {
	ServiceName string
}

type DiscoveryResponse struct {
	Instances []interface{} `json:"instances"`
	Error     string        `json:"error"`
}

func MakeDiscoveryEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DiscoveryRequest) // 通过断言说明类型
		instances, err := svc.DiscoveryService(ctx, req.ServiceName)
		var errString = ""
		if err != nil {
			errString = err.Error()
		}
		return &DiscoveryResponse{ // 注意这里要返回地址 可能是为了减少开销
			Instances: instances,
			Error:     errString,
		}, nil
	}
}

type HealthRequest struct{}

type HealthResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{
			Status: status,
		}, nil
	}
}
