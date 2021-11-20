package config

import (
	kitlog "github.com/go-kit/kit/log"
	"log"
	"os"
)


var Logger *log.Logger
var KitLogger kitlog.Logger


func init() {
	// 标准的向控制台输出的log.
	Logger = log.New(os.Stderr, "", log.LstdFlags)

	// 设置kit日志.
	KitLogger = kitlog.NewLogfmtLogger(os.Stderr)
	KitLogger = kitlog.With(KitLogger, "ts", kitlog.DefaultTimestampUTC)
	KitLogger = kitlog.With(KitLogger, "caller", kitlog.DefaultCaller)

}

