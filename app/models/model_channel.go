package models

import (
	"paymentserviceklink/app/enums"
)

type Channel struct {
	Id int64 `gorm:"column:id;primaryKey;autoIncrement"`
	//AggregatorId    int64                       `gorm:"column:aggregator_id;foreignKey"`
	Code            string                      `gorm:"column:code"`
	Name            string                      `gorm:"column:name"`
	PaymentMethod   enums.PaymentMethod         `gorm:"column:payment_method"`
	TransactionType string                      `gorm:"column:transaction_type"`
	Provider        enums.ProviderPaymentMethod `gorm:"column:provider"`
	Currency        enums.Currency              `gorm:"column:currency"`
	FeeType         enums.FeeType               `gorm:"column:fee_type"`
	FeeAmount       int64                       `gorm:"column:fee_amount"`
	FeePercentage   float32                     `gorm:"column:fee_percentage"`
	IsActive        bool                        `gorm:"column:is_active"`
	//IsEspay         bool                        `gorm:"column:is_espay"`
	ProductName  string          `gorm:"column:product_name"`
	ProductCode  string          `gorm:"column:product_code"`
	Instruction  string          `gorm:"column:instruction"`
	BankName     enums.Channel   `gorm:"column:bank_name"`
	BankCode     string          `gorm:"column:bank_code"`
	ChannelImage []*ChannelImage `gorm:"foreignKey:ChannelID;references:id"`
	BaseField
}

func (*Channel) TableName() string {
	return "channels"
}

func AllowedFilterColumnChannel() map[string]FilterColumn {
	return map[string]FilterColumn{
		"name": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "channels",
		},
		"platform_id": FilterColumn{
			Operator: []string{"eq"},
			Variant:  "number",
			Table:    "platform_channel",
		},
		"payment_method": FilterColumn{
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "channels",
		},
		"currency": FilterColumn{
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "channels",
		},
		"is_active": FilterColumn{
			Operator: []string{"eq", "ne"},
			Variant:  "boolean",
			Table:    "channels",
		},
		"product_name": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "channels",
		},
		"product_code": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "channels",
		},
		"bank_name": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "channels",
		},
		"bank_code": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "channels",
		},
	}
}
