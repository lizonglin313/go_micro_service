package main

import (
	"context"
	"flag"
	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"net"
	"os"
	"rpc/go_kit_gRPC/proto"
	"rpc/go_kit_gRPC/service"
)

func main() {
	flag.Parse()
	ctx := context.Background()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// 最原始的 Service: Concat(a, b string) string | Diff(a,b string) string
	var svc service.Service       // 声明接口
	svc = service.StringService{} // 实例化

	// 接入中间件
	svc = service.LoggingMiddleware(logger)(svc)

	// 根据服务创建 Endpoint
	// Endpoint 对 Service 做的封装是：
	// 1) 把 protoReq 转化成的 StringRequest 中的 a,b 拿出来交给 Service
	// 2) 把 Service 处理结果 string 包装成 StringResponse 交给 proto 做 protoResp 转换
	stringEndpoint := service.MakeStringEndpoint(svc)
	healthEndpoint := service.MakeHealthCheckEndpoint(svc)

	// 封装到一起
	endpts := service.StringEndpoints{
		StringEndpoint: stringEndpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	// 把 Endpoint 封装到 transport 中
	handler := service.NewStringServerGRPC(ctx, endpts)

	ls, _ := net.Listen("tcp", "127.0.0.1:8123")
	gRPCServer := grpc.NewServer()
	// 注册服务
	proto.RegisterStringServiceServer(gRPCServer, handler)
	// 开启服务
	gRPCServer.Serve(ls)
}
