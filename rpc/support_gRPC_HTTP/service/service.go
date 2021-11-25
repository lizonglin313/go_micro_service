package service

import "context"

// SumConcatService 定义服务接口.
type SumConcatService interface {
	Sum(ctx context.Context, a, b int) (reply int)
	Concat(ctx context.Context, a, b string) (reply string)
}

type sumConcatService struct {
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
