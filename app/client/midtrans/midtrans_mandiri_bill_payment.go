package midtrans

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/web"

	"github.com/rs/zerolog/log"
)

type MandiriBillPayment struct {
	midtrans *Midtrans
}

type MANDIRIVAResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	OrderId           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	BillKey           string `json:"bill_key"`
	BillerCode        string `json:"biller_code"`
	Currency          string `json:"currency"`
	ExpiryTime        string `json:"expiry_time"`
}

func NewMandiriBillPayment(midtrans *Midtrans) *MandiriBillPayment {
	return &MandiriBillPayment{
		midtrans: midtrans,
	}
}

func (m *MandiriBillPayment) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	paymentReq, ok := req.(PaymentRequest)
	if !ok {
		return nil, errors.New("invalid payment request type")
	}

	// Convert to Midtrans charge request
	chargeReq := m.createChargeRequest(paymentReq)

	result, err := m.midtrans.CreatePayment(ctx, chargeReq)
	if err != nil {
		log.Error().Interface("result", result).Err(err).Msg("Failed to create payment midtrans mandiri bill payment")
		return nil, err
	}

	return result, nil
}

func (m *MandiriBillPayment) ClientResponse(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error) {
	var mandiriVaResponse MANDIRIVAResponse

	err := json.Unmarshal(payment.GatewayResponse, &mandiriVaResponse)
	if err != nil {
		return nil, err
	}

	if mandiriVaResponse.StatusCode != "201" {
		return nil, helpers.NewErrorTrace(fmt.Errorf(mandiriVaResponse.StatusMessage), m.midtrans.serviceName).WithStatusCode(http.StatusPaymentRequired)
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
		PaymentMethod: channel.PaymentMethod,
		PaymentType:   enums.PAYMENT_BILL,
		PaymentDetail: web.PaymentDetail{
			Bank:            enums.CHANNEL_MANDIRI,
			Url:             nil,
			VaNumber:        "",
			BillKey:         mandiriVaResponse.BillKey,
			BIllCode:        mandiriVaResponse.BillerCode,
			TransactionTime: mandiriVaResponse.TransactionTime,
			ExpireTime:      mandiriVaResponse.ExpiryTime,
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

func (m *MandiriBillPayment) createChargeRequest(req PaymentRequest) map[string]interface{} {
	return map[string]interface{}{
		"payment_type": "echannel",
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": req.Amount,
		},
		//"item_details": []map[string]interface{}{
		//	{
		//		"id":       "a1",
		//		"price":    50000,
		//		"quantity": 2,
		//		"name":     "Apel",
		//	},
		//	{
		//		"id":       "a2",
		//		"price":    45000,
		//		"quantity": 1,
		//		"name":     "Jeruk",
		//	},
		//},
		"echannel": map[string]interface{}{
			"bill_info1": "Payment For:",
			"bill_info2": "debt",
			"bill_key":   "081211111111",
		},
	}
}

// Docs for Mandiri E-Channel (Bill Payment)
/*
JSON Attribute      | Description                                                        | Type     | Required
--------------------|--------------------------------------------------------------------|----------|----------
payment_type        | Set E-channel payment method. Value: `echannel`.                  | String   | Required
transaction_details | The details of the specific transaction such as `order_id` and `gross_amount`. | Object | Required
item_details        | Details of the item(s) purchased by the customer.                 | Object   | Optional
customer_details    | Details of the customer.                                           | Object   | Optional
echannel            | Charge details using Mandiri Bill Payment.                        | Object   | Required
*/
