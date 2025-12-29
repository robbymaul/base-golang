package strategy

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/web"
)

type PaymentStrategy interface {
	Pay(ctx context.Context, req any) (map[string]interface{}, error)
	CheckStatusPayment(ctx context.Context, req any) (any, error)
	CheckKey() (any, error)
	MapResponsePayment(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error)
	MapCheckStatusPayment(payment *models.Payments, statusPayment any) (*web.CheckStatusPaymentResponse, error)
}

type Strategy struct {
	SenangPay PaymentStrategy
	Midtrans  PaymentStrategy
	Espay     PaymentStrategy
}

func NewStrategy(senangPay PaymentStrategy, midtrans PaymentStrategy, espay PaymentStrategy) *Strategy {
	return &Strategy{SenangPay: senangPay, Midtrans: midtrans, Espay: espay}
}

func (s *Strategy) GetStrategy(configuration *models.Configuration) (PaymentStrategy, error) {
	log.Debug().Interface("configuration", configuration).Msg("get strategy")
	switch configuration.Aggregator.Slug {
	case enums.PROVIDER_PAYMENT_METHOD_MIDTRANS:
		return s.Midtrans, nil
	case enums.PROVIDER_PAYMENT_METHOD_SENANGPAY:
		return s.SenangPay, nil
	case enums.PROVIDER_PAYMENT_METHOD_ESPAY:
		return s.Espay, nil
	default:
		return nil, fmt.Errorf("unknown payment method")
	}
}
