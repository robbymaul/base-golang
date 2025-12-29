package validate

import (
	"fmt"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/web"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/zerolog/log"
)

func ValidationCreateAggregatorRequest(payload *web.CreateAggregatorRequest) error {
	log.Debug().Interface("payload", payload).Msg("validation create aggregator request")
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Name, validation.Required),
		validation.Field(&payload.Currency),
	)

	validAggregatorName, listAggregatorName := helpers.IsInList(
		payload.Name,
		enums.AGGREGATOR_NAME_MIDTRANS,
		enums.AGGREGATOR_NAME_ESPAY,
		enums.AGGREGATOR_NAME_SENANGPAY,
	)
	if !validAggregatorName {
		err = fmt.Errorf("invalid aggregator name please call administrator, %v", listAggregatorName)
		return err
	}

	validCurrency, currencyList := helpers.IsInList(payload.Currency, []enums.Currency{enums.CURRENCY_IDR, enums.CURRENCY_MYR}...)
	if !validCurrency {
		return fmt.Errorf("currency is not valid, example currency valid %v", currencyList)
	}

	//if len(payload.Channel) > 0 {
	//	err = ValidationCreateChannelRequest(payload.Channel)
	//	if err != nil {
	//		return err
	//	}
	//}

	if len(payload.Configuration) > 0 {
		err = ValidationCreateConfigurationRequest(payload.Configuration)
		if err != nil {
			return err
		}
	}

	return err
}

func ValidationUpdateAggregatorRequest(payload *web.AggregatorResponse) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Name, validation.Required),
	)

	validAggregatorName, listAggregatorName := helpers.IsInList(
		payload.Name,
		enums.AGGREGATOR_NAME_MIDTRANS,
		enums.AGGREGATOR_NAME_ESPAY,
		enums.AGGREGATOR_NAME_SENANGPAY,
	)
	if !validAggregatorName {
		err = fmt.Errorf("invalid aggregator name please call administrator, %v", listAggregatorName)
		return err
	}

	return err
}

func ValidationAggregatorQuery(aggregator string) error {
	if aggregator == "" {
		return nil
	}

	agg := enums.AggregatorName(aggregator)

	validAggregatorName, listAggregatorName := helpers.IsInList(
		agg,
		enums.AGGREGATOR_NAME_MIDTRANS,
		enums.AGGREGATOR_NAME_ESPAY,
		enums.AGGREGATOR_NAME_SENANGPAY,
	)
	if !validAggregatorName {
		err := fmt.Errorf("invalid aggregator name please call administrator, %v", listAggregatorName)
		return err
	}

	return nil
}
