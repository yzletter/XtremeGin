package JWTservice

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yzletter/XtremeGin/middlewares/JWT/JWTclaims"
	"github.com/yzletter/XtremeGin/middlewares/JWT/JWThandler"
	"github.com/yzletter/go-kit/setx"
	"net/http"
)

type JWTBuilder struct {
	ignorePaths *setx.Set[string] // 需要忽略的路由
	jwtHandler  JWThandler.JWTHandler
}

// NewJWTBuilder 构造函数
func NewJWTBuilder(jwtHandler JWThandler.JWTHandler) *JWTBuilder {
	hashPaths := setx.NewSet[string]()
	return &JWTBuilder{
		ignorePaths: hashPaths,
		jwtHandler:  jwtHandler,
	}
}

// Build 返回 Gin 框架中用于注册的 HandlerFunc 中间件，检测当前请求头部所带的 JWTToken 是否符合要求
func (builder *JWTBuilder) Build() gin.HandlerFunc {
	var checkJWTHandlerFunc func(ctx *gin.Context)

	checkJWTHandlerFunc =
		func(ctx *gin.Context) {
			// 获取当前路由
			nowPath := ctx.Request.URL.Path
			if builder.ignorePaths.Exist(nowPath) { // 判断当前路由是否需要忽略
				return
			}
			// 无需忽略，进行鉴权
			claims := &JWTclaims.UserClaims{}
			tokenStr := JWThandler.ExtractToken(ctx)                 // 取出上下文中的 token
			keyFunc := func(token *jwt.Token) (interface{}, error) { // 秘钥函数
				return builder.jwtHandler.SetAccessToken, nil
			}

			// 解析 Token 字符串, 将结果存到 claims 中, 并返回 Token
			token, err := jwt.ParseWithClaims(tokenStr, claims, keyFunc)

			if err != nil || !token.Valid || token == nil || claims.Uid == 0 {
				// 当前请求的用户没登录
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if claims.UserAgent != ctx.Request.UserAgent() {
				// 严重的安全问题，你是要监控
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			ctx.Set("claims", claims) // 将当前请求的用户数据存入上下文
		}

	return checkJWTHandlerFunc
}

// AddIgnorePath 将 path 加入忽略路径
func (builder *JWTBuilder) AddIgnorePath(path string) {
	builder.ignorePaths.Insert(path)
}
