package validate

import (
	"fmt"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/web"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/zerolog/log"
)

func ValidationCreatePaymentRequest(payload *web.CreatePaymentRequest) error {
	swapChannelKWallet1st(&payload)
	log.Debug().Interface("payload", payload).Msg("data payload after swap channel")

	err := validation.ValidateStruct(payload,
		validation.Field(&payload.ApiKey, validation.Required),
		validation.Field(&payload.SecretKey, validation.Required),
		//validation.Field(&payload.TransactionId, validation.Required),
		//validation.Field(&payload.Channel, validation.Required),
		//validation.Field(&payload.PaymentType, validation.Required),
		//validation.Field(&payload.Amount, validation.Required),
		validation.Field(&payload.Payment, validation.Required),
		validation.Field(&payload.CustomerId, validation.Required),
		//validation.Field(&payload.CustomerEmail, validation.Required, is.Email),
		//validation.Field(&payload.CustomerPhone, validation.Required),
	)
	if err != nil {
		return err
	}

	for _, pay := range payload.Payment {
		totalPaymentChannel := int64(0)

		for _, ch := range pay.Channel {
			err = validation.ValidateStruct(ch,
				validation.Field(&ch.Id, validation.Required),
				validation.Field(&ch.Amount, validation.Required),
			)
			if err != nil {
				return err
			}

			if ch.PaymentMethod == enums.PAYMENT_METEHOD_K_WALLET {
				if ch.Status != enums.KWalletStatusActive {
					return fmt.Errorf("k-wallet not active")
				}

				//if len(pay.Channel) > 1 && ch.Amount < ch.Balance {
				//	return fmt.Errorf("cannot split payment")
				//}
				//
				//if len(pay.Channel) <= 1 && ch.Amount > ch.Balance {
				//	return fmt.Errorf("k-wallet balance not enough")
				//}

				if ch.NoRekening == "" {
					return fmt.Errorf("rekening number not found or blank, please call administrator")
				}
			}

			totalPaymentChannel += ch.Amount
		}

		if totalPaymentChannel != pay.Amount {
			return fmt.Errorf("total payment amount channel not match with amount")
		}
	}

	//paymentTypeValid, paymentTypeList := helpers.IsInList(payload.PaymentType, enums.CODE_PAYMENT_TYPE_SALES_ORDER, enums.CODE_PAYMENT_TYPE_TOPUP_TOKEN, enums.CODE_PAYMENT_TYPE_TOPUP_WALLET)
	//if !paymentTypeValid {
	//	err = fmt.Errorf("invalid payment type please call administrator, %v", paymentTypeList)
	//	return err
	//}

	return nil
}

func swapChannelKWallet1st(payload **web.CreatePaymentRequest) {
	for _, pay := range (*payload).Payment {

		if len(pay.Channel) > 1 {
			if pay.Channel[0].PaymentMethod != enums.PAYMENT_METEHOD_K_WALLET {
				temp := pay.Channel[0]
				pay.Channel[0] = pay.Channel[1]
				pay.Channel[1] = temp
			}
		}
	}
}

func CheckStatusPaymentRequest(payload *web.CheckStatusPaymentRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.ApiKey, validation.Required),
		validation.Field(&payload.SecretKey, validation.Required),
		validation.Field(&payload.OrderId, validation.Required),
	)

	return err
}

func ValidationEspayInquiryRequest(payload *web.EspayInquiryRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.InquiryRequestId, validation.Required),
		validation.Field(&payload.CustomerNo, validation.Required),
		validation.Field(&payload.PartnerServiceId, validation.Required),
		validation.Field(&payload.TrxDateInit, validation.Required),
		validation.Field(&payload.VirtualAccountNo, validation.Required),
	)

	return err
}

func ValidationEspayPaymentNotificationRequest(payload *web.EspayPaymentNotificationRequest) error {
	//err := validation.ValidateStruct(payload,
	//	//validation.Field(&payload.RQUUID, validation.Required),
	//	//validation.Field(&payload.RSDateTime, validation.Required),
	//	//validation.Field(&payload.Signature, validation.Required),
	//	//validation.Field(&payload.CommCode, validation.Required),
	//	//validation.Field(&payload.OrderId, validation.Required),
	//	//validation.Field(&payload.CCY, validation.Required),
	//	//validation.Field(&payload.Amount, validation.Required),
	//	//validation.Field(&payload.DebitFromBank, validation.Required),
	//	//validation.Field(&payload.CreditToBank, validation.Required),
	//	//validation.Field(&payload.ProductCode, validation.Required),
	//	//validation.Field(&payload.PaymentDatetime, validation.Required),
	//	//validation.Field(&payload.PaymentRef, validation.Required),
	//)

	return nil
}

func ValidationEspayTopupNotificationRequest(payload *web.EspayTopupNotificationRequest) error {

	return nil
}

func ValidationGetDetailPaymentRequest(payload *web.GetDetailPaymentRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.OrderId, validation.Required),
	)

	return err
}
