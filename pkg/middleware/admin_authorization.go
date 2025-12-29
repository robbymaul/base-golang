package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	pkgjwt "paymentserviceklink/pkg/jwt"
)

func (a *Auth) AdminAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exists := ctx.Get(BearerToken)
		if !exists {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		response := value.(*pkgjwt.JwtResponse)

		ctx.Set(Session, response)
		ctx.Next()
	}
}
