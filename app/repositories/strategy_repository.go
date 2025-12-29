package repositories

import (
	"context"
	"paymentserviceklink/app/strategy"
)

func (rc *RepositoryContext) PaymentStrategyRepository(ctx context.Context, strategy strategy.PaymentStrategy, req any) (any, error) {
	result, err := strategy.Pay(ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}
