package validate

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/zerolog/log"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/web"
)

func ValidationCreateConfigurationRequest(payload []web.CreateConfigurationRequest) error {
	log.Debug().Interface("payload", payload).Msg("validation create configuration request")

	if len(payload) < 1 {
		return errors.New("payload length should be greater than zero, cannot empty")
	}

	for idx1, conf := range payload {
		err := ValidationConfigurationRequest(&conf)
		if err != nil {
			log.Debug().Err(err).Msg("validation create configuration request")
			return fmt.Errorf("%v on index %v", err, idx1)
		}
	}

	return nil
}

func ValidationConfigurationRequest(payload *web.CreateConfigurationRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.ConfigName, validation.Required),
		validation.Field(&payload.Aggregator, validation.Required),
		validation.Field(&payload.ConfigValue, validation.Required),
		validation.Field(&payload.ConfigJson, validation.Required),
	)

	configValueValid, configValueList := helpers.IsInList(payload.ConfigValue, enums.SANDBOX, enums.PRODUCTION)
	if !configValueValid {
		err = fmt.Errorf("invalid config value please call administrator, %v", configValueList)
		return err
	}

	if payload.Aggregator != nil {
		err = validation.ValidateStruct(payload.Aggregator,
			validation.Field(&payload.Aggregator.Id, validation.Required),
		)
		if err != nil {
			return err
		}
	}

	return err
}

func ValidationUpdateConfigurationRequest(payload *web.ResponseConfiguration) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.ConfigName, validation.Required),
		validation.Field(&payload.Aggregator, validation.Required),
		validation.Field(&payload.ConfigValue, validation.Required),
		validation.Field(&payload.ConfigJson, validation.Required),
	)
	if err != nil {
		return err
	}

	configValueValid, configValueList := helpers.IsInList(payload.ConfigValue, enums.SANDBOX, enums.PRODUCTION)
	if !configValueValid {
		err = fmt.Errorf("invalid config value please call administrator, %v", configValueList)
		return err
	}

	if payload.Aggregator != nil {
		err = validation.ValidateStruct(payload.Aggregator,
			validation.Field(&payload.Aggregator.Id, validation.Required),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidationConfigurationPlatformRequest(payload []web.DetailPlatformResponse) error {
	if len(payload) < 1 {
		return errors.New("payload length should be greater than zero, cannot empty")
	}

	seen := make(map[int64]bool)

	for idx1, platform := range payload {
		key := platform.Id

		if seen[platform.Id] {
			return fmt.Errorf("duplicate platform found at index %d: code=%s, name=%s", idx1, platform.Code, platform.Name)
		}

		seen[key] = true
	}

	for idx2, platform := range payload {
		err := ValidationPlatformRequest(&platform)
		if err != nil {
			return fmt.Errorf("%v on index %v", err, idx2)
		}
	}

	return nil
}

func ValidationArrayConfiguration(payload []web.ResponseConfiguration) error {
	seenConf := make(map[int64]bool)

	for idx2, conf := range payload {
		err := ValidationConfiguration(&conf)
		if err != nil {
			return fmt.Errorf("%v on index %v", err, idx2)
		}

		if seenConf[conf.AggregatorId] {
			return fmt.Errorf("duplicate configuration found at index %d: configName=%s, aggregatorId=%d", idx2, conf.ConfigName, conf.AggregatorId)
		}

		seenConf[conf.AggregatorId] = true
	}

	return nil
}

func ValidationConfiguration(payload *web.ResponseConfiguration) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Id, validation.Required),
		validation.Field(&payload.AggregatorId, validation.Required),
	)
	if err != nil {
		return err
	}

	return nil
}
