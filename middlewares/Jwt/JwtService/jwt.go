package JwtService

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yzletter/XtremeGin/middlewares/Jwt/JwtHandler"
	"github.com/yzletter/XtremeGin/middlewares/Jwt/Jwtclaims"
	"github.com/yzletter/go-kit/setx"
)

type JwtServiceBuilder struct {
	IgnorePaths *setx.Set[string]
	JwtHandler  JwtHandler.JwtHandler
}

// NewJwtServiceBuilder 构造函数
func NewJwtServiceBuilder(jwtHandler JwtHandler.JwtHandler) *JwtServiceBuilder {
	return &JwtServiceBuilder{
		IgnorePaths: setx.NewSet[string](),
		JwtHandler:  jwtHandler,
	}
}

// AddIgnorePath 添加需要忽略不进行鉴权的路径
func (jb *JwtServiceBuilder) AddIgnorePath(path string) *JwtServiceBuilder {
	jb.IgnorePaths.Insert(path)
	return jb
}

// Build 返回用于 Gin 注册的 gin.HandlerFunc
func (jb *JwtServiceBuilder) Build() gin.HandlerFunc {
	// todo 当前函数未进行错误处理
	var CheckRequestJWT func(ctx *gin.Context)

	CheckRequestJWT = func(ctx *gin.Context) {
		// 1. 判断当前路径是否需要忽略
		requestRoute := ctx.Request.URL.Path
		if jb.IgnorePaths.Exist(requestRoute) {
			return
		}
		// 2. 取出 tokenString
		tokenString := JwtHandler.ExtractToken(ctx)
		// 3. 将 Token 解析成 Claims
		targetAccessClaims := &Jwtclaims.AccessClaims{}
		keyFunc := func(token *jwt.Token) (interface{}, error) {
			return jb.JwtHandler.AccessTokenKey, nil
		}
		token, _ := jwt.ParseWithClaims(tokenString, targetAccessClaims, keyFunc)
		// 4. 判断 Token 是否正确
		if !token.Valid {
			return
		}
		// 5. 判断 Token 是否废弃
		if ifDiscarded := jb.JwtHandler.CheckTokenDiscarded(ctx, targetAccessClaims.SSid); ifDiscarded {
			return
		}
		// 6. 将请求的用户数据存入上下文
		ctx.Set("claims", targetAccessClaims)
		// 7. 退出
		return
	}

	return CheckRequestJWT
}
