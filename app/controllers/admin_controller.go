package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/services"
	"paymentserviceklink/app/validate"
	"paymentserviceklink/app/web"
	"strconv"
)

func (c *Controller) CreateAdminRoleController(ctx *gin.Context) {
	context := "AdminRole.CreateAdminRole"

	request := new(web.CreateAdminRolesRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationCreateAdminRolesRequest(request)
	if err != nil {
		log.Debug().Err(err).Msg("should validate create admin roles request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	adminService := services.NewAdminService(ctx, c.repo, c.cfg)

	err = adminService.CreateAdminRoleService(request)
	if err != nil {
		log.Error().Err(err).Msg("create admin role service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
	})
}

func (c *Controller) GetListAdminRoleController(ctx *gin.Context) {
	context := "AdminRole.GetListAdminRole"

	//page := ctx.DefaultQuery("page", "1")
	//perPage := ctx.DefaultQuery("per_page", "10")
	////search := ctx.DefaultQuery("search", "")
	//order := ctx.DefaultQuery("order", "asc")
	//isDelete := ctx.DefaultQuery("is_deleted", "false")
	//filter := ctx.Query("filter")
	//joinOperator := ctx.DefaultQuery("join_operator", "and")
	//
	//pages, err := pagination.New(page, perPage, 0, order, isDelete, filter, joinOperator)

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, "").WithStatusCode(http.StatusBadRequest))
		return
	}

	adminService := services.NewAdminService(ctx, c.repo, c.cfg)

	response, err := adminService.GetListAdminRoleService(pages)
	if err != nil {
		log.Error().Err(err).Msg("get list admin role service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) GetDetailAdminRoleController(ctx *gin.Context) {
	context := "AdminRole.GetDetailAdminRole"
	id := ctx.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	adminService := services.NewAdminService(ctx, c.repo, c.cfg)

	response, err := adminService.GetDetailAdminRoleService(idInt)
	if err != nil {
		log.Error().Err(err).Msg("get admin role failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) UpdateAdminRoleController(ctx *gin.Context) {
	context := "AdminRole.UpdateAdminRole"

	request := new(web.AdminRoleResponse)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationUpdateRoleRequest(request)
	if err != nil {
		log.Debug().Err(err).Msg("should validate update admin role request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	adminService := services.NewAdminService(ctx, c.repo, c.cfg)

	response, err := adminService.UpdateAdminRoleService(request)
	if err != nil {
		log.Error().Err(err).Msg("update admin role service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) CreateAdminUserController(ctx *gin.Context) {
	context := "AdminUser.CreateAdminUser"

	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	request := new(web.CreateAdminUserRequest)

	err = ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationCreateAdminUserRequest(request)
	if err != nil {
		log.Debug().Err(err).Msg("should validate create admin user request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	adminService := services.NewAdminService(ctx, c.repo, c.cfg)

	response, err := adminService.CreateAdminUserService(session, request)
	if err != nil {
		log.Error().Err(err).Msg("create admin user service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})

}

func (c *Controller) GetListAdminUserController(ctx *gin.Context) {
	context := "AdminUser.GetListAdminUser"
	//page := ctx.DefaultQuery("page", "1")
	//perPage := ctx.DefaultQuery("per_page", "10")
	//search := ctx.DefaultQuery("search", "")
	//order := ctx.DefaultQuery("order", "asc")
	//isDelete := ctx.DefaultQuery("is_deleted", "false")
	//filter := ctx.Query("filter")
	//joinOperator := ctx.DefaultQuery("join_operator", "and")

	//pages, err := pagination.New(page, perPage, 0, order, isDelete, filter, joinOperator)
	//if err != nil {
	//	log.Error().Err(err).Msg("pagination validation failed")
	//	helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
	//	return
	//}

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	adminService := services.NewAdminService(ctx, c.repo, c.cfg)

	response, err := adminService.GetListAdminUserService(pages)
	if err != nil {
		log.Error().Err(err).Msg("get list admin user service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})

}

func (c *Controller) GetDetailAdminUserController(ctx *gin.Context) {
	context := "AdminController.GetDetailAdminUser"
	id := ctx.Param("id")
	if id == "" {
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(errors.New("parameter id cannot blank or empty"), context).WithStatusCode(http.StatusBadRequest))
		return
	}

	adminService := services.NewAdminService(ctx, c.repo, c.cfg)

	response, err := adminService.GetDetailAdminUserService(id)
	if err != nil {
		log.Error().Err(err).Msg("get admin user service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) UpdateAdminUserController(ctx *gin.Context) {
	context := "AdminUser.UpdateAdminUser"
	request := new(web.DetailAdminResponse)

	id := ctx.Param("id")
	if id == "" {
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(errors.New("parameter id cannot blank or empty"), context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationUpdateAdminUserRequest(request)
	if err != nil {
		log.Debug().Err(err).Msg("should validate update admin user request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	adminService := services.NewAdminService(ctx, c.repo, c.cfg)

	response, err := adminService.UpdateAdminUserService(id, request)
	if err != nil {
		log.Error().Err(err).Msg("update admin user service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) DeleteAdminUserController(ctx *gin.Context) {
	context := "AdminController.DeleteAdminUser"
	id := ctx.Param("id")
	if id == "" {
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(errors.New("parameter id cannot blank or empty"), context).WithStatusCode(http.StatusBadRequest))
		return
	}

	adminService := services.NewAdminService(ctx, c.repo, c.cfg)

	response, err := adminService.DeleteAdminUserService(id)
	if err != nil {
		log.Error().Err(err).Msg("delete admin user service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}
