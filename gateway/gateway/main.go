package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {

}

// NewReverseProxy 创建反向代理.
// 实现的效果如下：
// /stringService/op/diff/a/b -> http:127.0.0.1:8123/op/diff/a/b
func NewReverseProxy(client *api.Client, logger log.Logger) *httputil.ReverseProxy {

	// 创建 director 转发 http 请求
	director := func(req *http.Request) {
		// 查询原始请求路径
		reqPath := req.URL.Path
		if reqPath == "" {
			logger.Log("requestUrl", "空url")
			return
		}

		// 划分url结构
		pathArray := strings.Split(reqPath, "/")
		serviceName := pathArray[1] // 拿到服务名

		// 调用consul api 查询服务实例列表
		result, _, err := client.Catalog().Service(serviceName, "", nil)
		if err != nil {
			logger.Log("ReverseProxy failed:", "query service instance error", err.Error())
			return
		}

		// 如果没有服务实例
		if len(result) == 0 {
			logger.Log("ReverseProxy failed:", "no such service instance", serviceName)
			return
		}

		// 重新组织请求路径 去掉服务名部分
		destPath := strings.Join(pathArray[2:], "/")

		// 随机选择服务实例
		tgt := result[rand.Int()%len(result)]
		logger.Log("service id", tgt.ServiceID)

		// 设置代理服务地址 也就是新的转发地址
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("%s:%d", tgt.ServiceAddress, tgt.ServicePort)
		req.URL.Path = "/" + destPath
	}
	return &httputil.ReverseProxy{Director: director}
}
