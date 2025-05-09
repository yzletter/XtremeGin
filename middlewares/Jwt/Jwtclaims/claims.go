package Jwtclaims

import "github.com/golang-jwt/jwt/v5"

// RefreshClaims 长 Token 声明
type RefreshClaims struct {
	Uid  int64
	SSid string
	jwt.RegisteredClaims
}

// AccessClaims 短 Token 声明
type AccessClaims struct {
	Uid       int64
	SSid      string
	UserAgent string
	jwt.RegisteredClaims
}
