package models

import (
	"paymentserviceklink/app/enums"
	"time"

	"github.com/shopspring/decimal"
)

type TopupTransaction struct {
	ID          int64                `gorm:"column:id;primaryKey;notNull;autoIncrement"`
	KWalletID   int64                `gorm:"column:k_wallet_id;foreignKey"`
	MemberId    string               `gorm:"column:member_id"`
	ChannelID   int64                `gorm:"column:channel_id;foreignKey"`
	Aggregator  enums.AggregatorName `gorm:"column:aggregator"`
	Merchant    string               `gorm:"column:merchant"`
	Amount      decimal.Decimal      `gorm:"column:amount"`
	FeeAdmin    decimal.Decimal      `gorm:"column:fee_admin"`
	Currency    enums.Currency       `gorm:"column:currency"`
	Symbol      enums.SymbolCurrency `gorm:"column:symbol"`
	ReferenceID string               `gorm:"column:reference_id"`
	Status      enums.PaymentStatus  `gorm:"column:status"`
	CompletedAt time.Time            `gorm:"column:completed_at"`
	Description string               `gorm:"column:description"`
	KWallet     *KWallet             `gorm:"foreignKey:KWalletID;references:ID"`
	BaseField
}

func (*TopupTransaction) TableName() string {
	return "topup_transaction"
}

func AllowedFilterColumnTopupTransaction() map[string]FilterColumn {
	return map[string]FilterColumn{
		"no_rekening": FilterColumn{
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "k_wallet",
		},
		"member_id": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "k_wallet",
		},
		"status": {
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "topup_transaction",
		},
		"currency": {
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "topup_transaction",
		},
	}
}
