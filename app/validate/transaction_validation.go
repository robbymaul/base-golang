package validate

import (
	"fmt"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/web"
	"time"
)

func ListTransactionRequestValidation(payload *web.ListTransactionRequest) error {
	if payload.Currency != "" {
		enumCurrency := enums.Currency(payload.Currency)
		validCurrency, currencyList := helpers.IsInList(enumCurrency, enums.CURRENCY_IDR, enums.CURRENCY_MYR)
		if !validCurrency {
			return fmt.Errorf("invalid currency parameter request, %v", currencyList)
		}
	}

	if payload.StartDate != "" {
		_, err := time.Parse("2006-01-02", payload.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start date parameter request, %v", err)
		}
	}

	if payload.EndDate != "" {
		_, err := time.Parse("2006-01-02", payload.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end date parameter request, %v", err)
		}
	}

	return nil
}
