package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"io"
	"net/http"
	"paymentserviceklink/app/validate"
	"paymentserviceklink/app/web"
)

func (a *Auth) ClientMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.ResponseWeb{
				Success: false,
				Message: err.Error(),
			})
			ctx.Abort()
			return
		}

		request := new(web.ClientRequest)
		if err = json.Unmarshal(bodyBytes, request); err != nil {
			ctx.JSON(http.StatusBadRequest, web.ResponseWeb{
				Success: false,
				Message: err.Error(),
			})
			ctx.Abort()
			return
		}

		err = validate.ValidationClientRequest(request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.ResponseWeb{
				Success: false,
				Message: err.Error(),
			})
			ctx.Abort()
			return
		}

		// Reset ulang body supaya controller bisa baca lagi
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		platform, err := a.repo.GetPlatformByApiKeyAndSecretKeyRepository(ctx, request)
		if err != nil {
			log.Error().Err(err).Msg("get platform by api key and secret key repository failed")
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusUnauthorized, web.ResponseWeb{
					Success: false,
					Message: fmt.Sprintf("platform api key and secret key is not yet registered"),
				})
				ctx.Abort()
				return
			}

			ctx.JSON(http.StatusInternalServerError, web.ResponseWeb{
				Success: false,
				Message: err.Error(),
			})
			ctx.Abort()
			return
		}

		log.Info().Interface("client", request).Interface("platform", platform).Msg("client request authorization")
		ctx.Set(Client, platform)
		ctx.Next()
	}
}
