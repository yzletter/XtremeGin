package XtremeGin

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yzletter/XtremeGin/middlewares/ratelimit"
	limiter "github.com/yzletter/XtremeGin/middlewares/ratelimit/limiter/slide_window_limiter"
	"time"
)

// 快速入门 QuickStart
func main() {
	// redis
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	// 服务端
	r := gin.Default()

	// 限流
	rateLimitHandlerFunc := ratelimit.NewRateLimitBuilder(limiter.NewRedisSlideWindowLimiter(rdb, time.Minute, 10)).Build()

	r.Use(rateLimitHandlerFunc)
}
