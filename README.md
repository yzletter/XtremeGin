# XtremeGin

## 文档

### 版本说明

- 1.x.0：为正式版
- 1.1.x：带后缀均为测试版本，不建议使用

## Gin 中间件库

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

	// 注册
	server.Use(rateLimitHandlerFunc)
}

```



## JWT 模块

博客地址：[JWT 模块设计与实现](https://yzletter.notion.site/JWT-Token-1ec89200bcae80acb735e0fb4914cae4?pvs=74)
#### 快速开始

``` go
package XtremeGin

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yzletter/XtremeGin/jwtx"
	"time"
)

// 快速入门 QuickStart
func main() {
	// redis 数据库
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	// Gin 初始化
	server := gin.Default()

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
	server.Use(jwtHandlerFunc)
}


```