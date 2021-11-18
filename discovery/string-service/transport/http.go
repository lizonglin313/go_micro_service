package transport

import (
	"context"
	endpts "discovery/endpoint"
	"discovery/string-service/endpoint"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// 构造httphandler 使用 discoverEndpoint
func MakeHttpHandler(ctx context.Context, endpoints endpts.DiscoveryEndpoints, logger log.Logger) http.Handler {
	// 实例化一个多路服用处理请求
	r := mux.NewRouter()

	// 定义处理器
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	// 处理 say hello 请求
	r.Methods("GET").Path("/say-hello").Handler(kithttp.NewServer(
		endpoints.SayHelloEndpoint,
		decodeSayHelloRequest,
		encodeJsonResponse,
		options...,
	))

	// 处理服务发现
	r.Methods("GET").Path("discovery").Handler(kithttp.NewServer(
		endpoints.DiscoveryEndpoint,
		decodeDiscoveryRequest,
		encodeJsonResponse,
		options...))

	// 处理健康检查
	r.Methods("GET").Path("health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJsonResponse,
		options...))

	return r
}

// decodeSayHelloRequest
// @Desc: 	对 sayhello 请求参数进行解码
// @Param:	_
// @Param:	r
// @Return:	interface{}
// @Return:	error
// @Notice:
func decodeSayHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpts.SayHelloRequest{}, nil
}

func decodeDiscoveryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	serviceName := r.URL.Query().Get("serviceName")
	if serviceName == "" {
		return nil, ErrorBadRequest
	}
	return endpts.DiscoveryRequest{ServiceName: serviceName}, nil
}

func decodeHealthCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

// encodeJsonResponse
// @Desc: 	将响应结果 编码为json的形式返回
// @Param:	ctx
// @Param:	w
// @Param:	response
// @Return:	error
// @Notice:
func encodeJsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// MakeHttpHandler make http handler use mux
//func MakeHttpHandler(ctx context.Context, endpoints endpoint.StringEndpoints, logger log.Logger) http.Handler {
//	r := mux.NewRouter()
//
//	options := []kithttp.ServerOption{
//		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
//		kithttp.ServerErrorEncoder(encodeError),
//	}
//
//	//  每个kithttp.Server方法
//	//	需要传入：
//	//	1) 该请求对应的处理器
//	//	2) 一个解析请求参数的 请求解码器
//	//	3) 一个编码请求相应的 响应编码器
//	r.Methods("POST").Path("/op/{type}/{a}/{b}").Handler(kithttp.NewServer(
//		endpoints.StringEndpoint,
//		decodeStringRequest,
//		encodeStringResponse,
//		options...,
//	))
//
//	// prometheus 这个是开源的做监控的组件
//	r.Path("/metrics").Handler(promhttp.Handler())
//
//	// create health check handler
//	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
//		endpoints.HealthCheckEndpoint,
//		decodeHealthCheckRequest,
//		encodeStringResponse,
//		options...,
//	))
//
//	return r
//}

// decodeStringRequest decode request params to struct
func decodeStringRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	requestType, ok := vars["type"]
	if !ok {
		return nil, ErrorBadRequest
	}

	pa, ok := vars["a"]
	if !ok {
		return nil, ErrorBadRequest
	}

	pb, ok := vars["b"]
	if !ok {
		return nil, ErrorBadRequest
	}

	return endpoint.StringRequest{
		RequestType: requestType,
		A:           pa,
		B:           pb,
	}, nil
}

// encodeStringResponse encode response to return
func encodeStringResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decodeHealthCheckRequest decode request
//func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
//	return endpoint.HealthRequest{}, nil
//}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
