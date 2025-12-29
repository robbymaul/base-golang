package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/services"
	"paymentserviceklink/app/validate"
	"paymentserviceklink/app/web"
	"paymentserviceklink/pkg/pagination"
	"time"
)

func (c *Controller) AdminGetListTransactionController(ctx *gin.Context) {
	context := "TransactionController.AdminGetListTransaction"
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "1")
	order := ctx.DefaultQuery("order", "desc")
	search := ctx.DefaultQuery("search", "")
	platformId := ctx.DefaultQuery("platform_id", "")
	currency := ctx.DefaultQuery("currency", "")
	startDate := ctx.DefaultQuery("start_date", "")
	endDate := ctx.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	isDelete := ctx.DefaultQuery("is_deleted", "false")
	filter := ctx.Query("filter")
	joinOperator := ctx.DefaultQuery("join_operator", "and")

	log.Debug().Interface("page", page).Interface("per_page", perPage).Interface("order", order).Interface("search", search).
		Interface("platform_id", platformId).Interface("currency", currency).Interface("start_date", startDate).
		Interface("end_date", endDate).Interface("context", context).Msg("pagination get list transactions")

	pages, err := pagination.New(page, perPage, 0, order, isDelete, filter, joinOperator)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, "").WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("pages", pages).Interface("context", context).Msg("pagination get list transactions")

	listRequest := &web.ListTransactionRequest{
		Currency:  currency,
		StartDate: startDate,
		EndDate:   endDate,
	}
	log.Debug().Interface("listRequest", listRequest).Interface("context", context).Msg("list transaction request")

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ListTransactionRequestValidation(listRequest)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	transactionService := services.NewTransactionService(ctx, c.repo, c.cfg)

	response, err := transactionService.AdminGetListTransactionService(pages, listRequest)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin get list transaction service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetDetailTransactionController(ctx *gin.Context) {
	context := "TransactionController.AdminGetDetailTransaction"

	transactionId := ctx.Param("id-transaction")

	transactionService := services.NewTransactionService(ctx, c.repo, c.cfg)

	response, err := transactionService.AdminGetDetailTransaction(transactionId)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("admin get detail transaction service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}
