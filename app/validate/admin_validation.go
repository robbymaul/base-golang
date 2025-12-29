package validate

import (
	"paymentserviceklink/app/web"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidationAdminLoginRequest(payload *web.AdminLoginRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.UsernameOrEmail, validation.Required),
		validation.Field(&payload.Password, validation.Required),
	)

	return err
}

func ValidationCreateAdminRolesRequest(payload *web.CreateAdminRolesRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Code, validation.Required),
		validation.Field(&payload.Name, validation.Required),
	)

	return err
}

func ValidationUpdateRoleRequest(payload *web.AdminRoleResponse) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Id, validation.Required),
		validation.Field(&payload.Code, validation.Required),
		validation.Field(&payload.Name, validation.Required),
	)
	return err
}

func ValidationCreateAdminUserRequest(payload *web.CreateAdminUserRequest) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Username, validation.Required),
		validation.Field(&payload.Email, validation.Required, is.Email),
		validation.Field(&payload.Password, validation.Required),
		validation.Field(&payload.FullName, validation.Required),
		validation.Field(&payload.Phone, validation.Required),
		validation.Field(&payload.RoleId, validation.Required),
	)
	return err
}

func ValidationUpdateAdminUserRequest(payload *web.DetailAdminResponse) error {
	err := validation.ValidateStruct(payload,
		validation.Field(&payload.Id, validation.Required),
		//validation.Field(&payload.Username, validation.Required),
		validation.Field(&payload.Email, validation.Required),
		validation.Field(&payload.FullName, validation.Required),
		validation.Field(&payload.Phone, validation.Required),
		validation.Field(&payload.RoleId, validation.Required),
	)
	return err
}
