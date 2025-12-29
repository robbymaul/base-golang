package midtrans

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/web"
	"strings"

	"github.com/rs/zerolog/log"
)

type Dana struct {
	midtrans *Midtrans
}

type DanaResponseSnap struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func NewDana(midtrans *Midtrans) *Dana {
	return &Dana{
		midtrans: midtrans,
	}
}

func (d *Dana) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	paymentReq, ok := req.(PaymentRequest)
	if !ok {
		return nil, errors.New("invalid payment request type")
	}

	// Convert to Midtrans charge request
	chargeReq := d.createChargeRequest(paymentReq)
	log.Debug().Interface("chargeReq", chargeReq).Msg("chargeReq pay()")

	result, err := d.midtrans.CreateSnapPayment(ctx, chargeReq)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Dana) ClientResponse(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error) {
	var danaResponse DanaResponseSnap

	err := json.Unmarshal(payment.GatewayResponse, &danaResponse)
	if err != nil {
		return nil, err
	}

	//if briVaResponse.StatusCode != "201" {
	//	return nil, helpers.NewErrorTrace(fmt.Errorf(briVaResponse.StatusMessage), b.midtrans.serviceName).WithStatusCode(http.StatusPaymentRequired)
	//}

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
			Bank: enums.CHANNEL_DANA,
			Url: []*web.Actions{{
				Name:   "snap url",
				Method: http.MethodGet,
				Url:    danaResponse.RedirectURL,
			}},
			//VaNumber:        briVaResponse.VaNumbers[0].VaNumber,
			//TransactionTime: briVaResponse.TransactionTime,
			//ExpireTime:      briVaResponse.ExpiryTime,
			Instruction: channel.Instruction,
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

func (d *Dana) createChargeRequest(req PaymentRequest) map[string]interface{} {
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

	//return map[string]interface{}{
	//	//"payment_type": "qris",
	//	"transaction_details": map[string]interface{}{
	//		"order_id":     req.OrderID,
	//		"gross_amount": req.Amount,
	//	},
	//	//"item_details": []map[string]interface{}{
	//	//	{
	//	//		"id":       "id1",
	//	//		"price":    275000,
	//	//		"quantity": 1,
	//	//		"name":     "Bluedio H+ Turbine Headphone with Bluetooth 4.1 -",
	//	//	},
	//	//},
	//	"payment_type": "ewallet",
	//	"customer_details": map[string]interface{}{
	//		"first_name": firstName,
	//		"last_name":  lastName,
	//		"email":      req.CustomerEmail,
	//		"phone":      req.CustomerPhone,
	//	},
	//	"enabled_payments": []string{"dana"},
	//	"dana": map[string]string{
	//		"callback": req.CallbackURL,
	//	},
	//}

	return map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": req.Amount,
		},
		"enabled_payments": &[]interface{}{
			"dana",
		},
		"customer_details": map[string]interface{}{
			"first_name": firstName,
			"last_name":  lastName,
			"email":      req.CustomerEmail,
			"phone":      req.CustomerPhone,
		},
		"callbacks": map[string]string{
			"finish": req.CallbackURL,
		},
	}
}
