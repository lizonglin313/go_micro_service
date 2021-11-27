package transport

import (
	"context"
	"encoding/json"
	transport_http "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"rpc/support_gRPC_HTTP/endpoint"
)

func MakeHttpHandler(endpoints endpoint.Endpoints) http.Handler {
	// mux := http.NewServeMux() // 构造复用器
	r := mux.NewRouter()
	// 设置工作路由
	r.Methods("POST").Path("/sum").Handler(
		transport_http.NewServer(
			endpoints.SumEndpoint,
			DecodeHttpSumRequest,
			EncodeHttpResponse,
		))

	r.Methods("POST").Path("/concat").Handler(
		transport_http.NewServer(
			endpoints.ConcatEndpoint,
			DecodeHttpConcatRequest,
			EncodeHttpResponse,
		))

	// 由于新增了healthCheck需要加上
	r.Methods("GET").Path("/health").Handler(
		transport_http.NewServer(
			endpoints.HealthCheckEndpoint,
			DecodeHttpHealthCheckRequest,
			EncodeHttpResponse,
		))

	return r
}

// DecodeHttpSumRequest 将HTTP形式的请求转化为结构体.
func DecodeHttpSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var realRequest endpoint.SumRequest
	if err := json.NewDecoder(r.Body).Decode(&realRequest); err != nil {
		return nil, err
	}
	// 返回后交给 endpoint 使用
	return realRequest, nil
}

func DecodeHttpConcatRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var realRequest endpoint.ConcatRequest
	if err := json.NewDecoder(r.Body).Decode(&realRequest); err != nil {
		return nil, err
	}
	// 返回后交给 endpoint 使用
	return realRequest, nil
}

func DecodeHttpHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpoint.HealthCheckRequest{}, nil
}

// EncodeHttpResponse 将endpoint响应结构体返回给Http.
func EncodeHttpResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8") // 设置编码格式
	return json.NewEncoder(w).Encode(response)
}
