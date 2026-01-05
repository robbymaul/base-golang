package middleware

import (
	"github.com/gin-gonic/gin"
)

func (a *Auth) ClientMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Next()
	}
}
