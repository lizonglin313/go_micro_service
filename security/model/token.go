package model

import "time"

// OAuth2Token 定义令牌结构.
type OAuth2Token struct {
	RefreshToken *OAuth2Token	// 刷新令牌
	TokenType string	// 令牌类型
	TokenValue string	// 令牌
	ExpiresTime *time.Time	// 令牌过期时间
}

func (oauth2Token *OAuth2Token) IsExpired() bool {
	return oauth2Token.ExpiresTime != nil &&
		oauth2Token.ExpiresTime.Before(time.Now())
}

type OAuth2Details struct {
	Client *ClientDetails
	User *UserDetails
}