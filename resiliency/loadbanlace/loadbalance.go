package loadbanlace

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"math/rand"
)

type LoadBalance interface {
	SelectService(services []*api.AgentService) (*api.AgentService, error)
}

type RandomLoadBalance struct {
}

// SelectService 随机的选择一个服务实例进行分发.
func (loadBalance *RandomLoadBalance) SelectService(services []*api.AgentService) (*api.AgentService, error) {
	if services == nil || len(services) == 0 {
		return nil, errors.New("service instances are not exist")
	}
	return services[rand.Intn(len(services))], nil
}
