package validate

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"paymentserviceklink/pkg/pagination"
)

func PaginationValidation(payload *pagination.Pages) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Page, validation.Required),
		validation.Field(&payload.PerPage, validation.Required),
		validation.Field(&payload.Sort, validation.Required),
	)

	return err
}
