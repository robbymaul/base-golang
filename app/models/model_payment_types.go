package models

import "paymentserviceklink/app/enums"

type PaymentTypes struct {
	Id          int64                 `gorm:"column:id;primaryKey;autoIncrement"`
	Code        enums.CodePaymentType `gorm:"column:code"`
	Name        string                `gorm:"column:name"`
	Description string                `gorm:"column:description"`
	IsActive    string                `gorm:"column:is_active"`
	BaseField
}

func (paymentTypes *PaymentTypes) TableName() string {
	return "payment_types"
}
