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
	"strings"

	"github.com/rs/zerolog/log"
)

type Qris struct {
	midtrans *Midtrans
}

type QRISResponse struct {
	Acquirer string `json:"acquirer"`
	Actions  []struct {
		Method string `json:"method"`
		Name   string `json:"name"`
		Url    string `json:"url"`
	} `json:"actions"`
	Currency          string `json:"currency"`
	ExpiryTime        string `json:"expiry_time"`
	FraudStatus       string `json:"fraud_status"`
	GrossAmount       string `json:"gross_amount"`
	MerchantId        string `json:"merchant_id"`
	OrderId           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	QrString          string `json:"qr_string"`
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	TransactionTime   string `json:"transaction_time"`
}

func NewQris(midtrans *Midtrans) *Qris {
	return &Qris{
		midtrans: midtrans,
	}
}

func (q *Qris) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	log.Debug().Interface("req", req).Msg("pay qris")
	paymentReq, ok := req.(PaymentRequest)
	if !ok {
		return nil, errors.New("invalid payment request type")
	}

	// Convert to Midtrans charge request
	chargeReq := q.createChargeRequest(paymentReq)
	log.Debug().Interface("chargeReq", chargeReq).Msg("charge request")

	result, err := q.midtrans.CreatePayment(ctx, chargeReq)
	if err != nil {
		return nil, err
	}
	log.Debug().Interface("result", result).Msg("result pay qris")

	return result, nil
}

func (q *Qris) ClientResponse(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error) {
	var qrisResponse QRISResponse

	err := json.Unmarshal(payment.GatewayResponse, &qrisResponse)
	if err != nil {
		return nil, err
	}

	if qrisResponse.StatusCode != "201" {
		return nil, helpers.NewErrorTrace(fmt.Errorf(qrisResponse.StatusMessage), q.midtrans.serviceName).WithStatusCode(http.StatusPaymentRequired)
	}

	actions := make([]*web.Actions, 0, len(qrisResponse.Actions))

	for _, action := range qrisResponse.Actions {
		actions = append(actions, &web.Actions{
			Name:   action.Name,
			Method: action.Method,
			Url:    action.Url,
		})
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
		PaymentType:   enums.PAYMENT_TYPE_QRIS,
		PaymentDetail: web.PaymentDetail{
			Bank:            enums.CHANNEL_QRIS,
			Url:             actions,
			VaNumber:        "",
			BillKey:         "",
			BIllCode:        "",
			TransactionTime: qrisResponse.TransactionTime,
			ExpireTime:      qrisResponse.ExpiryTime,
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

func (q *Qris) createChargeRequest(req PaymentRequest) map[string]interface{} {
	names := strings.Split(req.CustomerName, " ")
	var firstName, lastName string
	for i, name := range names {
		if i == 0 {
			firstName = name
		}
		if i == 1 {
			lastName = name
		} else {
			lastName += " " + name
		}
	}

	return map[string]interface{}{
		//"payment_type": "qris",
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": req.Amount,
		},
		//"item_details": []map[string]interface{}{
		//	{
		//		"id":       "id1",
		//		"price":    275000,
		//		"quantity": 1,
		//		"name":     "Bluedio H+ Turbine Headphone with Bluetooth 4.1 -",
		//	},
		//},
		"payment_type": "qris",
		"customer_details": map[string]interface{}{
			"first_name": firstName,
			"last_name":  lastName,
			"email":      req.CustomerEmail,
			"phone":      req.CustomerPhone,
		},
		"enabled_payments": []string{"other_qris"},
	}
}

// QRIS format JSON
//| Attribute             | Description                                                                        | Type   | Required |
//|-----------------------|------------------------------------------------------------------------------------|--------|----------|
//| payment_type          | Set Mandiri bill payment method. Value: `echannel`                                | String | Required |
//| transaction_details   | The details of the specific transaction such as `order_id` and `gross_amount`.    | Object | Required |
//| item_details          | Details of the item(s) purchased by the customer.                                 | Object | Optional |
//| customer_details      | Details of the customer.                                                           | Object | Optional |
//| echannel              | Charge details using Mandiri Bill Payment, including `bill_info1` and `bill_info2`. | Object | Required |
