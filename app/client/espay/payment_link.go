package espay

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/web"
	"time"

	"github.com/rs/zerolog/log"
)

type PaymentLinkResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	QrUrl           string `json:"qrUrl"`
}

type PaymentLink struct {
	Espay *Espay
}

func NewPaymentLink(e *Espay) *PaymentLink {
	return &PaymentLink{
		Espay: e,
	}
}

func (c *PaymentLink) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	timeNow := time.Now()
	log.Debug().Interface("req", req).Msg("pay credit card")
	paymentReq, ok := req.(PaymentRequest)
	if !ok {
		return nil, errors.New("invalid payment request type")
	}

	//signature
	signature := NewSignature(c.Espay.SignatureKey, c.Espay.MerchantCode)

	//hashSignature := signature.HashBashSignatureGenerate(paymentReq, enums.ENUM_SIGNATURE_ESPAY_SENDINVOICE)
	//log.Debug().Interface("signature", hashSignature).Msg("signature")

	chargeReq := c.createChargePaymentLinkRequest(paymentReq, timeNow, paymentReq)
	log.Debug().Interface("chargeReq", chargeReq).Msg("chargeReq host to host")

	asymmetricSignatureGenerate, err := signature.AsymmetricSignatureGenerate(chargeReq, c.Espay.PrivateKey)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate asymmetric signature")
		return nil, err
	}
	log.Debug().Interface("asymmetricSignatureGenerate", asymmetricSignatureGenerate).Msg("asymmetric signature generate")

	result, err := c.Espay.SendPaymentPaymentLink(ctx, chargeReq, asymmetricSignatureGenerate)
	if err != nil {
		log.Error().Err(err).Msg("send payment host to host espay")
		return nil, err
	}
	log.Debug().Interface("result", result).Msg("send payment host to host espay")

	return result, nil
}

func (c *PaymentLink) ClientResponse(payment *models.Payments, channel *models.Channel) (*web.PaymentResponse, error) {
	log.Debug().Interface("payment", payment).Msg("client response virtual account")

	var paymentLinkResponse PaymentLinkResponse

	err := json.Unmarshal(payment.GatewayResponse, &paymentLinkResponse)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal virtual account response")
		return nil, err
	}

	if paymentLinkResponse.ResponseCode != "2005400" {
		return nil, errors.New(paymentLinkResponse.ResponseMessage)
	}

	return &web.PaymentResponse{
		Id:            payment.Id,
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        payment.Status,
		Amount:        payment.TotalAmount.IntPart(),
		Currency:      payment.Currency,
		PaymentMethod: payment.Channel.PaymentMethod,
		PaymentType:   enums.PAYMENT_TYPE_REDIRECT,
		PaymentDetail: web.PaymentDetail{
			Bank: payment.Channel.BankName,
			Url: []*web.Actions{
				{
					Name:   "qr url",
					Method: http.MethodGet,
					Url:    paymentLinkResponse.QrUrl,
				},
			},
			VaNumber:        "",
			BillKey:         "",
			BIllCode:        "",
			TransactionTime: payment.CreatedAt.Format(time.DateTime),
			ExpireTime:      payment.ExpiredTime,
			Instruction:     channel.Instruction,
		},
		Customer: web.Customer{
			MemberId: payment.CustomerId,
			Name:     payment.CustomerName,
			Email:    payment.CustomerEmail,
			Phone:    payment.CustomerPhone,
		},
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}, nil
}

func (c *PaymentLink) createChargePaymentLinkRequest(req PaymentRequest, timeNow time.Time, paymentReq PaymentRequest) PaymentLinkRequest {
	return PaymentLinkRequest{
		//PartnerReferenceNo: req.OrderID,
		//MerchantId:         c.Espay.MerchantCode,
		//Amount: CCAmount{
		//	Value:    req.Amount,
		//	Currency: string(req.CCY),
		//},
		//AdditionalInfo: struct {
		//	ProductCode string `json:"productCode"`
		//}{
		//	ProductCode: paymentReq.ProductCode,
		//},
		//ValidityPeriod: timeNow,
	}
}
