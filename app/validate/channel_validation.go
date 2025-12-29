package validate

import (
	"fmt"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/web"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/zerolog/log"
)

func ValidationCreateChannelRequest(payload []web.CreateChannelRequest) error {
	var err error
	log.Debug().Interface("payload", payload).Msg("validation create channel request")
	for i, ch := range payload {
		err = ValidationCreateChannel(&ch)
		if err != nil {
			log.Debug().Err(err).Msg("validation create channel")
			return fmt.Errorf("%v, on index %v", err, i)
		}
	}

	return err
}

func ValidationCreateChannel(payload *web.CreateChannelRequest) error {
	log.Debug().Interface("payload", payload).Msg("validation create channel request")

	err := validation.ValidateStruct(payload,
		//validation.Field(&payload.Code, validation.Required),
		validation.Field(&payload.Name, validation.Required),
		validation.Field(&payload.PaymentMethod, validation.Required),
		//validation.Field(&payload.Provider, validation.Required),
		validation.Field(&payload.Currency, validation.Required),
		validation.Field(&payload.FeeType, validation.Required),
		validation.Field(&payload.BankName, validation.Required),
		validation.Field(&payload.ProductName, validation.Required),
		validation.Field(&payload.ProductCode, validation.Required),
		validation.Field(&payload.BankCode, validation.Required),
		validation.Field(&payload.BankName, validation.Required),
		//validation.Field(&payload.Aggregator, validation.Required),
	)
	if err != nil {
		return err
	}

	if payload.FeeType == enums.FEE_TYPE_PERCENTAGE {
		err = validation.ValidateStruct(payload,
			validation.Field(&payload.FeePercentage, validation.Required),
		)
		if err != nil {
			return err
		}
	}

	if payload.FeeType == enums.FEE_TYPE_FIXED {
		err = validation.ValidateStruct(payload,
			validation.Field(&payload.FeeFixed, validation.Required),
		)
		if err != nil {
			return err
		}
	}

	if payload.FeeType == enums.FEE_TYPE_FIXED_PERCENTAGE {
		err = validation.ValidateStruct(payload,
			validation.Field(&payload.FeeFixed, validation.Required),
			validation.Field(&payload.FeePercentage, validation.Required),
		)
		if err != nil {
			return err
		}
	}

	bankNameValid, bankNameList := helpers.IsInList(payload.BankName,
		enums.CHANNEL_BCA,
		enums.CHANNEL_BNI,
		enums.CHANNEL_BRI,
		enums.CHANNEL_CIMB,
		enums.CHANNEL_MANDIRI,
		enums.CHANNEL_PERMATA,
		enums.CHANNEL_GOPAY,
		enums.CHANNEL_SHOPEE,
		enums.CHANNEL_SENANGPAY,
		enums.CHANNEL_DANAMON,
		enums.CHANNEL_MAYBANK,
		enums.CHANNEL_K_WALLET,
		enums.CHANNEL_QRIS,
		enums.CHANNEL_DANA,
	)
	if !bankNameValid {
		err = fmt.Errorf("invalid bank name please, %v", bankNameList)
		return err
	}

	//codeValid, codeList := helpers.IsInList(payload.Code,
	//	enums.CODE_PAYMENT_METHOD_VA,
	//	enums.CODE_PAYMENT_METHOD_BANK_TRANSFER,
	//	//enums.CODE_PAYMENT_METHOD_E_WALLET,
	//	//enums.CODE_PAYMENT_METHOD_CREDIT_CARD,
	//	enums.CODE_PAYMENT_METHOD_QRIS,
	//)
	//if !codeValid {
	//	err = fmt.Errorf("invalid code payment method please, %v", codeList)
	//	return err
	//}

	//providerValid, providerList := helpers.IsInList(
	//	payload.Provider,
	//	enums.PROVIDER_PAYMENT_METHOD_MIDTRANS,
	//	enums.PROVIDER_PAYMENT_METHOD_SENANGPAY,
	//	enums.PROVIDER_PAYMENT_METHOD_ESPAY,
	//)
	//if !providerValid {
	//	err = fmt.Errorf("invalid provider payment method please, %v", providerList)
	//	return err
	//}

	currencyValid, currencyList := helpers.IsInList(payload.Currency, enums.CURRENCY_MYR, enums.CURRENCY_IDR)
	if !currencyValid {
		err = fmt.Errorf("invalid currency payment method please, %v", currencyList)
		return err
	}

	feeTypeValid, feeTypeList := helpers.IsInList(payload.FeeType, enums.FEE_TYPE_FIXED, enums.FEE_TYPE_PERCENTAGE, enums.FEE_TYPE_NONE, enums.FEE_TYPE_FIXED_PERCENTAGE)
	if !feeTypeValid {
		err = fmt.Errorf("invalid fee type please, %v", feeTypeList)
		return err
	}

	paymentMethodValid, paymentMethodList := helpers.IsInList(
		payload.PaymentMethod,
		//enums.PAYMENT_METHOD_BCA_VA,
		//enums.PAYMENT_METHOD_BNI_VA,
		//enums.PAYMENT_METHOD_BRI_VA,
		//enums.PAYMENT_METHOD_CIMB_VA,
		//enums.PAYMENT_METHOD_GOPAY,
		//enums.PAYMENT_METHOD_MANDIRI_VA,
		//enums.PAYMENT_METHOD_PERMATA_VA,
		//enums.PAYMENT_METHOD_SHOPEE_PAY,
		//enums.PAYMENT_METHOD_VA,
		enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT,
		enums.PAYMENT_METHOD_CREDIT_CARD,
		enums.PAYMENT_METHOD_QRIS,
		enums.PAYMENT_METHOD_GOPAY,
		enums.PAYMENT_METHOD_BANK_TRANSFER,
		enums.PAYMENT_METEHOD_K_WALLET,
		enums.PAYMENT_METHOD_E_WALLET,
		enums.PAYMENT_METHOD_RETAIL_OUTLET,
		enums.PAYMENT_METHOD_SENANGPAY,
		enums.PAYMENT_METHOD_QRIS,
	)
	if !paymentMethodValid {
		err = fmt.Errorf("invalid payment method please, %v", paymentMethodList)
		return err
	}

	//if payload.IsEspay {
	//	err = validation.ValidateStruct(payload,
	//		validation.Field(&payload.ProductName, validation.Required),
	//		validation.Field(&payload.ProductCode, validation.Required),
	//		validation.Field(&payload.BankCode, validation.Required),
	//	)
	//	return err
	//}

	for _, image := range payload.Image {
		err = validation.ValidateStruct(image,
			validation.Field(&image.FileName, validation.Required),
			validation.Field(&image.FileUrl, validation.Required),
			//validation.Field(&image.SizeType, validation.Required),
			//validation.Field(&image.Geometric, validation.Required),
		)

		if image.SizeType != "" {
			sizeValid, sizeList := helpers.IsInList(image.SizeType, enums.IMAGE_SIZE_TYPE_S, enums.IMAGE_SIZE_TYPE_M, enums.IMAGE_SIZE_TYPE_L, enums.IMAGE_SIZE_TYPE_XL)
			if !sizeValid {
				err = fmt.Errorf("invalid size image please, %v", sizeList)
				return err
			}
		}

		if image.Geometric != "" {
			geometricValid, geometricList := helpers.IsInList(image.Geometric, enums.IMAGE_CIRCLE, enums.IMAGE_ROUNDED, enums.IMAGE_SQUARE, enums.IMAGE_RECTANGLE)
			if !geometricValid {
				err = fmt.Errorf("invalid geometric image please, %v", geometricList)
				return err
			}
		}
	}

	return err
}

