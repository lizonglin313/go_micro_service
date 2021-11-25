package transport

import (
	"context"
	"encoding/json"
	transport_http "github.com/go-kit/kit/transport/http"
	"net/http"
	"rpc/support_gRPC_HTTP/endpoint"
)

func MakeHttpHandler(endpoints endpoint.Endpoints) http.Handler {
	mux := http.NewServeMux() // 构造复用器

	// 设置工作路由
	mux.Handle("/sum",
		transport_http.NewServer(
			endpoints.SumEndpoint,
			DecodeHttpSumRequest,
			EncodeHttpResponse,
		))

	mux.Handle("/concat",
		transport_http.NewServer(
			endpoints.ConcatEndpoint,
			DecodeHttpConcatRequest,
			EncodeHttpResponse,
		))

	return mux
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

// EncodeHttpResponse 将endpoint响应结构体返回给Http.
func EncodeHttpResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
