package jwthandler

import "github.com/gin-gonic/gin"

type Handler interface {
	CheckTokenDiscarded(ctx *gin.Context, SSid string) bool
	ClearToken(ctx *gin.Context) error
	SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
	SetAccessToken(ctx *gin.Context, uid int64, ssid string) error
	SetLoginToken(ctx *gin.Context, uid int64) error
	RefreshAccessToken(ctx *gin.Context)
}
