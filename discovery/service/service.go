package service

import (
	"context"
	"discovery/config"
	"discovery/discover"
	"errors"
)

//  Service
//  @Description: 真正本服务要做的接口
type Service interface {


	// HealthCheck check service health status
	HealthCheck() bool

	// sayHelloService
	SayHello() string

	//  discovery service from consul by serviceName
	// 直接使用与consul的交互接口
	DiscoveryService(ctx context.Context, serviceName string) ([]interface{}, error)

}


var ErrNotServiceInstances = errors.New("instances are not existed")

// 实现客户端的服务发现接口，有注册、注销、发现三个方法
type DiscoveryServiceImpl struct {
	// 需要使用接口里的方法
	discoveryClient discover.DiscoveryClient
}

// NewDiscoveryServiceImpl
// @Desc: 	实例的构造函数，构造一个服务实例并进行返回
//			这个服务实例是带着与consul交互的接口的
// @Param:	discoveryClient
// @Return:	Service
// @Notice:
func NewDiscoveryServiceImpl(discoveryClient discover.DiscoveryClient) Service  {
	return &DiscoveryServiceImpl{
		discoveryClient:discoveryClient,
	}
}

// SayHello
// @Desc: 	实现打招呼方法
// @Rece:	*DiscoveryServiceImpl
// @Return:	string
// @Notice:
func (*DiscoveryServiceImpl) SayHello() string {
	return "Hello World!"
}

// DiscoveryService
// @Desc: 	实现服务发现方法
// @Rece:	service
// @Param:	ctx
// @Param:	serviceName
// @Return:	[]interface{}
// @Return:	error
// @Notice:
func (service *DiscoveryServiceImpl) DiscoveryService(ctx context.Context, serviceName string) ([]interface{}, error)  {

	instances := service.discoveryClient.DiscoverServices(serviceName, config.Logger)

	if instances == nil || len(instances) == 0 {
		return nil, ErrNotServiceInstances
	}
	return instances, nil
}


// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (*DiscoveryServiceImpl) HealthCheck() bool {
	return true
}

