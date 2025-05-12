package XtremeGin

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yzletter/XtremeGin/middlewares/jwtx"
	"github.com/yzletter/XtremeGin/middlewares/ratelimitx"
	limiter "github.com/yzletter/XtremeGin/middlewares/ratelimitx/limiter/slide_window_limiter"
	"time"
)

// 快速入门 QuickStart
func main() {
	// redis 数据库
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	// Gin 初始化
	server := gin.Default()

	// 限流服务
	rateLimitHandlerFunc := ratelimitx.NewRateLimitBuilder(limiter.NewRedisSlideWindowLimiter(rdb, time.Minute, 10)).Build()

	// JWT服务
	RefreshTokenKey := "YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6a"
	AccessTokenKey := "YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"
	CtxClaimsName := "myClaims"
	IssuerName := "yzletter"
	AccessTokenDuration := time.Hour * 24 * 7
	RefreshTokenDuration := time.Hour * 24
	AccessTokenHeader := "x-access-token"
	RefreshTokenHeader := "x-refresh-token"
	RedisKeyPrefix := "users:ssid"

	jh := jwtx.
		NewJwtHandler(AccessTokenKey, RefreshTokenKey,
			CtxClaimsName, IssuerName, RedisKeyPrefix,
			AccessTokenHeader, RefreshTokenHeader,
			AccessTokenDuration, RefreshTokenDuration,
			rdb)
	jb := jwtx.NewJwtServiceBuilder(jh).
		AddIgnorePath("/ping").
		AddIgnorePath("/ping2").
		AddIgnorePath("/ping3").
		AddIgnorePath("/ping4").
		AddIgnorePath("/ping5")
	jwtHandlerFunc := jb.Build()

	// 注册
	server.Use(rateLimitHandlerFunc, jwtHandlerFunc)
}
