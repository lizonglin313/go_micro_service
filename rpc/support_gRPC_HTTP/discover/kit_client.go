package discover

import (
	"fmt"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"log"
	"strconv"
	"sync"
)

type KitDiscoverClient struct {
	ConsulHost string
	ConsulPort int
	client     consul.Client

	// 用于连接 consul 的配置
	config *api.Config
	mutex  sync.Mutex

	// 用于服务实例的本地缓存 减轻 consul 压力
	instancesMap sync.Map
}

func (consulClient *KitDiscoverClient) Register(serviceName, instanceId, healthCheckUrl string,
	instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool {
	// 1. 构建服务注册数据
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Port:    instancePort,
		Address: instanceHost,
		Meta:    meta,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:                       "15s",
		},
	}

	// 2. 进行注册
	err := consulClient.client.Register(serviceRegistration)
	if err != nil {
		// Ubuntu报错就在这里 因为传入的 logger 是临时定义的，所以在使用时会报空指针错误
		// logger.Println("Register Service Error!")
		fmt.Println("Register Service Error!")
		return false
	}
	// logger.Println("Register Service Success!")
	fmt.Println("Register Service Success!")
	return true
}

func (consulClient *KitDiscoverClient) DeRegister(instanceId string, logger *log.Logger) bool {
	// 1. 构建服务注销需要的数据
	serviceRegistration := &api.AgentServiceRegistration{
		ID: instanceId,
	}
	// 2. 向consul注销服务
	err := consulClient.client.Deregister(serviceRegistration)
	if err != nil {
		logger.Println("Deregister Service Error!")
		return false
	}
	logger.Println("Deregister Service Success!")
	return true
}

// DiscoverServices 通过向consul请求服务发现 得到服务实例.
func (consulClient *KitDiscoverClient) DiscoverServices(serviceName string, logger *log.Logger) []interface{} {
	//  先查看本地是否该服务已监控并缓存
	instanceList, ok := consulClient.instancesMap.Load(serviceName)
	if ok {
		// 如果本地有缓存的服务实例,直接返回
		return instanceList.([]interface{})	// 注意 一个服务可以有多个实例
	}
	// 如果本地没有，就从consul中请求
	// 首先申请锁
	consulClient.mutex.Lock()
	defer consulClient.mutex.Unlock()
	// 因为已经加锁,所以再次检查是否监控
	instanceList, ok = consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]interface{})
	} else {
		// 对这个服务进行监控,在有变化时更新缓存
		go func() {
			// 使用 consul 服务实例监控来监控某个服务名的服务实例列表变化
			params := make(map[string]interface{})
			// 根据服务名向Consul注册的Service类型的Watch监控机制
			params["type"] = "service"
			params["service"] = serviceName
			plan, _ := watch.Parse(params)	// 返回一个监控的plan
			// 这个Handler用来对已经被监控且发生变化的服务实例进行处理
			plan.Handler = func(u uint64, i interface{}) {
				if i == nil {
					return
				}
				v, ok := i.([]*api.ServiceEntry)
				if !ok {
					return // 数据异常，忽略
				}
				// 没有服务实例在线
				if len(v) == 0 {
					consulClient.instancesMap.Store(serviceName, []interface{}{})
				}
				var healthServices []interface{}
				for _, service := range v {
					if service.Checks.AggregatedStatus() == api.HealthPassing {
						healthServices = append(healthServices, service.Service)
					}
				}
				consulClient.instancesMap.Store(serviceName, healthServices)
			}
			defer plan.Stop()
			plan.Run(consulClient.config.Address)
		}()
	} // 从这里往上的部分是为了通过缓存减少注册中心负担,也可以直接去掉

	// 根据服务名请求服务实例列表
	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil {
		consulClient.instancesMap.Store(serviceName, []interface{}{})
		logger.Println("Discover Service Error!")
		return nil
	}
	instances := make([]interface{}, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = entries[i].Service
	}
	consulClient.instancesMap.Store(serviceName, instances)
	return instances
}

func NewKitDiscoverClient(consulHost string, consulPort int) (DiscoveryClient, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	apiClient, err := api.NewClient(consulConfig) // 这个client是对于consul api的client
	if err != nil {
		return nil, err
	}
	client := consul.NewClient(apiClient) // 这个才是真正的 consul client
	return &KitDiscoverClient{
		ConsulHost: consulHost,
		ConsulPort: consulPort,
		config:     consulConfig,
		client:     client,
	}, err

}
