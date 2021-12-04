// Package discover 定义和实现了与consul交互的接口.
package discover

import (
	"log"
)

// DiscoveryClient 包括服务注册，发现和注销 用于 与consul进行交互的接口.
type DiscoveryClient interface {

	// Register 向consul注册服务实例.
	Register(serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int,
		meta map[string]string, logger *log.Logger) bool

	// DeRegister 向consul注销服务实例.
	DeRegister(instanceId string, logger *log.Logger) bool

	// DiscoverServices 通过服务名向consul请求发现某个服务.
	DiscoverServices(serviceName string, logger *log.Logger) []interface{}
}

