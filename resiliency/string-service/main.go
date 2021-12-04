package main

import (
	"context"
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"os"
	"os/signal"
	"resiliency/discover"
	"resiliency/string-service/config"
	"resiliency/string-service/endpoint"
	"resiliency/string-service/service"
	"resiliency/string-service/transport"
	"strconv"
	"syscall"
)

func main() {
	// 获取命令行参数
	var (
		servicePort = flag.Int("service.port", 10085, "service port")
		serviceHost = flag.String("service.host", "127.0.0.1", "service host")
		consulPort  = flag.Int("consul.port", 8500, "consul port")
		consulHost  = flag.String("consul.host", "127.0.0.1", "consul host")
		serviceName = flag.String("service.name", "string", "service name")
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)
	var discoveryClient discover.DiscoveryClient

	discoveryClient, err := discover.NewKitDiscoverClient(*consulHost, *consulPort)
	if err != nil {
		config.Logger.Println("Get Consul Client failed")
		os.Exit(-1)
	}

	var svc service.Service
	svc = service.StringService{}
	stringEndpoint := endpoint.MakeStringEndpoint(svc)
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	endpts := endpoint.StringEndpoints{
		StringEndpoint:      stringEndpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	r := transport.MakeHttpHandler(ctx, endpts, config.KitLogger)

	instanceId := *serviceName + "-" + uuid.NewV4().String()

	go func() {
		config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
		if !discoveryClient.Register(*serviceName, instanceId, "/health",
			*serviceHost, *servicePort, nil, config.Logger) {
			config.Logger.Printf("string-service for service %s failed.\n", *serviceName)
			os.Exit(-1)
		}
		handler := r
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	discoveryClient.DeRegister(instanceId, config.Logger)
	config.Logger.Println(error)

}
