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
)

type Gopay struct {
	midtrans *Midtrans
}

type GopayResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	OrderId           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	ExpiryTime        string `json:"expiry_time"`
	Actions           []struct {
		Name   string        `json:"name"`
		Method string        `json:"method"`
		Url    string        `json:"url"`
		Fields []interface{} `json:"fields,omitempty"`
	} `json:"actions"`
	ChannelResponseCode    string `json:"channel_response_code"`
	ChannelResponseMessage string `json:"channel_response_message"`
	Currency               string `json:"currency"`
}

func NewGopay(midtrans *Midtrans) *Gopay {
	return &Gopay{
		midtrans: midtrans,
	}
}

func (g *Gopay) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	paymentReq, ok := req.(PaymentRequest)
	if !ok {
		return nil, errors.New("invalid payment request type")
	}

	// Convert to Midtrans charge request
	chargeReq := g.createChargeRequest(paymentReq)

	result, err := g.midtrans.CreatePayment(ctx, chargeReq)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (g *Gopay) ClientResponse(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error) {
	var gopayResponse GopayResponse

	err := json.Unmarshal(payment.GatewayResponse, &gopayResponse)
	if err != nil {
		return nil, err
	}

	if gopayResponse.StatusCode != "201" {
		return nil, helpers.NewErrorTrace(fmt.Errorf(gopayResponse.StatusMessage), g.midtrans.serviceName).WithStatusCode(http.StatusPaymentRequired)
	}

	actions := make([]*web.Actions, 0, len(gopayResponse.Actions))

	for _, action := range gopayResponse.Actions {
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
			Bank:            enums.CHANNEL_GOPAY,
			Url:             actions,
			VaNumber:        "",
			BillKey:         "",
			BIllCode:        "",
			TransactionTime: gopayResponse.TransactionTime,
			ExpireTime:      gopayResponse.ExpiryTime,
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

func (g *Gopay) createChargeRequest(req PaymentRequest) map[string]interface{} {
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
		"payment_type": req.Method,
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
		"customer_details": map[string]interface{}{
			"first_name": firstName,
			"last_name":  lastName,
			"email":      req.CustomerEmail,
			"phone":      req.CustomerPhone,
		},
		"gopay": map[string]interface{}{
			"enable_callback": true,
			"callback_url":    "someapps://callback",
		},
	}
}

// Gopay format JSON
//| Attribute            | Description                                                                 | Type   | Required |
//|----------------------|-----------------------------------------------------------------------------|--------|----------|
//| payment_type         | Set GoPay payment method. Value: gopay.                                     | String | Yes      |
//| transaction_details  | The details of the specific transaction such as order_id and gross_amount. | Object | Yes      |
//| item_details         | Details of the item(s) purchased by the customer.                           | Object | No       |
//| customer_details     | Details of the customer.                                                     | Object | No       |
//| gopay                | Charge details using GoPay.                                                  | Object | Yes      |

//{
//	"payment_type": req.PaymentType,
//	"transaction_details": map[string]interface{}{
//		"order_id":     req.OrderID,
//		"gross_amount": req.Amount,
//	},
//	"customer_details": map[string]interface{}{
//		"first_name": req.CustomerName,
//		"email":      req.CustomerEmail,
//		"phone":      req.CustomerPhone,
//	},
//	"item_details": []map[string]interface{}{
//		{
//			"id":       req.OrderID,
//			"price":    req.Amount,
//			"quantity": 1,
//			"name":     req.Description,
//		},
//	},
//}
