package XtremeGin

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yzletter/XtremeGin/middlewares/jwt/jwthandler"
	"github.com/yzletter/XtremeGin/middlewares/jwt/jwtservice"
	"github.com/yzletter/XtremeGin/middlewares/ratelimit"
	limiter "github.com/yzletter/XtremeGin/middlewares/ratelimit/limiter/slide_window_limiter"
	"time"
)

// 快速入门 QuickStart
func main() {
	// redis 数据库
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	// Gin 初始化
	server := gin.Default()

	// 限流服务
	rateLimitHandlerFunc := ratelimit.NewRateLimitBuilder(limiter.NewRedisSlideWindowLimiter(rdb, time.Minute, 10)).Build()

	// JWT服务
	refreshTokenKey := "YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6a"
	accessTokenKey := "YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"
	jh := jwthandler.NewJwtHandler(accessTokenKey, refreshTokenKey, rdb)
	jb := jwtservice.NewJwtServiceBuilder(jh).
		AddIgnorePath("/ping").
		AddIgnorePath("/ping2").
		AddIgnorePath("/ping3").
		AddIgnorePath("/ping4").
		AddIgnorePath("/ping5")
	jwtHandlerFunc := jb.Build()

	// 注册
	server.Use(rateLimitHandlerFunc, jwtHandlerFunc)
}
