package espay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/web"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

type VirtualAccount struct {
	Espay *Espay
}

type VirtualAccountResponse struct {
	RqUUID       string `json:"rq_uuid"`
	RsDateTime   string `json:"rs_datetime"`
	ErrorCode    any    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	VaNumber     string `json:"va_number"`
	Expired      string `json:"expired"`
	Description  string `json:"description"`
	TotalAmount  string `json:"total_amount"`
	Amount       string `json:"amount"`
	Fee          string `json:"fee"`
}

func NewVirtualAccount(espay *Espay) *VirtualAccount {
	return &VirtualAccount{Espay: espay}
}

func (v *VirtualAccount) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	log.Debug().Interface("request", req).Msg("pay virtual account")
	paymentReq, ok := req.(PaymentRequest)
	if !ok {
		return nil, errors.New("invalid payment request type")
	}
	log.Debug().Interface("paymentReq", paymentReq).Msg("payment request")

	// signature
	signature := NewSignature(v.Espay.SignatureKey, v.Espay.MerchantCode)

	signatureSendInvoice := signature.HashBashSignatureGenerate(paymentReq, enums.ENUM_SIGNATURE_ESPAY_SENDINVOICE)
	log.Debug().Interface("signature", signatureSendInvoice).Msg("signature")

	// Convert to Midtrans charge request
	chargeReq := v.createChargeSendInvoiceRequest(paymentReq, signatureSendInvoice)
	log.Debug().Interface("chargeReq", chargeReq).Msg("chargeReq")

	result, err := v.Espay.SendInvoicePaymentVa(ctx, chargeReq)
	if err != nil {
		log.Error().Err(err).Msg("failed to send charge request send invoice payment va")
		return nil, err
	}
	log.Debug().Interface("result", result).Msg("send invoice payment va success")

	return result, nil
}

func (v *VirtualAccount) ClientResponse(payment *models.Payments, channel *models.Channel) (*web.PaymentResponse, error) {
	log.Debug().Interface("payment", payment).Msg("client response virtual account")

	var vaResponse VirtualAccountResponse

	err := json.Unmarshal(payment.GatewayResponse, &vaResponse)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal virtual account response")
		return nil, err
	}

	if vaResponse.ErrorCode != "0000" {
		return nil, errors.New(vaResponse.ErrorMessage)
	}

	feeAdmin, err := decimal.NewFromString(vaResponse.Fee)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert fee admin")
		return nil, err
	}

	amount, err := decimal.NewFromString(vaResponse.Amount)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert amount")
		return nil, err
	}

	totalAmount := amount.Add(feeAdmin)

	return &web.PaymentResponse{
		Id:            payment.Id,
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        payment.Status,
		Amount:        amount.IntPart(),
		FeeAdmin:      feeAdmin.IntPart(),
		TotalAmount:   totalAmount.IntPart(),
		Currency:      payment.Currency,
		PaymentMethod: payment.Channel.PaymentMethod,
		PaymentType:   enums.PAYMENT_TYPE_VA,
		PaymentDetail: web.PaymentDetail{
			Bank:            payment.Channel.BankName,
			Url:             nil,
			VaNumber:        vaResponse.VaNumber,
			BillKey:         "",
			BIllCode:        "",
			TransactionTime: vaResponse.RsDateTime,
			ExpireTime:      vaResponse.Expired,
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

func (v *VirtualAccount) createChargeSendInvoiceRequest(req PaymentRequest, signature string) map[string]string {
	return map[string]string{
		"rq_uuid":     req.RQUUID,                           // uuid -> clear
		"rq_datetime": req.RQDateTime.Format(time.DateTime), // time transaction -> clear
		"order_id":    req.OrderID,                          // order id -> clear
		"amount":      fmt.Sprint(req.Amount),
		"ccy":         string(req.CCY),
		"comm_code":   v.Espay.MerchantCode,
		"remark1":     req.CustomerPhone,
		"remark2":     req.CustomerName,
		"remark3":     req.CustomerEmail,
		"remark4":     v.Espay.MerchantName,
		"update":      "N",
		"bank_code":   req.BankCode,
		"va_expired":  fmt.Sprint(req.VaExpired),
		"signature":   signature,
	}
}
