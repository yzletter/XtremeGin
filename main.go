package XtremeGin

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yzletter/XtremeGin/middlewares/ratelimit"
	limiter "github.com/yzletter/XtremeGin/middlewares/ratelimit/limiter/slide_window_limiter"
	"time"
)

// QuickStart
func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	r := gin.Default()

	rateLimitHandler := ratelimit.NewRateLimitBuilder(limiter.NewRedisSlideWindowLimiter(rdb, time.Minute, 10)).Build()

	r.Use(rateLimitHandler)
}
