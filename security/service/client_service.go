package service

import (
	"context"
	"errors"
	"security/model"
)

var (
	ErrClientNotExist = errors.New("client Id is not exist")
	ErrClientSecret   = errors.New("invalid client secret")
)

type ClientDetailService interface {
	// GetClientDetailByClientId 根据客户端Id和密钥返回客户端信息.
	GetClientDetailByClientId(ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error)
}

// InMemoryClientDetailsService 现将客户端信息内存存储.
type InMemoryClientDetailsService struct {
	clientDetailsDict map[string]*model.ClientDetails
}

// GetClientDetailByClientId 实现客户端服务接口.
func (service *InMemoryClientDetailsService) GetClientDetailByClientId(
	ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error) {
	clientDetails, ok := service.clientDetailsDict[clientId]
	if ok {
		if clientDetails.ClientSecret == clientSecret {	// 判断密钥
			return clientDetails, nil
		} else {
			return nil, ErrClientSecret
		}
	} else {
		return nil, ErrClientNotExist
	}
}

// NewInMemoryClientDetailService 客户端服务的构造方法.
func NewInMemoryClientDetailService(clientDetailsList []*model.ClientDetails) *InMemoryClientDetailsService {
	clientDetailsDict := make(map[string]*model.ClientDetails)

	// 将 List 的数据存入 map: clientId <-> clientDetail
	if clientDetailsList != nil {
		for _, value := range clientDetailsList {
			clientDetailsDict[value.ClientId] = value
		}
	}

	return &InMemoryClientDetailsService{
		clientDetailsDict: clientDetailsDict,
	}
}
