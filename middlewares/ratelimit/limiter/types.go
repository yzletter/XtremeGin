package limiter

import "context"

// Limiter 每种限流方法，都需要实现一个 Limit 函数
type Limiter interface {
	Limit(ctx context.Context, key string) (bool, error)
}
