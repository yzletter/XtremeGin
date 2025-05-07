package JWThandler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yzletter/XtremeGin/middlewares/JWT/JWTclaims"
	"net/http"
	"strings"
	"time"
)

type JWTHandler struct {
	refreshTokenKey []byte // 长 Token 秘钥
	accessTokenKey  []byte // 短 Token 秘钥
}

// NewJWTHandler 构造函数
func NewJWTHandler(refreshTokenKey, accessTokenKey string) JWTHandler {
	return JWTHandler{
		refreshTokenKey: []byte(refreshTokenKey),
		accessTokenKey:  []byte(accessTokenKey),
	}
}

// SetRefreshToken 设置长 Token
func (jh JWTHandler) SetRefreshToken(ctx *gin.Context, uid int64) error {
	claims := JWTclaims.UserClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2400)), // 过期时间
			Issuer:    "yzletter",                                           // 签发人
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(jh.refreshTokenKey)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}

	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

// SetAccessToken 设置短 Token
func (jh JWTHandler) SetAccessToken(ctx *gin.Context, uid int64) error {

	claims := JWTclaims.UserClaims{
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 过期时间
			Issuer:    "yzletter",                                         // 签发人
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(jh.accessTokenKey)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}

	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

// ExtractToken 从上下文中获取 JWTToken
func ExtractToken(ctx *gin.Context) string {
	headerString := ctx.GetHeader("Authorization")         // Bearer + 具体的 Token 内容, eg: Bearer aaaaa.bbbbb.ccccc
	headerContents := strings.SplitN(headerString, " ", 2) // 对字符串进行切分取 Token 内容
	if len(headerContents) != 2 {
		return ""
	}
	return headerContents[1]
}
