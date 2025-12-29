package strategy

import (
	"context"
	"fmt"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/web"
)

type EspayPayment interface {
	Pay(ctx context.Context, req any) (map[string]interface{}, error)
	ClientResponse(payment *models.Payments, channel *models.Channel) (*web.PaymentResponse, error)
}

type EspayStrategy struct {
	VirtualAccount EspayPayment
	CreditCard     EspayPayment
	Qris           EspayPayment
	PaymentLink    EspayPayment
}

func NewEspayStrategy(virtualAccount EspayPayment, creditCard EspayPayment, qris EspayPayment, paymentLink EspayPayment) *EspayStrategy {
	return &EspayStrategy{
		VirtualAccount: virtualAccount,
		CreditCard:     creditCard,
		Qris:           qris,
		PaymentLink:    paymentLink,
	}
}

func (e *EspayStrategy) GetStrategy(method enums.PaymentMethod) (EspayPayment, error) {
	switch method {
	case enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT:
		return e.VirtualAccount, nil
	case enums.PAYMENT_METHOD_CREDIT_CARD:
		return e.CreditCard, nil
	case enums.PAYMENT_METHOD_QRIS:
		return e.Qris, nil
	default:
		return nil, fmt.Errorf("unknown payment method: %s", method)
	}
}
