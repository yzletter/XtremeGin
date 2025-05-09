package JwtHandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type JwtHandler struct {
	AccessTokenKey  []byte
	RefreshTokenKey []byte
	RedisClient     redis.Cmdable
}

func NewJwtHandler() *JwtHandler {
	return &JwtHandler{}
}

// CheckToken 判断 Token 是否被废弃
func (jh *JwtHandler) CheckToken(ctx *gin.Context, SSid string) bool {
	cnt, _ := jh.RedisClient.Exists(ctx, fmt.Sprintf("users:ssid:%s", SSid)).Result()
	if cnt > 0 {
		return true
	}
	return false
}

// ExtractToken 从上下文取出 Token
func ExtractToken(ctx *gin.Context) string {
	return ""
}
