package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	slog "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"rpc/support_gRPC_HTTP/discover"
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
	consulPort := flag.Int("consul.port", 8500, "consul port")
	consulHost := flag.String("consul.host", "127.0.0.1", "consul host")
	serviceName := flag.String("service.name", "SumConcatService", "service name")

	servicePort := flag.Int("service.port", 8123, "service port")
	serviceHost := flag.String("service.host", "127.0.0.1", "service host")

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
		SumEndpoint:         endpoint.MakeSumEndpoint(svc),
		ConcatEndpoint:      endpoint.MakeConcatEndpoint(svc),
		HealthCheckEndpoint: endpoint.MakeHealthCheckEndpoint(svc),
	}

	errchan := make(chan error)

	// 构建 consul 客户端
	var discoveryClient discover.DiscoveryClient
	instanceId := *serviceName + "-" + uuid.NewV4().String()
	discoveryClient, err := discover.NewKitDiscoverClient(*consulHost, *consulPort)
	if err != nil {
		logger.Log("Get Consul Client failed:", instanceId)
		os.Exit(-1)
	}

	// 处理 HTTP 服务
	// 注意 使用 HTTP 用 POST方法发送JSON数据
	go func() {
		logger := log.With(logger, "transport", "HTTP")
		logger.Log("addr", *httpAddr)

		httpHandler := transport.MakeHttpHandler(endpoints)
		//启动前执行注册
		if !discoveryClient.Register(*serviceName, instanceId, "/health", *serviceHost,  *servicePort, nil, &slog.Logger{}){
			slog.Printf("string-service for service %s failed.", serviceName)
			// 注册失败，服务启动失败
			os.Exit(-1)
		}
		errchan <- http.ListenAndServe(*httpAddr, httpHandler)
	}()

	// 启一个 goroutine 等待接收 ctrl+c 中断信号
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errchan <- fmt.Errorf("%s", <-c)
		// 退出时向consul注销服务
		discoveryClient.DeRegister(instanceId, &slog.Logger{})
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
