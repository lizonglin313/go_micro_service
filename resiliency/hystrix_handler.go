package main

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"net/http/httputil"
	"resiliency/discover"
	"resiliency/loadbalance"
	"strings"
	"sync"
)

var (
	ErrNoInStance = errors.New("query service instance error")
)

// HystrixHandler 实现http.Handler接口.
// 表明可以处理http服务
type HystrixHandler struct {
	// 记录hystrix是否已经配置
	hystrixs      map[string]bool
	hystrixsMutex *sync.Mutex
	// 配置服务发现/负载均衡/日志
	discoveryClient discover.DiscoveryClient
	loadbalance     loadbalance.LoadBalance
	logger          *log.Logger
}

// NewHystrixHandler 构造函数.
func NewHystrixHandler(discoveryClient discover.DiscoveryClient,
	loadbalance loadbalance.LoadBalance, logger *log.Logger) *HystrixHandler {
	return &HystrixHandler{
		hystrixs:        make(map[string]bool),
		hystrixsMutex:   &sync.Mutex{},
		discoveryClient: discoveryClient,
		loadbalance:     loadbalance,
		logger:          logger,
	}
}

func (hystrixHandler *HystrixHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// 处理URL获取服务名称
	reqPath := req.URL.Path
	if reqPath == "" {
		return
	}
	pathArray := strings.Split(reqPath, "/")
	fmt.Println(pathArray) // 用于调试
	serviceName := pathArray[1]
	if serviceName == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	// 如果该服务没有在hystrix配置 将该服务加入hystrix配置
	if _, ok := hystrixHandler.hystrixs[serviceName]; !ok {
		hystrixHandler.hystrixsMutex.Lock()                     // map加锁读写
		if _, ok := hystrixHandler.hystrixs[serviceName]; !ok { // 这是为了防止加锁前有写
			// 将serviceName作为hystrix命令名
			hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{
				// 这里可以进行命令自定义
			})
			hystrixHandler.hystrixs[serviceName] = true
		}
		hystrixHandler.hystrixsMutex.Unlock()
	}

	// 以同步方式进行链路监控
	err := hystrix.Do(serviceName, func() error {
		// 调用consul api查询服务实例列表
		instances := hystrixHandler.discoveryClient.DiscoverServices(serviceName, hystrixHandler.logger)
		instanceList := make([]*api.AgentService, len(instances))
		for i := 0; i < len(instances); i++ {
			instanceList[i] = instances[i].(*api.AgentService)
		}

		// 负载均衡选择实例
		selectInstance, err := hystrixHandler.loadbalance.SelectService(instanceList)
		if err != nil {
			return ErrNoInStance // 返回 没有可以使用的实例
		}

		// 用该实例创建转发器: 接收原始请求 返回 转发请求的路径
		// 用于构建反向代理 告诉反向代理 进来的请求应该转发到哪个服务路径
		director := func(req *http.Request) {
			// 组织请求路径
			destPath := strings.Join(pathArray[2:], "/")

			hystrixHandler.logger.Println("service Id ", selectInstance.ID)

			// 设置代理服务地址信息
			req.URL.Scheme = "http"
			req.URL.Host = fmt.Sprintf("%s:%d", selectInstance.Address, selectInstance.Port)
			req.URL.Path = "/" + destPath

			fmt.Println(req.URL.Host + req.URL.Path) // 用于调试
		}

		var proxyError error

		// 返回代理异常 用于记录 hystrix.Do 执行失败
		// 用于构建反向代理 告诉反向代理异常该做什么
		errorHandler := func(ew http.ResponseWriter, er *http.Request, err error) {
			proxyError = err
		}

		// 创建反向代理
		proxy := &httputil.ReverseProxy{
			Director: director,
			ErrorHandler: errorHandler,
		}

		// 反向代理开始工作 进行监听
		proxy.ServeHTTP(rw, req)

		return proxyError	// 将异常情况返回给 hystrix
	}, func(err error) error {
		hystrixHandler.logger.Println("proxy error", err)
		return errors.New("fallback excute")	// 该出错处理函数回调执行
	})

	// 如果 hystrix.Do 执行异常
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}

}
