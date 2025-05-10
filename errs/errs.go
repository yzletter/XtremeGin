package errs

import "errors"

var (
	ErrSetAccessToken  = errors.New("XtremeGin:设置 AccessToken 错误")
	ErrSetRefreshToken = errors.New("XtremeGin:设置 RefreshToken 错误")
	ErrRedisSetSSid    = errors.New("XtremeGin:Redis 设置 SSid 错误")
)
