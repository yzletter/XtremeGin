# XtremeGin

Gin 中间件库

- 限流模块
- JWT 模块

## 限流模块
博客地址：[限流模块设计与实现](https://yzletter.notion.site/Gin-Redis-1df89200bcae80b4802dfe20582bda46)

### 基于 Redis 的滑动窗口限流

#### 快速开始

``` go
package XtremeGin

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yzletter/XtremeGin/middlewares/ratelimit"
	limiter "github.com/yzletter/XtremeGin/middlewares/ratelimit/limiter/slide_window_limiter"
	"time"
)

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	r := gin.Default()
	
	rateLimitHandler := ratelimit.NewRateLimitBuilder(limiter.NewRedisSlideWindowLimiter(rdb, time.Minute, 10)).Build()

	r.Use(rateLimitHandler)
}
```



## JWT 模块

博客地址：[JWT 模块设计与实现](https://yzletter.notion.site/JWT-Token-1ec89200bcae80acb735e0fb4914cae4?pvs=74)