func GetListChannelRequest(payload *web.GetListChannelRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Currency, validation.Required),
		validation.Field(&payload.MemberId, validation.Required),
	)

	currencyValid, currencyList := helpers.IsInList(payload.Currency, enums.CURRENCY_MYR, enums.CURRENCY_IDR)
	if !currencyValid {
		err = fmt.Errorf("invalid currency payment method please, %v", currencyList)
		return err
	}

	return err
}

// ValidationUpdateChannelRequest validates the update request for a payment method.
// It ensures that required fields such as Id, Code, Name, Provider, Currency, and FeeType
// are present. Additionally, it checks if the Code, Provider, Currency, and FeeType values
// are within allowed lists. If the payment method is of type Espay, it also validates that
// ProductName, ProductCode, and BankCode are provided. Returns an error if any validation
// fails, otherwise returns nil.
func ValidationUpdateChannelRequest(payload *web.DetailChannelResponse) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Id, validation.Required),
		validation.Field(&payload.Code, validation.Required),
		validation.Field(&payload.Name, validation.Required),
		//validation.Field(&payload.Provider, validation.Required),
		validation.Field(&payload.Currency, validation.Required),
		validation.Field(&payload.FeeType, validation.Required),
		//validation.Field(&payload.FeeAmount, validation.Required),
		validation.Field(&payload.BankName, validation.Required),
	)

	if payload.FeeType == enums.FEE_TYPE_PERCENTAGE {
		err = validation.ValidateStruct(payload,
			validation.Field(&payload.FeePercentage, validation.Required),
		)
		if err != nil {
			return err
		}
	}

	if payload.FeeType == enums.FEE_TYPE_FIXED {
		err = validation.ValidateStruct(payload,
			validation.Field(&payload.FeeFixed, validation.Required),
		)
	}

	if payload.FeeType == enums.FEE_TYPE_FIXED_PERCENTAGE {
		err = validation.ValidateStruct(payload,
			validation.Field(&payload.FeeFixed, validation.Required),
			validation.Field(&payload.FeePercentage, validation.Required),
		)
	}

	bankNameValid, bankNameList := helpers.IsInList(payload.BankName,
		enums.CHANNEL_BCA,
		enums.CHANNEL_BNI,
		enums.CHANNEL_BRI,
		enums.CHANNEL_CIMB,
		enums.CHANNEL_MANDIRI,
		enums.CHANNEL_PERMATA,
		enums.CHANNEL_GOPAY,
		enums.CHANNEL_SHOPEE,
		enums.CHANNEL_SENANGPAY,
		enums.CHANNEL_DANAMON,
		enums.CHANNEL_MAYBANK,
		enums.CHANNEL_K_WALLET,
	)
	if !bankNameValid {
		err = fmt.Errorf("invalid bank name please, %v", bankNameList)
		return err
	}

	//codeValid, codeList := helpers.IsInList(payload.Code,
	//	enums.CODE_PAYMENT_METHOD_VA,
	//	enums.CODE_PAYMENT_METHOD_BANK_TRANSFER,
	//	//enums.CODE_PAYMENT_METHOD_E_WALLET,
	//	//enums.CODE_PAYMENT_METHOD_CREDIT_CARD,
	//	enums.CODE_PAYMENT_METHOD_QRIS,
	//)
	//if !codeValid {
	//	err = fmt.Errorf("invalid code payment method please, %v", codeList)
	//	return err
	//}

	//providerValid, providerList := helpers.IsInList(
	//	payload.Provider,
	//	enums.PROVIDER_PAYMENT_METHOD_SENANGPAY,
	//	enums.PROVIDER_PAYMENT_METHOD_MIDTRANS,
	//	enums.PROVIDER_PAYMENT_METHOD_ESPAY,
	//)
	//if !providerValid {
	//	err = fmt.Errorf("invalid provider payment method please, %v", providerList)
	//	return err
	//}

	currencyValid, currencyList := helpers.IsInList(payload.Currency, enums.CURRENCY_IDR, enums.CURRENCY_MYR)
	if !currencyValid {
		err = fmt.Errorf("invalid currency payment method please, %v", currencyList)
		return err
	}

	feeTypeValid, feeTypeList := helpers.IsInList(payload.FeeType, enums.FEE_TYPE_PERCENTAGE, enums.FEE_TYPE_NONE, enums.FEE_TYPE_FIXED)
	if !feeTypeValid {
		err = fmt.Errorf("invalid fee type please, %v", feeTypeList)
		return err
	}

	paymentMethodValid, paymentMethodList := helpers.IsInList(
		payload.PaymentMethod,
		//enums.PAYMENT_METHOD_BCA_VA,
		//enums.PAYMENT_METHOD_BNI_VA,
		//enums.PAYMENT_METHOD_BRI_VA,
		//enums.PAYMENT_METHOD_CIMB_VA,
		//enums.PAYMENT_METHOD_GOPAY,
		//enums.PAYMENT_METHOD_MANDIRI_VA,
		//enums.PAYMENT_METHOD_PERMATA_VA,
		//enums.PAYMENT_METHOD_SHOPEE_PAY,
		//enums.PAYMENT_METHOD_VA,
		enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT,
		enums.PAYMENT_METHOD_CREDIT_CARD,
		enums.PAYMENT_METHOD_QRIS,
		enums.PAYMENT_METHOD_BANK_TRANSFER,
		enums.PAYMENT_METHOD_GOPAY,
		enums.PAYMENT_METEHOD_K_WALLET,
		enums.PAYMENT_METHOD_E_WALLET,
		enums.PAYMENT_METHOD_RETAIL_OUTLET,
		enums.PAYMENT_METHOD_SENANGPAY,
	)
	if !paymentMethodValid {
		err = fmt.Errorf("invalid payment method please, %v", paymentMethodList)
		return err
	}

	//if payload.IsEspay {
	//	err = validation.ValidateStruct(payload,
	//		validation.Field(&payload.ProductName, validation.Required),
	//		validation.Field(&payload.ProductCode, validation.Required),
	//		validation.Field(&payload.BankCode, validation.Required),
	//	)
	//}

	return err
}

