package jwtx

import "time"

type HandlerConfig struct {
	AccessTokenKey       []byte        `AccessToken 秘钥`
	RefreshTokenKey      []byte        `RefreshToken 秘钥`
	AccessTokenDuration  time.Duration `AccessToken 过期时间`
	RefreshTokenDuration time.Duration `RefreshToken 过期时间`
	AccessTokenHeader    string        `AccessToken 请求头名`
	RefreshTokenHeader   string        `RefreshToken 请求头名`
	AuthorizationHeader  string        `认证信息请求头`
	CtxClaimsName        string        `CTX 存储用户信息的 Claims 名 : claims`
	IssuerName           string        `JWT 签名人 : yzletter`
	RedisKeyPrefix       string        `Redis Key 前缀 : users:ssid `
}

func (config *HandlerConfig) init() {
	if config.AccessTokenDuration == 0 {
		config.AccessTokenDuration = time.Hour * 24
	}
	if config.RefreshTokenDuration == 0 {
		config.RefreshTokenDuration = time.Hour * 24 * 7
	}
	if config.AccessTokenHeader == "" {
		config.AccessTokenHeader = "x-access-token"
	}
	if config.RefreshTokenHeader == "" {
		config.RefreshTokenHeader = "x-refresh-token"
	}
	if config.AuthorizationHeader == "" {
		config.AuthorizationHeader = "Authorization"
	}
	if config.CtxClaimsName == "" {
		config.AuthorizationHeader = "myClaims"
	}
	if config.IssuerName == "" {
		config.AuthorizationHeader = "yzletter"
	}
	if config.RedisKeyPrefix == "" {
		config.AuthorizationHeader = "users:ssid"
	}
}
