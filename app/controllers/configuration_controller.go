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

func (c *Controller) AdminCreateConfigurationController(ctx *gin.Context) {
	context := "Configuration.CreateConfiguration"

	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	var request []web.CreateConfigurationRequest

	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		log.Debug().Err(err).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationCreateConfigurationRequest(request)
	if err != nil {
		log.Debug().Err(err).Msg("should validate create platform configuration request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	configurationService := services.NewConfigurationService(ctx, c.repo, c.cfg)

	response, err := configurationService.AdminCreateConfigurationService(session, request)
	if err != nil {
		log.Error().Err(err).Msg("admin create platform configuration failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetListConfigurationController(ctx *gin.Context) {
	context := "Configuration.AdminGetListConfiguration"

	//page := ctx.DefaultQuery("page", "1")
	//perPage := ctx.DefaultQuery("per_page", "10")
	//search := ctx.DefaultQuery("search", "")
	//order := ctx.DefaultQuery("order", "asc")
	aggregator := ctx.DefaultQuery("aggregator", "")
	//isDelete := ctx.DefaultQuery("is_deleted", "false")
	//filter := ctx.Query("filter")
	//joinOperator := ctx.DefaultQuery("join_operator", "and")

	//log.Debug().Interface("page", page).Interface("per_page", perPage).
	//	Interface("order", order).Interface("context", context).Msg("admin get list configuration controller")

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("pages", pages).Interface("context", context).Msg("pages get list configurations")

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, "").WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationAggregatorQuery(aggregator)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("aggregator query validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, "").WithStatusCode(http.StatusBadRequest))
		return
	}

	configurationService := services.NewConfigurationService(ctx, c.repo, c.cfg)

	response, err := configurationService.AdminGetListConfigurationService(pages, aggregator)
	if err != nil {
		log.Error().Err(err).Msg("get list platform configuration failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetConfigurationController(ctx *gin.Context) {
	context := "Configuration.AdminGetConfiguration"

	configurationId := ctx.Param("id-configuration")
	intConfigurationId, err := strconv.Atoi(configurationId)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should convert id configuration to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("configuration_id", configurationId).Interface("int configuration_id", intConfigurationId).
		Interface("context", context).Msg("admin get configuration controller")

	configurationService := services.NewConfigurationService(ctx, c.repo, c.cfg)

	response, err := configurationService.AdminGetConfigurationService(int64(intConfigurationId))
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin get configuration failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminUpdateConfigurationController(ctx *gin.Context) {
	context := "Configuration.AdminUpdateConfiguration"

	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	configurationId := ctx.Param("id-configuration")
	intConfigurationId, err := strconv.Atoi(configurationId)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should convert id configuration to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("configuration_id", configurationId).Interface("int configuration_id", intConfigurationId).
		Interface("context", context).Msg("admin update configuration controller")

	request := new(web.ResponseConfiguration)

	err = ctx.ShouldBindJSON(request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("should unmarshal json request body")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationUpdateConfigurationRequest(request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("validation update configuration request failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	configurationService := services.NewConfigurationService(ctx, c.repo, c.cfg)

	response, err := configurationService.AdminUpdateConfigurationService(session, int64(intConfigurationId), request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin update configuration service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminAssignmentConfigurationToPlatformController(ctx *gin.Context) {
	context := "Configuration.AdminAssignmentConfigurationToPlatform"

	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	configurationId := ctx.Param("id-configuration")
	intConfigurationId, err := strconv.Atoi(configurationId)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should convert id configuration to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("configuration_id", configurationId).Interface("int configuration_id", intConfigurationId).
		Interface("context", context).Msg("admin assignment configuration to configuration controller")

	var request []web.DetailPlatformResponse

	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("should unmarshal json request body")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("payload", request).Interface("context", context).Msg("data payload request body")

	err = validate.ValidationConfigurationPlatformRequest(request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("validation assignment configuration request failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	configurationService := services.NewConfigurationService(ctx, c.repo, c.cfg)

	response, err := configurationService.AdminAssignmentConfigurationToPlatformService(session, int64(intConfigurationId), request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin assignment configuration to platform service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminRemovalPlatformFromConfigurationController(ctx *gin.Context) {
	context := "Configuration.AdminRemovalPlatformFromConfiguration"

	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	configurationId := ctx.Param("id-configuration")
	intConfigurationId, err := strconv.Atoi(configurationId)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should convert id configuration to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("configuration_id", configurationId).Interface("int configuration_id", intConfigurationId).
		Interface("context", context).Msg("admin removal platform from configuration controller")

	var request []web.DetailPlatformResponse

	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("should unmarshal json request body")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("payload", request).Interface("context", context).Msg("data payload request body")

	err = validate.ValidationConfigurationPlatformRequest(request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("validation removal platform request failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	configurationService := services.NewConfigurationService(ctx, c.repo, c.cfg)

	response, err := configurationService.AdminRemovalPlatformFromConfigurationService(session, int64(intConfigurationId), request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin removal platform from configuration service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetListPlatformInConfigurationController(ctx *gin.Context) {
	context := "Configuration.AdminGetListPlatformInConfiguration"

	//page := ctx.DefaultQuery("page", "1")
	//perPage := ctx.DefaultQuery("per_page", "1")
	//order := ctx.DefaultQuery("order", "asc")
	//search := ctx.DefaultQuery("search", "")
	//isDelete := ctx.DefaultQuery("is_delete", "false")

	configurationId := ctx.Param("id-configuration")
	intConfigurationId, err := strconv.Atoi(configurationId)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should convert id configuration to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	//log.Debug().Interface("page", page).Interface("per_page", perPage).Interface("order", order).Interface("search", search).
	//	Interface("configuration_id", configurationId).Interface("int configuration_id", intConfigurationId).
	//	Interface("context", context).Msg("admin get list platform in configuration controller")

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("pages", pages).Interface("context", context).Msg("pagination get list platform")

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, "").WithStatusCode(http.StatusBadRequest))
		return
	}

	configurationService := services.NewConfigurationService(ctx, c.repo, c.cfg)

	response, err := configurationService.AdminGetListPlatformInConfigurationService(int64(intConfigurationId), pages)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin get list platform in configuration service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}
