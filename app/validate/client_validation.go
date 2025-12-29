package validate

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"paymentserviceklink/app/web"
)

func ValidationClientRequest(payload *web.ClientRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.ApiKey, validation.Required),
		validation.Field(&payload.SecretKey, validation.Required),
	)

	return err
}
