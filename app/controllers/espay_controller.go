package controllers

import (
	"bytes"
	"io"
	"net/http"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/services"
	"paymentserviceklink/app/validate"
	"paymentserviceklink/app/web"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (c *Controller) EspayValidationInquiryController(ctx *gin.Context) {
	//rqUUID := ctx.DefaultQuery("rq_uuid", "")
	//rqDateTime := ctx.DefaultQuery("rq_datetime", "")
	//senderId := ctx.DefaultQuery("sender_id", "")
	//receiverId := ctx.DefaultQuery("receiver_id", "")
	//password := ctx.DefaultQuery("password", "")
	//commCode := ctx.DefaultQuery("comm_code", "")
	//memberCode := ctx.DefaultQuery("member_code", "")
	//orderId := ctx.DefaultQuery("order_id", "")
	//signature := ctx.DefaultQuery("signature", "")

	request := new(web.EspayInquiryRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Error().Err(err).Msg("error parsing request body espay validation inquiry")
		espayCode, message := enums.CreateEspayCodeResponse(http.StatusBadRequest, enums.PAYMENT, enums.ESPAY_INVALID_MISSING_FIELD_FORMAT, err)
		ctx.JSON(http.StatusOK, &web.EspayInquiryResponse{
			ResponseCode:       espayCode,
			ResponseMessage:    message,
			VirtualAccountData: nil,
		})
		return
	}
	log.Debug().Interface("request", request).Msg("parse request body espay validation inquiry controller")

	err = validate.ValidationEspayInquiryRequest(request)
	if err != nil {
		log.Error().Err(err).Msg("error validating request body espay validation inquiry")
		espayCode, message := enums.CreateEspayCodeResponse(http.StatusBadRequest, enums.PAYMENT, enums.ESPAY_INVALID_MANDATORY_FIELD, err)
		ctx.JSON(http.StatusOK, &web.EspayInquiryResponse{
			ResponseCode:       espayCode,
			ResponseMessage:    message,
			VirtualAccountData: nil,
		})
		return
	}

	//fmt.Println(rqUUID, senderId, receiverId, password, commCode, memberCode, signature)

	paymentService := services.NewPaymentService(ctx, c.repo, c.cfg)

	response := paymentService.EspayValidationInquiryService(request)
	log.Debug().Interface("response", response).Msg("espay validation inquiry service response")

	//dtAmount := "15000.00"

	// pesan
	//msg := "Transaction will be Process!"

	// format response sesuai Espay
	//response := fmt.Sprintf("0;Success;%s;%s;IDR;%s;%s", orderId, dtAmount, msg, rqDateTime)

	//ctx.Data(http.StatusOK, "text/xml; charset=utf-8", []byte(response))
	ctx.JSON(http.StatusOK, response)
}

/*
param :
*/
func (c *Controller) EspayPaymentNotificationController(ctx *gin.Context) {
	request := new(web.EspayPaymentNotificationRequest)

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		log.Error().Err(err).Msg("error parsing request body espay payment notification")
		espayCode, message := enums.CreateEspayCodeResponse(http.StatusBadRequest, enums.PAYMENT, enums.ESPAY_INVALID_MISSING_FIELD_FORMAT, err)
		ctx.JSON(http.StatusOK, &web.EspayPaymentNotificationResponse{
			ResponseCode:    espayCode,
			ResponseMessage: message,
		})
		return
	}
	log.Debug().Interface("request", request).Msg("parse request body espay payment notification controller")

	err = validate.ValidationEspayPaymentNotificationRequest(request)
	if err != nil {
		log.Error().Err(err).Msg("error validating request body espay payment notification")
		espayCode, message := enums.CreateEspayCodeResponse(http.StatusBadRequest, enums.PAYMENT, enums.ESPAY_INVALID_MISSING_FIELD_FORMAT, err)
		ctx.JSON(http.StatusOK, &web.EspayPaymentNotificationResponse{
			ResponseCode:    espayCode,
			ResponseMessage: message,
		})
		return
	}

	paymentService := services.NewPaymentService(ctx, c.repo, c.cfg)

	response := paymentService.EspayPaymentNotificationService(request)
	log.Debug().Interface("response", response).Msg("espay payment notification service response")

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) EspayTopupNotificationController(ctx *gin.Context) {
	// Baca seluruh body sekaligus
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
	}
	defer ctx.Request.Body.Close()

	// Logging body mentah (dalam bentuk string/JSON)
	log.Debug().Str("raw_body", string(bodyBytes)).Msg("Espay notification raw request body")

	// Buat reader baru agar body bisa dibaca ulang oleh ShouldBindJSON
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	request := new(web.EspayTopupNotificationRequest)

	err = ctx.ShouldBind(request)
	if err != nil {
		log.Error().Err(err).Msg("error parsing request body espay topup notification")
		espayCode, message := enums.CreateEspayCodeResponse(http.StatusBadRequest, enums.PAYMENT, enums.ESPAY_INVALID_MISSING_FIELD_FORMAT, err)
		ctx.JSON(http.StatusOK, &web.EspayPaymentNotificationResponse{
			ResponseCode:    espayCode,
			ResponseMessage: message,
		})
		return
	}
	log.Debug().Interface("request", request).Msg("parse request body espay topup notification controller")

	err = validate.ValidationEspayTopupNotificationRequest(request)
	if err != nil {
		log.Error().Err(err).Msg("error validating request body espay topup notification")
		espayCode, message := enums.CreateEspayCodeResponse(http.StatusBadRequest, enums.PAYMENT, enums.ESPAY_INVALID_MISSING_FIELD_FORMAT, err)
		ctx.JSON(http.StatusOK, &web.EspayPaymentNotificationResponse{
			ResponseCode:    espayCode,
			ResponseMessage: message,
		})
		return
	}

	paymentService := services.NewPaymentService(ctx, c.repo, c.cfg)

	response := paymentService.EspayTopupNotificationService(request)
	log.Debug().Interface("response", response).Msg("espay topup notification service response")

	ctx.JSON(http.StatusOK, response)
}
