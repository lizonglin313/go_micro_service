package service

import "context"

// SumConcatService 定义服务接口.
// 实现: type Endpoints struct & type sumConcatService struct
type SumConcatService interface {
	Sum(ctx context.Context, a, b int) (reply int)
	Concat(ctx context.Context, a, b string) (reply string)
	HealthCheck(ctx context.Context, status bool) bool
}

type sumConcatService struct {
}

// HealthCheck 健康检查先都返回true.
func (s sumConcatService) HealthCheck(ctx context.Context, status bool) bool {
	return true
}

func (s sumConcatService) Sum(ctx context.Context, a, b int) (reply int) {
	return a + b
}

func (s sumConcatService) Concat(ctx context.Context, a, b string) (reply string) {
	return a + b
}

// NewSumConcatService 返回服务实例.
func NewSumConcatService() SumConcatService {
	return sumConcatService{}
}
