package models

import "paymentserviceklink/app/enums"

type Aggregator struct {
	Id          int64                       `gorm:"column:id;primaryKey;autoIncrement"`
	Name        enums.AggregatorName        `gorm:"column:name"`
	Slug        enums.ProviderPaymentMethod `gorm:"column:slug"`
	Description string                      `gorm:"column:description"`
	IsActive    bool                        `gorm:"column:is_active"`
	Currency    enums.Currency              `gorm:"column:currency"`
	BaseField
}

func (*Aggregator) TableName() string {
	return "aggregators"
}

func AllowedFilterColumnAggregator() map[string]FilterColumn {
	return map[string]FilterColumn{
		"name": {
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "aggregators",
		},
		"is_active": {
			Operator: []string{"eq"},
			Variant:  "boolean",
			Table:    "aggregators",
		},
		"currency": {
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "aggregators",
		},
	}
}
