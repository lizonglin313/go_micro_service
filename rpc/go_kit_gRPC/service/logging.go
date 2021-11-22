package service

import (
	"context"
	"github.com/go-kit/log"
	"time"
)

// loggingMiddleware 接入服务的中间件.
// 为什么logging中间件也要实现接口呢？
// 因为构造 Endpoint 需要一个 Service 类型的参数
type loggingMiddleware struct {
	Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{next, logger}
	}
}

func (mw loggingMiddleware) Concat(ctx context.Context, a, b string) (ret string, err error) {
	// 接入中间件功能
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Concat",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now()) // 直接调用

	// 再去执行真正的功能
	ret, err = mw.Service.Concat(ctx, a, b)
	return
}

func (mw loggingMiddleware) Diff(ctx context.Context, a, b string) (ret string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Diff",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.Diff(ctx, a, b)
	return
}
