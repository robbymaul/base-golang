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
	"time"
)

type ShopeePay struct {
	midtrans *Midtrans
}

type ShopeePayResponse struct {
	StatusCode             string `json:"status_code"`
	StatusMessage          string `json:"status_message"`
	ChannelResponseCode    string `json:"channel_response_code"`
	ChannelResponseMessage string `json:"channel_response_message"`
	TransactionId          string `json:"transaction_id"`
	OrderId                string `json:"order_id"`
	MerchantId             string `json:"merchant_id"`
	GrossAmount            string `json:"gross_amount"`
	Currency               string `json:"currency"`
	PaymentType            string `json:"payment_type"`
	TransactionTime        string `json:"transaction_time"`
	TransactionStatus      string `json:"transaction_status"`
	FraudStatus            string `json:"fraud_status"`
	ExpiryTime             string `json:"expiry_time"`
	Actions                []struct {
		Name   string `json:"name"`
		Method string `json:"method"`
		Url    string `json:"url"`
	} `json:"actions"`
}

func NewShopeePay(midtrans *Midtrans) *ShopeePay {
	return &ShopeePay{
		midtrans: midtrans,
	}
}

func (s *ShopeePay) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	paymentReq, ok := req.(PaymentRequest)
	if !ok {
		return nil, errors.New("invalid payment request type")
	}

	// Convert to Midtrans charge request
	chargeReq := s.createChargeRequest(paymentReq)

	result, err := s.midtrans.CreatePayment(ctx, chargeReq)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ShopeePay) ClientResponse(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error) {
	var shopeepayResponse ShopeePayResponse

	err := json.Unmarshal(payment.GatewayResponse, &shopeepayResponse)
	if err != nil {
		return nil, err
	}

	if shopeepayResponse.StatusCode != "201" {
		return nil, helpers.NewErrorTrace(fmt.Errorf(shopeepayResponse.StatusMessage), s.midtrans.serviceName).WithStatusCode(http.StatusPaymentRequired)
	}

	actions := make([]*web.Actions, 0, len(shopeepayResponse.Actions))

	for _, action := range shopeepayResponse.Actions {
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
		PaymentType:   enums.PAYMENT_TYPE_REDIRECT,
		PaymentDetail: web.PaymentDetail{
			Bank:            enums.CHANNEL_SHOPEE,
			Url:             actions,
			VaNumber:        "",
			BillKey:         "",
			BIllCode:        "",
			TransactionTime: shopeepayResponse.TransactionTime,
			ExpireTime:      shopeepayResponse.ExpiryTime,
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

func (s *ShopeePay) createChargeRequest(req PaymentRequest) map[string]interface{} {
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
		"payment_type": "shopeepay",
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": req.Amount,
		},
		"item_details": []map[string]interface{}{
			{
				"id":       time.Now().Unix(),
				"price":    req.Amount,
				"quantity": 1,
				"name":     "Brown sugar boba milk tea",
			},
		},
		"customer_details": map[string]interface{}{
			"first_name": firstName,
			"last_name":  lastName,
			"email":      req.CustomerEmail,
			"phone":      req.CustomerPhone,
		},
		"shopeepay": map[string]interface{}{
			"callback_url": "https://midtrans.com/",
		},
	}
}

// ShopeePay format JSON
// | JSON Attribute       | Description                                                      | Type   | Required |
// |----------------------|------------------------------------------------------------------|--------|----------|
// | payment_type         | Set ShopeePay payment method. Value: shopeepay                  | String | Yes      |
// | transaction_details  | The details of the specific transaction such as order_id and gross_amount | Object | Yes |
// | item_details         | Details of the item(s) purchased by the customer                | Object | Yes      |
// | customer_details     | Details of the customer                                          | Object | Yes      |
// | shopeepay            | Charge details using ShopeePay                                  | Object | Yes      |
