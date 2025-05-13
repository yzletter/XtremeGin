package XtremeGin

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yzletter/XtremeGin/jwtx"
	"github.com/yzletter/XtremeGin/ratelimitx"
	limiterx "github.com/yzletter/XtremeGin/ratelimitx/limiterx/slide_window_limiter"
	"time"
)

// 快速入门 QuickStart
func main() {
	// redis 数据库
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	// Gin 初始化
	server := gin.Default()

	// 限流服务
	rateLimitHandlerFunc := ratelimitx.NewRateLimitBuilder(limiterx.NewRedisSlideWindowLimiter(redisClient, time.Minute, 10)).Build()

	// JWT 服务
	handlerConfig := &jwtx.HandlerConfig{
		AccessTokenKey:  []byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6a"), // AccessToken 秘钥
		RefreshTokenKey: []byte("YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"), // RefreshToken 秘钥
	}

	jwtHandlerFunc := jwtx.NewJwtServiceBuilder(jwtx.NewJwtHandler(handlerConfig, redisClient), false).
		AddIgnorePath("/ping").
		AddIgnorePath("/ping2").
		AddIgnorePath("/ping3").
		AddIgnorePath("/ping4").
		AddIgnorePath("/ping5").Build()

	// 注册
	server.Use(rateLimitHandlerFunc, jwtHandlerFunc)
}
