package XtremeGin

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	JWTservice "github.com/yzletter/XtremeGin/middlewares/JWT"
	"github.com/yzletter/XtremeGin/middlewares/JWT/JWThandler"
	"github.com/yzletter/XtremeGin/middlewares/ratelimit"
	limiter "github.com/yzletter/XtremeGin/middlewares/ratelimit/limiter/slide_window_limiter"
	"time"
)

// QuickStart
func main() {
	// redis
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	r := gin.Default()

	// 限流
	rateLimitHandlerFunc := ratelimit.NewRateLimitBuilder(limiter.NewRedisSlideWindowLimiter(rdb, time.Minute, 10)).Build()

	refreshTokenKey := "YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6T"
	accessTokenKey := "YTsKHvuxjcQ3jGXrSXH27JvnA3XTkJ6a"
	jwtHandler := JWThandler.NewJWTHandler(refreshTokenKey, accessTokenKey)
	jwtHandlerFunc := JWTservice.NewJWTBuilder(jwtHandler).Build()

	r.Use(rateLimitHandlerFunc, jwtHandlerFunc)
}
