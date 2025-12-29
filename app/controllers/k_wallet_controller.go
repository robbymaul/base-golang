package controllers

import (
	"errors"
	"net/http"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/services"
	"paymentserviceklink/app/validate"
	"paymentserviceklink/app/web"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (c *Controller) AdminGetListKWalletController(ctx *gin.Context) {
	context := "KWalletController.AdminGetListKWallet"

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

	kWalletService := services.NewKWalletService(ctx, c.repo, c.cfg)

	response, err := kWalletService.AdminGetListKWalletService(pages)
	if err != nil {
		log.Error().Err(err).Msg("admin get list k-wallet service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetDetailKWalletController(ctx *gin.Context) {
	noRekening := ctx.Param("no-rekening")

	kWalletService := services.NewKWalletService(ctx, c.repo, c.cfg)

	response, err := kWalletService.AdminGetDetailKWalletService(noRekening)
	if err != nil {
		log.Error().Err(err).Interface("context", c.context).Msg("get detail k-wallet service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetListKWalletTransactionController(ctx *gin.Context) {
	context := "KWalletController.AdminGetListKWalletTransaction"
	noRekening := ctx.Param("no-rekening")

	if noRekening == "" {
		log.Error().Interface("no-rekening", noRekening).Msg("no-rekening is empty")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(errors.New("parameter no-rekening is empty"), context).WithStatusCode(http.StatusBadRequest))
		return
	}

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	kWalletService := services.NewKWalletService(ctx, c.repo, c.cfg)

	response, err := kWalletService.AdminGetListKWalletTransactionService(noRekening, pages)
	if err != nil {
		log.Error().Err(err).Msg("admin get list k-wallet transaction service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) AdminGetListTopupKWalletController(ctx *gin.Context) {
	context := "KWalletController.AdminGetListTopupKWallet"

	pages, err := c.Pagination(ctx)
	if err != nil {
		log.Error().Err(err).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.PaginationValidation(pages)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("pagination validation failed")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	kWalletService := services.NewKWalletService(ctx, c.repo, c.cfg)

	response, err := kWalletService.AdminGetListTopupKWalletService(pages)
	if err != nil {
		log.Error().Err(err).Msg("admin get list topup k-wallet service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) CreateKWalletController(ctx *gin.Context) {
	context := "KWalletController.CreateKWallet"
	request := new(web.CreateKWalletRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Error().Err(err).Msg("should bind json error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationCreateKWallet(request)
	if err != nil {
		log.Error().Err(err).Msg("validation create k-wallet error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	kWalletService := services.NewKWalletService(ctx, c.repo, c.cfg)

	response, err := kWalletService.CreateKWalletService(request)
	if err != nil {
		log.Error().Err(err).Msg("create k-wallet service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) GetKWalletMemberController(ctx *gin.Context) {
	context := "KWalletController.GetKWalletMemberController"
	request := new(web.GetKWalletRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Error().Err(err).Msg("should bind json error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationGetKWallet(request)
	if err != nil {
		log.Error().Err(err).Msg("validation get k-wallet error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	kWalletService := services.NewKWalletService(ctx, c.repo, c.cfg)

	response, err := kWalletService.GetKWalletMemberService(request)
	if err != nil {
		log.Error().Err(err).Msg("create k-wallet service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) GetListKWalletTransactionMemberController(ctx *gin.Context) {
	context := "KWalletController.GetListKWalletTransactionMember"
	request := new(web.GetListKWalletTransactionRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Error().Err(err).Msg("should bind json error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationGetListKWalletTransaction(request)
	if err != nil {
		log.Error().Err(err).Msg("validation get list k-wallet transaction error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	kWalletService := services.NewKWalletService(ctx, c.repo, c.cfg)

	response, err := kWalletService.GetListKWalletTransactionMemberService(request)
	if err != nil {
		log.Error().Err(err).Msg("get list k-wallet transaction member service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})

}

func (c *Controller) GetVirtualAccountKWalletController(ctx *gin.Context) {
	context := "KWalletController.GetVirtualAccountKWalletController"

	request := new(web.GetVirtualAccountKWalletRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Error().Err(err).Msg("should bind json error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationGetVirtualAccountKWallet(request)
	if err != nil {
		log.Error().Err(err).Msg("validation get virtual account k-wallet error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	kWalletService := services.NewKWalletService(ctx, c.repo, c.cfg)

	response, err := kWalletService.GetVirtualAccountKWalletService(request)
	if err != nil {
		log.Error().Err(err).Msg("get virtual account k-wallet service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) CreateTopupKWalletController(ctx *gin.Context) {
	context := "KWalletController.CreateTopupKWallet"
	request := new(web.CreateTopupKWalletRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Error().Err(err).Msg("should bind json error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
	}

	err = validate.ValidationCreateTopupKWallet(request)
	if err != nil {
		log.Error().Err(err).Msg("validation create topup k-wallet error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
	}

	kWalletService := services.NewKWalletService(ctx, c.repo, c.cfg)

	response, err := kWalletService.CreateTopupKWalletService(request)
	if err != nil {
		log.Error().Err(err).Msg("create topup k-wallet service error")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}
