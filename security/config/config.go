package config

import (
	kitlog "github.com/go-kit/kit/log"
	"log"
	"os"
)

var Logger *log.Logger
var KitLogger kitlog.Logger		// Logger 事实上是一个接口

// init 方法 在使用该包之前 自动调用.
func init() {
	// 设置系统 log 输出方式等信息
	Logger = log.New(os.Stderr, "", log.LstdFlags)

	// 设置kit log 输出方式等信息
	KitLogger = kitlog.NewLogfmtLogger(os.Stderr)
	KitLogger = kitlog.With(KitLogger, "ts", kitlog.DefaultTimestampUTC)
	KitLogger = kitlog.With(KitLogger, "caller", kitlog.DefaultCaller)
}

