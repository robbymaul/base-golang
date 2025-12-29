package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog/log"
	"net/http"
	"paymentserviceklink/app/services"
	"paymentserviceklink/app/web"
)

func (c *Controller) WebhookCallbackSenangpayPaymentNotification(ctx *gin.Context) {
	//context := "Webhook.WebhookCallbackSenangpayPaymentNotification"
	request := new(web.WebhookCallbackSenangpay)

	err := ctx.ShouldBindWith(request, binding.Form)
	if err != nil {
		log.Debug().Err(err).Msg("failed to bind senangpay form data")
		ctx.String(http.StatusOK, "OK")
		return
	}
	log.Debug().Interface("payload", request).Msg("payload webhook callback senangpay notification")

	paymentService := services.NewPaymentService(ctx, c.repo, c.cfg)

	_ = paymentService.CallbackSenangpayNotificationService(request)

	ctx.String(http.StatusOK, "OK")
}
