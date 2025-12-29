package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"paymentserviceklink/app/services"
	"paymentserviceklink/app/web"
)

func (c *Controller) HealthController(ctx *gin.Context) {
	healthService := services.NewHealthService(ctx, c.startTime, c.repo, c.cfg)

	response := healthService.GetHealthService()

	ctx.JSON(http.StatusOK, web.ResponseWeb{
		Message: "OK",
		Success: true,
		Data:    response,
	})
}
