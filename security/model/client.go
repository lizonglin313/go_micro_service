// Package model 定义客户端信息.
package model

type ClientDetails struct {
	ClientId string	// 客户端标识
	ClientSecret string	// 客户端密钥
	AccessTokenValiditySeconds int	// 访问令牌有效时间 秒
	RefreshTokenValiditySeconds int	// 刷新令牌有效时间 秒
	RegisteredRedirectUri string	// 重定向地址 授权码类型中使用
	AuthorizedGrantTypes []string	// 可以使用的授权类型
}

