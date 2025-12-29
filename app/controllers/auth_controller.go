package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/services"
	"paymentserviceklink/app/validate"
	"paymentserviceklink/app/web"
	pkgjwt "paymentserviceklink/pkg/jwt"
	"paymentserviceklink/pkg/middleware"
)

func (c *Controller) AdminLoginController(ctx *gin.Context) {
	c.context = "AdminAuth.AdminLoginController"
	request := new(web.AdminLoginRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, c.context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationAdminLoginRequest(request)
	if err != nil {
		log.Debug().Err(err).Msg("should validate admin login request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, c.context).WithStatusCode(http.StatusBadRequest))
		return
	}

	healthService := services.NewAdminService(ctx, c.repo, c.cfg)

	response, err := healthService.AdminLoginService(request)
	if err != nil {
		log.Debug().Err(err).Msg("admin login service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminMeController(ctx *gin.Context) {
	c.context = "AdminAuth.AdminLoginController"

	adminService := services.NewAdminService(ctx, c.repo, c.cfg)

	response, err := adminService.GetDetailAdmin()
	if err != nil {
		log.Debug().Err(err).Msg("get detail admin error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) GetAdminSession(ctx context.Context) (*pkgjwt.JwtResponse, error) {
	c.context = "Auth.GetAdminSession"

	log.Debug().Ctx(ctx).Interface("context", c.context).Msg("get admin session")

	value := ctx.Value(middleware.Session)
	if value == nil {
		return nil, helpers.NewErrorTrace(fmt.Errorf("unauthorized"), c.context).WithStatusCode(http.StatusUnauthorized)
	}
	log.Debug().Interface("value", value).Interface("context", c.context).Msg("get admin session")

	session, ok := value.(*pkgjwt.JwtResponse)
	if !ok {
		return nil, helpers.NewErrorTrace(fmt.Errorf("unauthorized"), c.context).WithStatusCode(http.StatusUnauthorized)
	}
	log.Debug().Interface("ok", ok).Interface("session", session).Interface("context", c.context).Msg("value context jwt response")

	return session, nil
}
