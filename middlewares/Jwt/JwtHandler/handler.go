package JwtHandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yzletter/XtremeGin/middlewares/Jwt/Jwtclaims"
	"net/http"
	"strings"
	"time"
)

type JwtHandler struct {
	AccessTokenKey  []byte
	RefreshTokenKey []byte
	RedisClient     redis.Cmdable
}

func NewJwtHandler(accessTokenKey, refreshTokenKey string, redisClient redis.Cmdable) *JwtHandler {
	return &JwtHandler{
		AccessTokenKey:  []byte(accessTokenKey),
		RefreshTokenKey: []byte(refreshTokenKey),
		RedisClient:     redisClient,
	}
}

// CheckToken 判断 Token 是否被废弃
func (jh *JwtHandler) CheckTokenDiscarded(ctx *gin.Context, SSid string) bool {
	cnt, _ := jh.RedisClient.Exists(ctx, fmt.Sprintf("users:ssid:%s", SSid)).Result()
	if cnt > 0 {
		return true
	}
	return false
}

// SetLoginToken 设置登录 Token, 包括长 token 和短 token
func (jh *JwtHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	// todo 未进行错误处理
	ssid := uuid.New().String()
	jh.SetRefreshToken(ctx, uid, ssid)
	jh.SetAccessToken(ctx, uid, ssid)
	return nil
}

// SetRefreshToken 设置长 token
func (jh *JwtHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	// todo 未进行错误处理
	// 1. 携带信息的声明
	myClaims := &Jwtclaims.RefreshClaims{
		Uid:  uid,
		SSid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "yzletter",                                             // 签名人
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 过期时间
		},
	}
	// 2. 生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, myClaims)
	// 3. 对 Token 进行加密
	tokenString, _ := token.SignedString(jh.RefreshTokenKey) // 用长 token 秘钥进行加密
	// 4. 将 token 放入上下文
	ctx.Set("x-refresh-token", tokenString)
	return nil
}

// SetAccessToken 设置短 token
func (jh *JwtHandler) SetAccessToken(ctx *gin.Context, uid int64, ssid string) error {
	// todo 未进行错误处理
	// 1. 携带信息的声明
	myClaims := &Jwtclaims.AccessClaims{
		Uid:       uid,
		SSid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "yzletter",                                         // 签名人
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 过期时间
		},
	}
	// 2. 生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, myClaims)
	// 3. 对 Token 进行加密
	tokenString, _ := token.SignedString(jh.AccessTokenKey) // 用短 token 秘钥进行加密
	// 4. 将 token 放入上下文
	ctx.Set("x-access-token", tokenString)
	return nil
}

// ClearToken 将 token 废弃
func (jh *JwtHandler) ClearToken() {

}

// RefreshAccessToken 若长 token 未过期, 则刷新短 token
func (jh *JwtHandler) RefreshAccessToken(ctx *gin.Context) {
	// todo 未进行错误处理
	// 1. 取出 tokenString
	tokenString := ExtractToken(ctx)
	// 2. 用 refreshTokenKey 解析到声明中
	targetClaims := &Jwtclaims.RefreshClaims{}
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		return jh.RefreshTokenKey, nil
	}
	token, _ := jwt.ParseWithClaims(tokenString, targetClaims, keyfunc)
	// 3. 判断 token 是否废弃
	if ifDiscarded := jh.CheckTokenDiscarded(ctx, targetClaims.SSid); ifDiscarded {
		return
	}
	// 4. 判断 token 是否正确
	if !token.Valid {
		return
	}
	// 5. 设置新的 access token
	jh.SetAccessToken(ctx, targetClaims.Uid, targetClaims.SSid)
	// 6. 返回
	ctx.String(http.StatusOK, "刷新短token成功")
}

// ExtractToken 从上下文取出 tokenString
func ExtractToken(ctx *gin.Context) string {
	// todo 未进行错误处理
	headerString := ctx.GetHeader("Authorization")
	headerStringSegs := strings.SplitN(headerString, " ", 2)
	return headerStringSegs[1]
}
