package validate

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"paymentserviceklink/app/web"
)

func ValidationCreatePlatformRequest(payload *web.CreatePlatformRequest) error {
	err := validation.ValidateStruct(payload,
		//validation.Field(&payload.Code, validation.Required),
		validation.Field(&payload.Name, validation.Required),
		validation.Field(&payload.Description, validation.Required),
	)

	//validCode, codeList := helpers.IsInList(payload.Code, enums.CODE_PLATFORM_KNET, enums.CODE_PLATFORM_SMS)
	//if !validCode {
	//	err = fmt.Errorf("invalid code platform please call administrator, %v", codeList)
	//}

	seenCh := make(map[string]bool)

	for idx1, ch := range payload.Channel {
		err = ValidationChannel(ch)
		if err != nil {
			return fmt.Errorf("%v on index %v", err, idx1)
		}

		// bikin key gabungan dari field yg harus unik
		key := fmt.Sprintf("%s|%s|%s|%s", ch.Code, ch.PaymentMethod, ch.Currency, ch.BankName)

		if seenCh[key] {
			return fmt.Errorf("duplicate channel found at index %d: code=%s, paymentMethod=%s, currency=%s, bankName=%s",
				idx1, ch.Code, ch.PaymentMethod, ch.Currency, ch.BankName)
		}
		seenCh[key] = true

	}

	seenConf := make(map[int64]bool)

	for idx2, conf := range payload.Configuration {
		err = ValidationConfiguration(conf)
		if err != nil {
			return fmt.Errorf("%v on index %v", err, idx2)
		}

		if seenConf[conf.AggregatorId] {
			return fmt.Errorf("duplicate configuration found at index %d: configName=%s, aggregatorId=%d", idx2, conf.ConfigName, conf.AggregatorId)
		}

		seenConf[conf.AggregatorId] = true
	}

	return err
}

func ValidationUpdatePlatformRequest(payload *web.DetailPlatformResponse) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Id, validation.Required),
		//validation.Field(&payload.Code, validation.Required),
		validation.Field(&payload.Name, validation.Required),
		validation.Field(&payload.Description, validation.Required),
	)

	//validCode, codeList := helpers.IsInList(payload.Code, enums.CODE_PLATFORM_KNET, enums.CODE_PLATFORM_SMS)
	//if !validCode {
	//	err = fmt.Errorf("invalid code platform please call administrator, %v", codeList)
	//}

	return err
}

func ValidationPlatformRequest(payload *web.DetailPlatformResponse) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Id, validation.Required),
	)

	return err
}
