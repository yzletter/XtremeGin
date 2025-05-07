package JWTclaims

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	Uid       int64  // 要放进去 token 里面的数据
	UserAgent string // 请求
	jwt.RegisteredClaims
}
