package controllers

import (
	"net/http"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/services"
	"paymentserviceklink/app/validate"
	"paymentserviceklink/app/web"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (c *Controller) CreatePaymentController(ctx *gin.Context) {
	context := "Payment.CreatePaymentController"
	request := new(web.CreatePaymentRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("payload", request).Interface("context", context).Msg("data payload create payment request")

	err = validate.ValidationCreatePaymentRequest(request)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should validate create payment request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	paymentService := services.NewPaymentService(ctx, c.repo, c.cfg)

	response, err := paymentService.CreatePaymentService(request)
	if err != nil {
		log.Error().Err(err).Msg("create payment service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})

}

func (c *Controller) GetDetailPaymentController(ctx *gin.Context) {
	context := "PaymentController.GetDetailPayment"
	request := new(web.GetDetailPaymentRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	err = validate.ValidationGetDetailPaymentRequest(request)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("validation get detail payment request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	paymentService := services.NewPaymentService(ctx, c.repo, c.cfg)

	response, err := paymentService.GetDetailPaymentService(request)
	if err != nil {
		log.Error().Err(err).Msg("get detail payment service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) CheckStatusPaymentController(ctx *gin.Context) {
	context := "PaymentController.CheckStatusPayment"
	request := new(web.CheckStatusPaymentRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("payload", request).Interface("context", context).Msg("data payload check status payment request")

	err = validate.CheckStatusPaymentRequest(request)
	if err != nil {
		log.Debug().Err(err).Interface("context", context).Msg("should validate create payment request")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}

	paymentService := services.NewPaymentService(ctx, c.repo, c.cfg)

	response, err := paymentService.CheckStatusPaymentService(request)
	if err != nil {
		log.Error().Err(err).Msg("check status payment service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}

func (c *Controller) CheckKeyMidtransController(ctx *gin.Context) {
	//context := "Payment.CheckKeyMidtransController"

	paymentService := services.NewPaymentService(ctx, c.repo, c.cfg)

	response, err := paymentService.CheckKeyMidtransService()
	if err != nil {
		log.Error().Err(err).Msg("check key midtrans service failed")
		helpers.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}
