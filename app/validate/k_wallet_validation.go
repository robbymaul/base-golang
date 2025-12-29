package validate

import (
	"fmt"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/web"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidationGetKWallet(payload *web.GetKWalletRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.MemberId, validation.Required),
	)

	return err
}

func ValidationGetListKWalletTransaction(payload *web.GetListKWalletTransactionRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.MemberId, validation.Required),
		validation.Field(&payload.NoRekening, validation.Required),
	)

	return err
}

func ValidationGetVirtualAccountKWallet(payload *web.GetVirtualAccountKWalletRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.MemberId, validation.Required),
		validation.Field(&payload.NoRekening, validation.Required),
	)
	return err
}

func ValidationCreateKWallet(payload *web.CreateKWalletRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.MemberID, validation.Required),
		validation.Field(&payload.CustomerPhone, validation.Required),
		validation.Field(&payload.CustomerName, validation.Required),
		validation.Field(&payload.CustomerEmail, validation.Required),
		validation.Field(&payload.Currency, validation.Required),
	)

	validCurrency, currencyList := helpers.IsInList(payload.Currency, []enums.Currency{enums.CURRENCY_IDR, enums.CURRENCY_MYR}...)
	if !validCurrency {
		return fmt.Errorf("currency is not valid, example currency valid %v", currencyList)
	}

	return err
}

func ValidationCreateTopupKWallet(payload *web.CreateTopupKWalletRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.MemberId, validation.Required),
		validation.Field(&payload.NoRekening, validation.Required),
		//validation.Field(&payload.Amount, validation.Required),
		validation.Field(&payload.Channel, validation.Required),
	)

	var validAmount int64 = 10000

	//if payload.Amount < validAmount {
	//	return fmt.Errorf("amount is not valid, amount must be greater than %v", validAmount)
	//}

	err = validation.ValidateStruct(payload.Channel,
		validation.Field(&payload.Channel.PaymentMethod, validation.Required),
		validation.Field(&payload.Channel.Amount, validation.Required),
	)

	if payload.Channel.Amount < validAmount {
		return fmt.Errorf("amount is not valid, amount must be greater than %v", validAmount)
	}

	return err
}
