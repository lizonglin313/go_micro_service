package discover

import "log"

// DiscoveryClient 定义与consul进行交互的接口.
type DiscoveryClient interface {
	Register(serviceName, instanceId, healthCheckUrl string,
		instanceHost string, instancePort int,
		meta map[string]string,
		logger *log.Logger) bool

	DeRegister(instanceId string, logger *log.Logger) bool

	DiscoverServices(serviceName string, logger *log.Logger) []interface{}
}
