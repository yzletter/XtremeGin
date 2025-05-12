package limiter

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window_script.lua
var luaSlideWindowScript string // luaSlideWindowScript 滑动窗口算法 lua 脚本

// RedisSlideWindowLimiter 基于 Redis 的滑动窗口限流方法
type RedisSlideWindowLimiter struct {
	cmd      redis.Cmdable // cmd RedisClient
	interval time.Duration // interval 窗口大小
	rate     int           // rate 阈值
}

// NewRedisSlideWindowLimiter 构造函数: 需要传入一个 RedisClient 缓存、窗口大小、阈值 ——> 表示 interval 时间内该 IP 最多通过 rate 个请求
func NewRedisSlideWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) *RedisSlideWindowLimiter {
	return &RedisSlideWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

// Limit 具体限流执行
func (limiter *RedisSlideWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	// ARGV 参数
	windowScale := limiter.interval.Milliseconds()
	maxRate := limiter.rate
	nowTime := time.Now().UnixMilli()
	return limiter.cmd.Eval(ctx, luaSlideWindowScript, []string{key}, windowScale, maxRate, nowTime).Bool()
}
