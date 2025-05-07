package JWTservice

import (
	"github.com/gin-gonic/gin"
	"github.com/yzletter/go-kit/setx"
)

type JWTBuilder struct {
	ignorePaths *setx.Set[string] // 需要忽略的路由
}

// NewJWTBuilder 构造函数
func NewJWTBuilder() *JWTBuilder {
	hashPaths := setx.NewSet[string]()
	return &JWTBuilder{
		ignorePaths: hashPaths,
	}
}

// Build 返回 Gin 框架中用于注册的 HandlerFunc 中间件，检测当前请求的 JWTToken 是否符合要求
func (builder *JWTBuilder) Build() gin.HandlerFunc {
	// todo

	var checkJWTVaildHandlerFunc func(ctx *gin.Context)

	checkJWTVaildHandlerFunc = func(ctx *gin.Context) {
		// 获取当前路由
		nowPath := ctx.Request.URL.Path
		// 判断是否需要忽略
		if builder.ignorePaths.Exist(nowPath) {
			// 当前路由在需要忽略的路由中
			return
		}
		// 无需忽略，进行鉴权

	}

	return checkJWTVaildHandlerFunc
}

// AddIgnorePath 将 path 加入忽略路径
func (builder *JWTBuilder) AddIgnorePath(path string) {
	builder.ignorePaths.Insert(path)
}
