package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/services"
	"paymentserviceklink/app/validate"
	"paymentserviceklink/app/web"
	"strconv"
)

func (c *Controller) AdminCreatePlatformController(ctx *gin.Context) {
	context := "PlatformController.AdminCreatePlatform"

	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	request := new(web.CreatePlatformRequest)

	err = ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("payload", request).Interface("context", context).Msg("admin create platform controller")

	err = validate.ValidationCreatePlatformRequest(request)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should validate create platform request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminCreatePlatformService(session, request)
	if err != nil {
		log.Error().Err(err).Msg("create platform service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetListPlatformController(ctx *gin.Context) {
	context := "PlatformController.AdminGetListPlatform"

	//page := ctx.DefaultQuery("page", "1")
	//perPage := ctx.DefaultQuery("per_page", "10")
	//order := ctx.DefaultQuery("order", "desc")
	//isDelete := ctx.DefaultQuery("is_deleted", "false")
	//filter := ctx.Query("filter")
	//joinOperator := ctx.DefaultQuery("join_operator", "and")

	//log.Debug().Interface("page", page).Interface("per_page", perPage).
	//	Interface("order", order).Interface("context", context).Msg("admin get list platform controller")

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("pages", pages).Interface("context", context).Msg("pages pagination get list platform")

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, "").WithStatusCode(http.StatusBadRequest))
		return
	}

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminGetListPlatformService(pages)
	if err != nil {
		log.Error().Err(err).Msg("get list platform service failed")
		helpers.ErrorResponse(ctx, err)
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetDetailPlatformController(ctx *gin.Context) {
	context := "PlatformController.AdminGetDetailPlatform"

	platformId := ctx.Param("id-platform")
	platformIdInt, err := strconv.Atoi(platformId)
	if err != nil {
		log.Debug().Err(err).Msg("should convert id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("platform_id", platformId).Interface("context", context).Msg("admin get detail platform controller")

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminGetDetailPlatformService(int64(platformIdInt))
	if err != nil {
		log.Error().Err(err).Msg("get platform service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminUpdatePlatformController(ctx *gin.Context) {
	context := "PlatformController.AdminUpdatePlatform"

	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	platformId := ctx.Param("id-platform")
	platformIdInt, err := strconv.Atoi(platformId)
	if err != nil {
		log.Debug().Err(err).Msg("should convert id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("platform_id", platformId).Interface("context", context).Msg("admin update detail platform controller")

	request := new(web.DetailPlatformResponse)

	err = ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("payload", request).Interface("context", context).Msg("data payload update request")

	err = validate.ValidationUpdatePlatformRequest(request)
	if err != nil {
		log.Debug().Err(err).Msg("should validate update platform request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminUpdatePlatformService(session, int64(platformIdInt), request)
	if err != nil {
		log.Error().Err(err).Msg("update platform service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminUpdatePlatformSecretKeyController(ctx *gin.Context) {
	context := "PlatformController.AdminUpdatePlatformSecretKey"
	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	platformId := ctx.Param("id-platform")
	platformIdInt, err := strconv.Atoi(platformId)
	if err != nil {
		log.Debug().Err(err).Msg("should convert id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("platform_id", platformId).Interface("context", context).Msg("admin update platform secret key controller")

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminUpdatePlatformSecretKeyService(session, int64(platformIdInt))
	if err != nil {
		log.Error().Err(err).Msg("update platform secret service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminAssignmentPlatformConfigurationController(ctx *gin.Context) {
	context := "PlatformController.AdminAssignmentPlatformConfiguration"
	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	platformId := ctx.Param("id-platform")
	platformIdInt, err := strconv.Atoi(platformId)
	if err != nil {
		log.Debug().Err(err).Msg("should convert id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("platform_id", platformId).Interface("context", context).Msg("admin assignment platform configuration controller")

	var request []web.ResponseConfiguration

	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationArrayConfiguration(request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("validation array configuration error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminAssignmentPlatformConfigurationService(session, int64(platformIdInt), request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin assign platform configuration service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminRemovalPlatformConfigurationController(ctx *gin.Context) {
	context := "PlatformController.AdminRemovalPlatformConfiguration"

	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	platformId := ctx.Param("id-platform")
	platformIdInt, err := strconv.Atoi(platformId)
	if err != nil {
		log.Debug().Err(err).Msg("should convert id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("platform_id", platformId).Interface("context", context).Msg("admin removal platform configuration controller")

	var request []web.ResponseConfiguration

	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationArrayConfiguration(request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("validation array configuration error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminRemovalPlatformConfigurationService(session, int64(platformIdInt), request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin removal platform configuration service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetListConfigurationInPlatformController(ctx *gin.Context) {
	context := "PlatformController.AdminGetListConfigurationInPlatform"

	//page := ctx.DefaultQuery("page", "1")
	//perPage := ctx.DefaultQuery("per_page", "1")
	//sort := ctx.DefaultQuery("sort", "asc")
	//search := ctx.DefaultQuery("search", "")
	//isDelete := ctx.DefaultQuery("is_delete", "false")

	platformId := ctx.Param("id-platform")
	platformIdInt, err := strconv.Atoi(platformId)
	if err != nil {
		log.Debug().Err(err).Msg("should convert id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	//log.Debug().Interface("page", page).Interface("per_page", perPage).Interface("sort", sort).Interface("search", search).
	//	Interface("configuration_id", platformId).Interface("int configuration_id", platformIdInt).
	//	Interface("context", context).Msg("admin get list configuration in platform controller")

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("pages", pages).Interface("context", context).Msg("pagination get list configuration")

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, "").WithStatusCode(http.StatusBadRequest))
		return
	}

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminGetListConfigurationInPlatformService(int64(platformIdInt), pages)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin get list configuration in platform service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminAssignmentPlatformChannelController(ctx *gin.Context) {
	context := "PlatformController.AdminAssignmentPlatformChannel"
	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	platformId := ctx.Param("id-platform")
	platformIdInt, err := strconv.Atoi(platformId)
	if err != nil {
		log.Debug().Err(err).Msg("should convert id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("platform_id", platformId).Interface("context", context).Msg("admin assignment platform channel controller")

	var request []web.DetailChannelResponse

	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("should bind json request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationArrayChannel(request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("validation array channel error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminAssignmentPlatformChannelService(session, int64(platformIdInt), request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin assignment platform channel service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminAssignmentPlatformListChannelController(ctx *gin.Context) {
	context := "PlatformController.AdminAssignmentPlatformListChannel"
	platformId := ctx.Param("id-platform")
	platformIdInt, err := strconv.Atoi(platformId)
	if err != nil {
		log.Debug().Err(err).Msg("should convert id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("platform_id", platformId).Interface("context", context).Msg("admin assignment platform list channel controller")

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminAssignmentPlatformListChannelService(int64(platformIdInt))
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin assignment platform list channel service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminRemovalPlatformChannelController(ctx *gin.Context) {
	context := "PlatformController.AdminRemovalPlatformChannel"
	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	platformId := ctx.Param("id-platform")
	platformIdInt, err := strconv.Atoi(platformId)
	if err != nil {
		log.Debug().Err(err).Msg("should convert id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("platform_id", platformId).Interface("context", context).Msg("admin removal platform configuration controller")

	var request []web.DetailChannelResponse

	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("should bind json request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationArrayChannel(request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("validation array channel error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminRemovalPlatformChannelService(session, int64(platformIdInt), request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin removal platform channel service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetListChannelInPlatformController(ctx *gin.Context) {
	context := "PlatformController.AdminGetListChannelInPlatform"

	//page := ctx.DefaultQuery("page", "1")
	//perPage := ctx.DefaultQuery("per_page", "1")
	//order := ctx.DefaultQuery("order", "asc")
	//search := ctx.DefaultQuery("search", "")
	//isDelete := ctx.DefaultQuery("is_delete", "false")
	//filter := ctx.Query("filter")
	//joinOperator := ctx.DefaultQuery("join_operator", "and")

	platformId := ctx.Param("id-platform")
	platformIdInt, err := strconv.Atoi(platformId)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should convert id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	//log.Debug().Interface("page", page).Interface("per_page", perPage).Interface("order", order).
	//	Interface("configuration_id", platformId).Interface("int configuration_id", platformIdInt).
	//	Interface("context", context).Msg("admin get list configuration in platform controller")

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("pages", pages).Interface("context", context).Msg("pagination get list channels")

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, "").WithStatusCode(http.StatusBadRequest))
		return
	}

	platformService := services.NewPlatformService(ctx, c.repo, c.cfg)

	response, err := platformService.AdminGetListChannelInPlatformService(int64(platformIdInt), pages)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin get list channel in platform service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}
