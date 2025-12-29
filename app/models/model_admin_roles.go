package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"paymentserviceklink/app/enums"
)

type AdminRoles struct {
	Id          int64                `gorm:"column:id;primaryKey;autoIncrement"`
	Code        enums.CodeAdminRole  `gorm:"column:code"`
	Name        enums.NameAdminRoles `gorm:"column:name"`
	Description string               `gorm:"column:description"`
	Permissions *Permissions         `gorm:"column:permissions"`
	IsActive    bool                 `gorm:"column:is_active"`
	BaseField
}

type Permissions struct {
}

func (adminRoles *AdminRoles) TableName() string {
	return "admin_roles"
}

// Implement driver.Valuer
func (p *Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Implement sql.Scanner
func (p *Permissions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, p)
}

func AllowedFilterColumnAdminRole() map[string]FilterColumn {
	return map[string]FilterColumn{
		"name": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "admin_roles",
		},
		"is_active": {
			Operator: []string{"eq"},
			Variant:  "boolean",
			Table:    "admin_roles",
		},
	}
}
