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
)

type PermataVa struct {
	midtrans *Midtrans
}

func NewPermataVa(midtrans *Midtrans) *PermataVa {
	return &PermataVa{
		midtrans: midtrans,
	}
}

type PermataVAResponse struct {
	Currency          string `json:"currency"`
	OrderId           string `json:"order_id"`
	ExpiryTime        string `json:"expiry_time"`
	MerchantId        string `json:"merchant_id"`
	StatusCode        string `json:"status_code"`
	FraudStatus       string `json:"fraud_status"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	StatusMessage     string `json:"status_message"`
	TransactionId     string `json:"transaction_id"`
	TransactionTime   string `json:"transaction_time"`
	PermataVaNumber   string `json:"permata_va_number"`
	TransactionStatus string `json:"transaction_status"`
}

func (p *PermataVa) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	paymentReq, ok := req.(PaymentRequest)
	if !ok {
		return nil, errors.New("invalid payment request type")
	}

	// Convert to Midtrans charge request
	chargeReq := p.createChargeRequest(paymentReq)

	result, err := p.midtrans.CreatePayment(ctx, chargeReq)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PermataVa) ClientResponse(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error) {
	var permataVaResponse PermataVAResponse

	err := json.Unmarshal(payment.GatewayResponse, &permataVaResponse)
	if err != nil {
		return nil, err
	}

	if permataVaResponse.StatusCode != "201" {
		return nil, helpers.NewErrorTrace(fmt.Errorf(permataVaResponse.StatusMessage), p.midtrans.serviceName).WithStatusCode(http.StatusPaymentRequired)
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
			Bank:            enums.CHANNEL_PERMATA,
			Url:             nil,
			VaNumber:        permataVaResponse.PermataVaNumber,
			TransactionTime: permataVaResponse.TransactionTime,
			ExpireTime:      permataVaResponse.ExpiryTime,
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

func (p *PermataVa) createChargeRequest(req PaymentRequest) map[string]interface{} {
	return map[string]interface{}{
		"payment_type": "bank_transfer",
		"bank_transfer": map[string]interface{}{
			"bank": "permata",
			"permata": map[string]interface{}{
				"recipient_name": p.midtrans.MerchantId,
			},
		},
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": req.Amount,
		},
	}
}

// Docs for Permata VA
/*
JSON Attribute      | Description                                                        | Type     | Required
--------------------|--------------------------------------------------------------------|----------|----------
payment_type        | Set Bank Transfer payment method. Value: `bank_transfer`.         | String   | Required
bank_transfer       | Charge details using bank transfer.                               | Object   | Required
transaction_details | The details of the specific transaction such as `order_id` and `gross_amount`. | Object | Required
item_details        | Details of the item(s) purchased by the customer.                 | Object   | Optional
customer_details    | Details of the customer.                                           | Object   | Optional
*/
