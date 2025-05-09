package JwtHandler

import "github.com/gin-gonic/gin"

type Handler interface {
	CheckToken(ctx *gin.Context, SSid string) bool
}
