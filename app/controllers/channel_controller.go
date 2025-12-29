package controllers

import (
	"net/http"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/services"
	"paymentserviceklink/app/validate"
	"paymentserviceklink/app/web"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (c *Controller) AdminCreateChannelController(ctx *gin.Context) {
	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("get admin session failed")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("admin session")

	context := "Channel.CreateChannel"

	var request []web.CreateChannelRequest

	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		log.Debug().Err(err).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("payload", request).Msg("payload create channel request")

	err = validate.ValidationCreateChannelRequest(request)
	if err != nil {
		log.Debug().Err(err).Msg("should validate create payment method request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	channelService := services.NewChannelService(ctx, c.repo, c.cfg)

	response, err := channelService.AdminCreateChannelService(session, request)
	if err != nil {
		log.Error().Err(err).Msg("create payment method failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetListChannelController(ctx *gin.Context) {
	context := "Channel.GetListChannel"

	//page := ctx.DefaultQuery("page", "1")
	//perPage := ctx.DefaultQuery("per_page", "10")
	//search := ctx.DefaultQuery("search", "")
	//order := ctx.DefaultQuery("order", "desc")
	//isDelete := ctx.DefaultQuery("is_deleted", "false")
	//filter := ctx.Query("filter")
	//joinOperator := ctx.DefaultQuery("join_operator", "and")
	//currency := enums.Currency(ctx.DefaultQuery("currency", ""))

	//log.Debug().Interface("page", page).Interface("per_page", perPage).
	//	Interface("order", order).
	//	Interface("context", context).
	//	Msg("parameter admin get list channel controller")

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	paymentMethodService := services.NewChannelService(ctx, c.repo, c.cfg)

	response, err := paymentMethodService.AdminGetListChannelService(pages)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin get list payment method failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) GetListPlatformChannelController(ctx *gin.Context) {
	context := "Channel.GetListChannel"

	//page := ctx.DefaultQuery("page", "1")
	//perPage := ctx.DefaultQuery("per_page", "10")
	//search := ctx.DefaultQuery("search", "")
	//order := ctx.DefaultQuery("order", "desc")
	currency := enums.Currency(ctx.DefaultQuery("currency", ""))
	idPlatform := ctx.Param("id-platform")
	//isDelete := ctx.DefaultQuery("is_delete", "false")

	intIdPlatform, err := strconv.Atoi(idPlatform)
	if err != nil {
		log.Debug().Err(err).Msg("should convert id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	request := &web.GetListChannelRequest{
		Currency: currency,
	}

	//pages, err := pagination.New(page, perPage, 0, order, isDelete, "", "")
	//if err != nil {
	//	log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
	//	helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
	//	return
	//}

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.GetListChannelRequest(request)
	if err != nil {
		log.Error().Err(err).Msg("request validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	paymentMethodService := services.NewChannelService(ctx, c.repo, c.cfg)

	response, err := paymentMethodService.GetListChannelService(int64(intIdPlatform), pages, request)
	if err != nil {
		log.Error().Err(err).Msg("get list payment method failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetDetailChannelController(ctx *gin.Context) {
	context := "Channel.GetDetailChannel"

	channelId := ctx.Param("id-channel")

	intChannelId, err := strconv.Atoi(channelId)
	if err != nil {
		log.Debug().Err(err).Msg("should convert channel id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("channel_id", channelId).
		Interface("int channel id", intChannelId).
		Interface("context", context).
		Msg("admin get detail channel controller")

	paymentMethodService := services.NewChannelService(ctx, c.repo, c.cfg)

	response, err := paymentMethodService.AdminGetDetailChannelService(int64(intChannelId))
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin get detail channel service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminUpdateChannelController(ctx *gin.Context) {
	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("get admin session failed")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("admin session")

	context := "Channel.UpdateChannel"

	request := new(web.DetailChannelResponse)
	channelId := ctx.Param("id-channel")

	intChannelId, err := strconv.Atoi(channelId)
	if err != nil {
		log.Warn().Err(err).Interface("context", context).Msg("should convert channel id to int")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("channel_id", channelId).
		Interface("context", context).Msg("admin update channel controller")

	err = ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationUpdateChannelRequest(request)
	if err != nil {
		log.Debug().Err(err).Msg("should validate update payment method request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("payload", request).Interface("context", context).Msg("data payload update channel")

	paymentMethodService := services.NewChannelService(ctx, c.repo, c.cfg)

	response, err := paymentMethodService.AdminUpdateChannelService(session, request, int64(intChannelId))
	if err != nil {
		log.Error().Err(err).Msg("admin update payment method failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) ClientGetListChannelController(ctx *gin.Context) {
	context := "Channel.GetListChannel"

	request := new(web.GetListChannelRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.GetListChannelRequest(request)
	if err != nil {
		log.Error().Err(err).Msg("request validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	paymentMethodService := services.NewChannelService(ctx, c.repo, c.cfg)

	response, err := paymentMethodService.ClientGetListChannelService(request)
	if err != nil {
		log.Error().Err(err).Msg("get list payment method failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}
