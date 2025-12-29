package espay

import (
	"context"
	"encoding/json"
	"errors"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/web"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	CC_ENUM_URL_TYPE_PAY_RETURN = "PAY_RETURN"
	CC_ENUM_PAY_TYPE_REDIRECT   = "REDIRECT"
)

type CCResponse struct {
	ResponseCode       string `json:"responseCode"`
	ResponseMessage    string `json:"responseMessage"`
	PartnerReferenceNo string `json:"partnerReferenceNo"`
	ApprovalCode       string `json:"approvalCode"`
	WebRedirectUrl     string `json:"webRedirectUrl"`
}

type CreditCard struct {
	Espay *Espay
}

func NewCreditCard(espay *Espay) *CreditCard {
	return &CreditCard{
		Espay: espay,
	}
}

func (c *CreditCard) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
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

	chargeReq := c.createChargeHostToHostRequest(paymentReq, timeNow)
	log.Debug().Interface("chargeReq", chargeReq).Msg("chargeReq host to host")

	asymmetricSignatureGenerate, err := signature.AsymmetricSignatureGenerate(chargeReq, c.Espay.PrivateKey)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate asymmetric signature")
		return nil, err
	}
	log.Debug().Interface("asymmetricSignatureGenerate", asymmetricSignatureGenerate).Msg("asymmetric signature generate")

	result, err := c.Espay.SendPaymentHostToHost(ctx, chargeReq, asymmetricSignatureGenerate)
	if err != nil {
		log.Error().Err(err).Msg("send payment host to host espay")
		return nil, err
	}
	log.Debug().Interface("result", result).Msg("send payment host to host espay")

	return result, nil
}

func (c *CreditCard) ClientResponse(payment *models.Payments, channel *models.Channel) (*web.PaymentResponse, error) {
	log.Debug().Interface("payment", payment).Msg("client response virtual account")

	var ccResponse CCResponse

	err := json.Unmarshal(payment.GatewayResponse, &ccResponse)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal virtual account response")
		return nil, err
	}

	if ccResponse.ResponseCode != "2005400" {
		return nil, errors.New(ccResponse.ResponseMessage)
	}

	return &web.PaymentResponse{
		Id:            payment.Id,
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        payment.Status,
		Amount:        payment.Amount.IntPart(),
		FeeAdmin:      payment.FeeAmount.IntPart(),
		TotalAmount:   payment.TotalAmount.IntPart(),
		Currency:      payment.Currency,
		PaymentMethod: payment.Channel.PaymentMethod,
		PaymentType:   enums.PAYMENT_TYPE_REDIRECT,
		PaymentDetail: web.PaymentDetail{
			Bank:            payment.Channel.BankName,
			Url:             nil,
			VaNumber:        ccResponse.WebRedirectUrl,
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

func (c *CreditCard) createChargeHostToHostRequest(req PaymentRequest, timeNow time.Time) CreditCardRequest {
	return CreditCardRequest{
		PartnerReferenceNo: req.OrderID,
		MerchantId:         c.Espay.MerchantCode,
		SubMerchantId:      c.Espay.ApiKey,
		Amount: CCAmount{
			Value:    req.Amount,
			Currency: string(req.CCY),
		},
		UrlParam: CCUrlParam{
			Url:        req.ReturnUrl,
			Type:       CC_ENUM_URL_TYPE_PAY_RETURN,
			IsDeeplink: "N",
		},
		ValidUpTo:         timeNow,
		PointOfInitiation: "",
		PayOptionDetails: CCPayOptionDetails{
			PayMethod: req.BankCode,
			PayOption: req.ProductCode,
			TransAmount: CCTransAmount{
				Value:    req.Amount,
				Currency: string(req.CCY),
			},
			CCFeeAmount: CCFeeAmount{
				Value:    req.FeeAmount,
				Currency: string(req.CCY),
			},
		},
		AdditionalInfo: CCAdditionalInfo{
			PayType:       CC_ENUM_PAY_TYPE_REDIRECT,
			UserId:        req.CustomerID,
			UserName:      req.CustomerName,
			UserEmail:     req.CustomerEmail,
			UserPhone:     req.CustomerPhone,
			BuyerId:       "",
			ProductCode:   req.ProductCode,
			BalanceType:   "",
			BankCardToken: "",
		},
	}
}
