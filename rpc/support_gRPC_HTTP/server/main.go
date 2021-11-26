package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"rpc/support_gRPC_HTTP/endpoint"
	"rpc/support_gRPC_HTTP/proto"
	"rpc/support_gRPC_HTTP/service"
	"rpc/support_gRPC_HTTP/transport"
	"syscall"
)

// main 同时实现HTTP和gRPC的服务.
func main() {
	httpAddr := flag.String("HTTP", ":8123", "HTTP Server")
	gRPCAddr := flag.String("gRPC", ":9123", "gRPC Server")
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	logger.Log("msg", "Server Start...")
	defer logger.Log("msg", "Server Closed")

	svc := service.NewSumConcatService()

	endpoints := endpoint.Endpoints{
		SumEndpoint:    endpoint.MakeSumEndpoint(svc),
		ConcatEndpoint: endpoint.MakeConcatEndpoint(svc),
	}

	errchan := make(chan error)

	// 启一个 goroutine 等待接收 ctrl+c 中断信号
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errchan <- fmt.Errorf("%s", <-c)
	}()

	// 处理 HTTP 服务
	// 注意 使用 HTTP 用 POST方法发送JSON数据
	go func() {
		logger := log.With(logger, "transport", "HTTP")
		logger.Log("addr", *httpAddr)

		httpHandler := transport.MakeHttpHandler(endpoints)
		errchan <- http.ListenAndServe(*httpAddr, httpHandler)
	}()

	// 处理 gRPC 服务
	go func() {
		logger := log.With(logger, "transport", "gRPC")
		logger.Log("addr", *gRPCAddr)

		// 创建TCP的监听接口
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errchan <- err
			return
		}

		gRPCHandler := transport.MakeGRPCServer(endpoints)
		gRPCServer := grpc.NewServer()
		// 注册服务
		proto.RegisterSumConcatServer(gRPCServer, gRPCHandler)
		// 让服务监听这个接口
		errchan <- gRPCServer.Serve(listener)
	}()

	logger.Log("exit", <-errchan)

}
