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

type BriVa struct {
	midtrans *Midtrans
}

type BRIVAResponse struct {
	Currency  string `json:"currency"`
	OrderId   string `json:"order_id"`
	VaNumbers []struct {
		Bank     string `json:"bank"`
		VaNumber string `json:"va_number"`
	} `json:"va_numbers"`
	ExpiryTime        string `json:"expiry_time"`
	MerchantId        string `json:"merchant_id"`
	StatusCode        string `json:"status_code"`
	FraudStatus       string `json:"fraud_status"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
}

func NewBriVa(midtrans *Midtrans) *BriVa {
	return &BriVa{
		midtrans: midtrans,
	}
}

func (b *BriVa) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	paymentReq, ok := req.(PaymentRequest)
	if !ok {
		return nil, errors.New("invalid payment request type")
	}

	// Convert to Midtrans charge request
	chargeReq := b.createChargeRequest(paymentReq)

	result, err := b.midtrans.CreatePayment(ctx, chargeReq)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b *BriVa) ClientResponse(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error) {
	var briVaResponse BRIVAResponse

	err := json.Unmarshal(payment.GatewayResponse, &briVaResponse)
	if err != nil {
		return nil, err
	}

	if briVaResponse.StatusCode != "201" {
		return nil, helpers.NewErrorTrace(fmt.Errorf(briVaResponse.StatusMessage), b.midtrans.serviceName).WithStatusCode(http.StatusPaymentRequired)
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
		PaymentType:   enums.PAYMENT_TYPE_VA,
		PaymentDetail: web.PaymentDetail{
			Bank:            enums.Channel(briVaResponse.VaNumbers[0].Bank),
			Url:             nil,
			VaNumber:        briVaResponse.VaNumbers[0].VaNumber,
			TransactionTime: briVaResponse.TransactionTime,
			ExpireTime:      briVaResponse.ExpiryTime,
			Instruction:     channel.Instruction,
		},
		Customer: web.Customer{
			MemberId: payment.CustomerId,
			Name:     payment.CustomerName,
			Email:    payment.CustomerEmail,
			Phone:    payment.CustomerPhone,
		},
		CreatedAt: nil,
		UpdatedAt: nil,
	}, nil
}

func (b *BriVa) createChargeRequest(req PaymentRequest) map[string]interface{} {
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
		"payment_type": "bank_transfer",
		"transaction_details": map[string]interface{}{
			"gross_amount": req.Amount,
			"order_id":     req.OrderID,
		},
		"customer_details": map[string]interface{}{
			"email":      req.CustomerEmail,
			"first_name": firstName,
			"last_name":  lastName,
			"phone":      req.CustomerPhone,
		},
		//"item_details": []map[string]interface{}{
		//	{
		//		"id":       "1388998298204",
		//		"price":    5000,
		//		"quantity": 1,
		//		"name":     "Ayam Zozozo",
		//	},
		//	{
		//		"id":       "1388998298205",
		//		"price":    5000,
		//		"quantity": 1,
		//		"name":     "Ayam Xoxoxo",
		//	},
		//},
		"bank_transfer": map[string]interface{}{
			"bank": "bri",
			//"va_number": "111111",
		},
	}
}

// BRI VA JSON Request Format
/*
   JSON Attribute       | Description                                                                 | Type   | Required
   ---------------------|-----------------------------------------------------------------------------|--------|---------
   payment_type         | Set Bank Transfer payment method. Value: bank_transfer                     | String | Required
   transaction_details  | The details of the specific transaction such as order_id and gross_amount. | Object | Required
   customer_details     | Details of the customer.                                                    | Object | Optional
   item_details         | Details of the item(s) purchased by the customer.                          | Object | Optional
   bank_transfer        | Charge details using bank transfer.                                        | Object | Required
*/
