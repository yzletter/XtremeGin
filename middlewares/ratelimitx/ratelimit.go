package ratelimitx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yzletter/XtremeGin/middlewares/ratelimitx/limiter"
	"net/http"
)

// RateLimitBuilder 限流器
type RateLimitBuilder struct {
	prefix  string          // prefix 用于构造 Redis 中的 key
	limiter limiter.Limiter // limiter 为采用的限流策略
}

// NewRateLimitBuilder 构造函数, 所采用的 limiter 由外部依赖注入
func NewRateLimitBuilder(limiter limiter.Limiter) *RateLimitBuilder {
	return &RateLimitBuilder{
		prefix:  "ip-limiter",
		limiter: limiter,
	}
}

// Build 返回 Gin 框架中用于注册的 HandlerFunc 中间件
func (builder *RateLimitBuilder) Build() gin.HandlerFunc {
	// 要返回的 HandlerFunc
	var rateLimitGinHandlerFunc func(ctx *gin.Context)

	rateLimitGinHandlerFunc =
		func(ctx *gin.Context) {
			ifLimited, err := builder.limit(ctx)
			if err != nil {
				// 这边记录日志，限流出错（一般为 Redis 出错）
				// 处理事故：
				// 保守做法：由于借助 Redis 进行限流, Redis 崩溃了, 为了防止系统崩溃，直接限流 ——> ctx.AbortWithStatus(http.StatusInternalServerError)
				// 激进做法：虽然 Redis 崩溃了, 但为了不影响用户体验，直接放行 ——> ctx.Next()
				ctx.AbortWithStatus(http.StatusInternalServerError) // 这里采用保守做法
				return
			}

			if ifLimited { // 当前访问被限流了
				ctx.AbortWithStatus(http.StatusTooManyRequests)
				return
			}
			ctx.Next()
		}

	return rateLimitGinHandlerFunc
}

// limit 通过滑动窗口算法 lua 脚本执行限流, 若进行限流, 返回 true
func (builder *RateLimitBuilder) limit(ctx *gin.Context) (bool, error) {
	redisKey := fmt.Sprintf("%s:%s", builder.prefix, ctx.ClientIP()) // prefix + IP 构成 Redis 中的 key

	return builder.limiter.Limit(ctx, redisKey)
}

// UpdatePrefix 修改 builder 的 prefix
func (builder *RateLimitBuilder) UpdatePrefix(prefix string) {
	builder.prefix = prefix
}
