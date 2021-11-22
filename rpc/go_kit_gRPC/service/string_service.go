package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const StrMaxSize = 1024

var (
	ErrMaxSize = errors.New("maximum size of 1024 bytes exceeded")

	ErrStrValue = errors.New("maximum size of 1024 bytes exceeded")
)

// Service 定义服务接口 任何实现这两个功能的服务都是 String 服务.
type Service interface {
	Concat(ctx context.Context, a, b string) (string, error)

	Diff(ctx context.Context, a, b string) (string, error)
}

// StringService 实现上述接口.
type StringService struct {
}

func (s StringService) Concat(ctx context.Context, a, b string) (string, error) {
	if len(a)+len(b) > StrMaxSize {
		return "", ErrMaxSize
	}
	fmt.Printf("StringService Concat return %s", a+b)
	return a + b, nil
}

func (s StringService) Diff(ctx context.Context, a, b string) (string, error) {
	if len(a) < 1 || len(b) < 1 {
		return "", nil
	}
	res := ""
	if len(a) >= len(b) {
		for _, char := range b {
			if strings.Contains(a, string(char)) {
				res = res + string(char)
			}
		}
	} else {
		for _, char := range a {
			if strings.Contains(b, string(char)) {
				res = res + string(char)
			}
		}
	}
	return res, nil
}

// ServiceMiddleware 用于接入日志中间件 还是很形象的.
// 它接收一个服务 进行某些处理后（中间过程）再把服务交出
// 进行后续的服务处理
type ServiceMiddleware func(Service) Service