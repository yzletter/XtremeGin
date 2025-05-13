package jwtx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yzletter/XtremeGin/errs"
	"net/http"
	"strings"
	"time"
)

type JwtHandler struct {
	AccessTokenKey       []byte        `AccessToken 秘钥`
	RefreshTokenKey      []byte        `RefreshToken 秘钥`
	AccessTokenDuration  time.Duration `AccessToken 过期时间`
	RefreshTokenDuration time.Duration `RefreshToken 过期时间`
	AccessTokenHeader    string        `AccessToken 请求头名`
	RefreshTokenHeader   string        `RefreshToken 请求头名`
	CtxClaimsName        string        `CTX 存储用户信息的 Claims 名 : claims`
	IssuerName           string        `JWT 签名人 : yzletter`
	RedisKeyPrefix       string        `Redis Key 前缀 : users:ssid `
	RedisClient          redis.Cmdable `Redis 缓存`
}

// NewJwtHandler 构造函数
func NewJwtHandler(config HandlerConfig, RedisClient redis.Cmdable) *JwtHandler {
	return &JwtHandler{
		AccessTokenKey:       config.AccessTokenKey,
		RefreshTokenKey:      config.RefreshTokenKey,
		AccessTokenDuration:  config.AccessTokenDuration,
		RefreshTokenDuration: config.RefreshTokenDuration,
		AccessTokenHeader:    config.AccessTokenHeader,
		RefreshTokenHeader:   config.RefreshTokenHeader,
		CtxClaimsName:        config.CtxClaimsName,
		IssuerName:           config.IssuerName,
		RedisKeyPrefix:       config.RedisKeyPrefix,
		RedisClient:          RedisClient,
	}
}

// CheckTokenDiscarded 判断 Token 是否被废弃
func (jh *JwtHandler) CheckTokenDiscarded(ctx *gin.Context, SSid string) bool {
	cnt, _ := jh.RedisClient.Exists(ctx, MakeRedisKey(jh.RedisKeyPrefix, SSid)).Result()
	if cnt > 0 {
		return true
	}
	return false
}

// SetLoginToken 设置登录 Token, 包括长 token 和短 token
func (jh *JwtHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := jh.SetRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = jh.SetAccessToken(ctx, uid, ssid)
	return err
}

// SetRefreshToken 设置长 token
func (jh *JwtHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	// 1. 携带信息的声明
	myClaims := &RefreshClaims{
		Uid:  uid,
		SSid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jh.IssuerName,                                               // 签名人
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jh.RefreshTokenDuration)), // 过期时间
		},
	}
	// 2. 生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, myClaims)
	// 3. 对 Token 进行加密
	tokenString, err := token.SignedString(jh.RefreshTokenKey) // 用长 token 秘钥进行加密
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return errs.ErrSetRefreshToken
	}
	// 4. 将 token 放入上下文
	ctx.Header(jh.RefreshTokenHeader, tokenString)
	return nil
}

// SetAccessToken 设置短 token
func (jh *JwtHandler) SetAccessToken(ctx *gin.Context, uid int64, ssid string) error {
	// 1. 携带信息的声明
	myClaims := &AccessClaims{
		Uid:       uid,
		SSid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jh.IssuerName,                                              // 签名人
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jh.AccessTokenDuration)), // 过期时间
		},
	}
	// 2. 生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, myClaims)
	// 3. 对 Token 进行加密
	tokenString, err := token.SignedString(jh.AccessTokenKey) // 用短 token 秘钥进行加密
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return errs.ErrSetAccessToken
	}
	// 4. 将 token 放入上下文
	ctx.Header(jh.AccessTokenHeader, tokenString)
	return nil
}

// ClearToken 将 token 废弃
func (jh *JwtHandler) ClearToken(ctx *gin.Context) error {
	// 1. 设置前端请求头的长短 token 为非法值
	SetHeaderNil(ctx, jh)
	// 2. 获取当前请求的 SSid
	myClaims := ctx.MustGet(jh.CtxClaimsName).(*AccessClaims)
	// 3. 在 Redis 中记录当前 SSid 废弃
	err := jh.RedisClient.Set(ctx, MakeRedisKey(jh.RedisKeyPrefix, myClaims.SSid), "", jh.RefreshTokenDuration).Err()
	if err != nil {
		return errs.ErrRedisSetSSid
	}
	return nil
}

// RefreshAccessToken 若长 token 未过期, 则刷新短 token
func (jh *JwtHandler) RefreshAccessToken(ctx *gin.Context) {
	// 1. 取出 tokenString
	tokenString := ExtractToken(ctx, "Authorization")
	// 2. 用 refreshTokenKey 解析到声明中
	targetClaims := &RefreshClaims{}
	keyFunc := MakeKeyFunc(jh.RefreshTokenKey)
	token, err := jwt.ParseWithClaims(tokenString, targetClaims, keyFunc)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 3. 判断 token 是否废弃
	if ifDiscarded := jh.CheckTokenDiscarded(ctx, targetClaims.SSid); ifDiscarded {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 4. 判断 token 是否正确
	if !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 5. 设置新的 access token
	err = jh.SetAccessToken(ctx, targetClaims.Uid, targetClaims.SSid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 6. 返回
	ctx.String(http.StatusOK, "刷新 accessToken 成功")
}

// SetHeaderNil 将 JWT 请求头设为空
func SetHeaderNil(ctx *gin.Context, jh *JwtHandler) {
	ctx.Header(jh.AccessTokenHeader, "")
	ctx.Header(jh.RefreshTokenHeader, "")
}

// ExtractToken 从上下文取出 tokenString
func ExtractToken(ctx *gin.Context, HeaderName string) string {
	// "Authorization"
	headerString := ctx.GetHeader(HeaderName)
	headerStringSeg := strings.SplitN(headerString, " ", 2)
	if len(headerStringSeg) != 2 {
		return ""
	}
	return headerStringSeg[1]
}

// MakeRedisKey 返回用于 Redis 查询的 Key
func MakeRedisKey(prefix, SSid string) string {
	return fmt.Sprintf("%s:%s", prefix, SSid)
}

// MakeKeyFunc 返回用于解析 JWT Token 的函数
func MakeKeyFunc(RefreshTokenKey []byte) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		return RefreshTokenKey, nil
	}
}
