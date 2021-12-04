package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"net/http"
	"net/url"
	"resiliency/discover"
	"resiliency/loadbanlace"
	"resiliency/string-service/config"
	"strconv"
)

const (
	StringServiceCommandName = "String.string"
	StringService            = "string"
)

var (
	ErrHystrixFallbackExecute = errors.New("hystrix fall back execute")
)

type Service interface {
	UseStringService(operationType, a, b string) (string, error)
	HealthCheck() bool
}

type UseStringService struct {
	discoveryClient discover.DiscoveryClient
	loadbalance     loadbanlace.LoadBalance
}

type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

// NewUseStringService 可以定制服务发现客户端和负载平衡策略.
func NewUseStringService(client discover.DiscoveryClient, lb loadbanlace.LoadBalance) Service {
	hystrix.ConfigureCommand(StringServiceCommandName, hystrix.CommandConfig{
		// 设置请求阈值为 5
		RequestVolumeThreshold: 5,
	})
	return &UseStringService{
		discoveryClient: client,
		loadbalance:     lb,
	}
}

//func (s UseStringService) UseStringService(operationType, a, b string) (string, error) {
//	var operationResult string
//	var err error
//
//	instances := s.discoveryClient.DiscoverServices(StringService, config.Logger)
//	instanceList := make([]*api.AgentService, len(instances))
//	for i := 0; i < len(instances); i++ {
//		// 断言类型
//		instanceList[i] = instances[i].(*api.AgentService)
//	}
//
//	// 使用复杂均衡选择实例执行
//	selectInstance, err := s.loadbalance.SelectService(instanceList)
//	if err == nil {
//		config.Logger.Printf("current string-service ID is %s and address:port is %s:%s",
//			selectInstance.ID, selectInstance.Address, strconv.Itoa(selectInstance.Port))
//		requestUrl := url.URL{ // 构造请求url
//			Scheme: "http",
//			Host:   selectInstance.Address + ":" + strconv.Itoa(selectInstance.Port),
//			Path:   "/op/" + operationType + "/" + a + "/" + b,
//		}
//		// 进行请求
//		resp, err := http.Post(requestUrl.String(), "", nil)
//		if err == nil {
//			result := &StringResponse{}
//			err = json.NewDecoder(resp.Body).Decode(result)
//			if err == nil && result.Error == nil {
//				operationResult = result.Result
//			}
//		}
//	}
//	return operationResult, err
//}

func (s UseStringService) UseStringService(operationType, a, b string) (string, error) {
	var operationResult string
	err := hystrix.Do(StringServiceCommandName, func() error {
		instances := s.discoveryClient.DiscoverServices(StringService, config.Logger)
		fmt.Println("LENNNNNNNNNG is %d", len(instances))
		// 随机选取实例执行
		instanceList := make([]*api.AgentService, len(instances))
		for i := 0; i < len(instances); i++ {
			instanceList[i] = instances[i].(*api.AgentService)
		}

		selectInstance, err := s.loadbalance.SelectService(instanceList)
		if err != nil {
			config.Logger.Println(err.Error())
			return err
		}
		config.Logger.Printf("current string-service ID is %s and address:port is %s:%s\n", selectInstance.ID, selectInstance.Address, strconv.Itoa(selectInstance.Port))
		requestUrl := url.URL{
			Scheme: "http",
			Host:   selectInstance.Address + ":" + strconv.Itoa(selectInstance.Port),
			Path:   "/op/" + operationType + "/" + a + "/" + b,
		}

		resp, err := http.Post(requestUrl.String(), "", nil)
		if err != nil {
			return err
		}
		result := &StringResponse{}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		} else if result.Error != nil {
			return result.Error
		}

		operationResult = result.Result
		return nil
	},
		func(e error) error {
			// 这是定义的一个简单的失败回滚函数
			// 如果发生错误 如该服务的熔断器已经打开 则直接返回错误 进行服务熔断
			return ErrHystrixFallbackExecute
		})
	return operationResult, err
}

func (u UseStringService) HealthCheck() bool {
	return true
}

type ServiceMiddleware func(Service) Service
