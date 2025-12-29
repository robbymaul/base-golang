package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"paymentserviceklink/app/client/midtrans"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/services"
)

func (c *Controller) MidtransPaymentNotification(ctx *gin.Context) {
	context := "MidtransController.MidtransPaymentNotification"
	request := new(midtrans.CheckStatusPaymentResponse)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Error().Err(err).Interface("context", "MidtransController.MidtransPaymentNotification").Msg("should bind json body request error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusBadRequest))
		return
	}
	log.Debug().Interface("payload", request).Interface("context", context).Msg("midtrans payment notification controller payload")

	paymentService := services.NewPaymentService(ctx, c.repo, c.cfg)

	_, err = paymentService.MidtransPaymentNotificationService(request)
	if err != nil {
		log.Error().Err(err).Interface("context", context).Msg("midtrans payment notification service error")
		helpers.ErrorResponse(ctx, helpers.NewErrorTrace(err, context).WithStatusCode(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
