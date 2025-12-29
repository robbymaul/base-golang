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

func (c *Controller) AdminCreateAggregatorController(ctx *gin.Context) {
	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Interface("context", c.context).Msg("get admin session")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("session admin")

	c.context = "Aggregator.AdminCreateAggregator"

	request := new(web.CreateAggregatorRequest)

	err = ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Interface("context", c.context).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, c.context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("payload", request).Interface("context", c.context).Msg("admin create aggregator controller payload")

	err = validate.ValidationCreateAggregatorRequest(request)
	if err != nil {
		log.Debug().Err(err).Interface("context", c.context).Msg("should validate create aggregator request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, c.context).WithStatusCode(http.StatusBadRequest))
		return
	}

	aggregatorService := services.NewAggregatorService(ctx, c.repo, c.cfg)

	response, err := aggregatorService.AdminCreateAggregatorService(session, request)
	if err != nil {
		log.Error().Err(err).Interface("context", c.context).Msg("create aggregator service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetAllAggregatorController(ctx *gin.Context) {
	context := "Aggregator.AdminGetAllAggregator"
	//page := ctx.DefaultQuery("page", "1")
	//perPage := ctx.DefaultQuery("per_page", "10")
	//order := ctx.DefaultQuery("order", "asc")
	////search := ctx.DefaultQuery("search", "")
	//isDelete := ctx.DefaultQuery("is_deleted", "false")
	//filter := ctx.Query("filter")
	//joinOperator := ctx.DefaultQuery("join_operator", "and")

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	aggregatorService := services.NewAggregatorService(ctx, c.repo, c.cfg)

	response, err := aggregatorService.AdminGetAllAggregatorService(pages)
	if err != nil {
		log.Error().Err(err).Interface("context", c.context).Msg("get all aggregator service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetAggregatorController(ctx *gin.Context) {
	c.context = "Aggregator.AdminGetDetailAggregator"
	id := ctx.Param("id-aggregator")
	parseInt, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("convert id to int failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, c.context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("id param", id).Interface("parse int", parseInt).Interface("context", c.context).
		Msg("admin get aggregator controller")

	aggregatorService := services.NewAggregatorService(ctx, c.repo, c.cfg)

	response, err := aggregatorService.AdminGetAggregatorService(int64(parseInt))
	if err != nil {
		log.Error().Err(err).Interface("context", c.context).Msg("get detail aggregator service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminUpdateAggregatorController(ctx *gin.Context) {
	c.context = "Aggregator.AdminUpdateAggregator"

	session, err := c.GetAdminSession(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("get admin session failed")
		helpers.ErrorResponse(ctx, err)
		return
	}
	log.Debug().Interface("session data", session).Msg("admin session")

	id := ctx.Param("id-aggregator")
	parseInt, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("convert id to int failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, c.context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("param id", id).Interface("parse id", parseInt).Interface("context", c.context).
		Msg("admin update aggregator controller")

	request := new(web.AggregatorResponse)

	err = ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, c.context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationUpdateAggregatorRequest(request)
	if err != nil {
		log.Debug().Err(err).Msg("should validate update aggregator request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, c.context).WithStatusCode(http.StatusBadRequest))
		return
	}

	aggregatorService := services.NewAggregatorService(ctx, c.repo, c.cfg)

	response, err := aggregatorService.AdminUpdateAggregatorService(session, int64(parseInt), request)
	if err != nil {
		log.Error().Err(err).Msg("update aggregator service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}