func ValidationChannel(payload *web.DetailChannelResponse) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Id, validation.Required),
		//va/\lidation.Field(&payload.Code, validation.Required),
		validation.Field(&payload.PaymentMethod, validation.Required),
		validation.Field(&payload.Currency, validation.Required),
		validation.Field(&payload.BankName, validation.Required),
	)
	if err != nil {
		return err
	}

	//codeValid, codeList := helpers.IsInList(payload.Code,
	//	enums.CODE_PAYMENT_METHOD_VA,
	//	enums.CODE_PAYMENT_METHOD_BANK_TRANSFER,
	//	//enums.CODE_PAYMENT_METHOD_E_WALLET,
	//	//enums.CODE_PAYMENT_METHOD_CREDIT_CARD,
	//	enums.CODE_PAYMENT_METHOD_QRIS,
	//)
	//if !codeValid {
	//	err = fmt.Errorf("invalid code payment method please, %v", codeList)
	//	return err
	//}

	currencyValid, currencyList := helpers.IsInList(payload.Currency, enums.CURRENCY_IDR, enums.CURRENCY_MYR)
	if !currencyValid {
		err = fmt.Errorf("invalid currency payment method please, %v", currencyList)
		return err
	}

	paymentMethodValid, paymentMethodList := helpers.IsInList(
		payload.PaymentMethod,
		//enums.PAYMENT_METHOD_BCA_VA,
		//enums.PAYMENT_METHOD_BNI_VA,
		//enums.PAYMENT_METHOD_BRI_VA,
		//enums.PAYMENT_METHOD_CIMB_VA,
		//enums.PAYMENT_METHOD_GOPAY,
		//enums.PAYMENT_METHOD_MANDIRI_VA,
		//enums.PAYMENT_METHOD_PERMATA_VA,
		//enums.PAYMENT_METHOD_SHOPEE_PAY,
		//enums.PAYMENT_METHOD_VA,
		enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT,
		enums.PAYMENT_METHOD_CREDIT_CARD,
		enums.PAYMENT_METHOD_QRIS,
		enums.PAYMENT_METHOD_GOPAY,
		enums.PAYMENT_METHOD_BANK_TRANSFER,
		enums.PAYMENT_METEHOD_K_WALLET,
		enums.PAYMENT_METHOD_E_WALLET,
		enums.PAYMENT_METHOD_RETAIL_OUTLET,
		enums.PAYMENT_METHOD_SENANGPAY,
	)
	if !paymentMethodValid {
		err = fmt.Errorf("invalid payment method please, %v", paymentMethodList)
		return err
	}

	bankNameValid, bankNameList := helpers.IsInList(payload.BankName,
		enums.CHANNEL_CIMB,
		enums.CHANNEL_SENANGPAY,
		enums.CHANNEL_CIMB,
		enums.CHANNEL_SHOPEE,
		enums.CHANNEL_BCA,
		enums.CHANNEL_GOPAY,
		enums.CHANNEL_MANDIRI,
		enums.CHANNEL_PERMATA,
		enums.CHANNEL_BNI,
		enums.CHANNEL_BRI,
		enums.CHANNEL_DANAMON,
		enums.CHANNEL_MAYBANK,
		enums.CHANNEL_K_WALLET,
	)
	if !bankNameValid {
		err = fmt.Errorf("invalid bank name please, %v", bankNameList)
		return err
	}

	return nil
}

func ValidationArrayChannel(payload []web.DetailChannelResponse) error {
	seenCh := make(map[string]bool)

	for idx1, ch := range payload {
		err := ValidationChannel(&ch)
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

	return nil
}
