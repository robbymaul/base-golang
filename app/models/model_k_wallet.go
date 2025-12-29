package models

import (
	"paymentserviceklink/app/enums"

	"github.com/shopspring/decimal"
)

type KWallet struct {
	ID             int64                    `gorm:"column:id;primaryKey;notNull;autoIncrement"`
	MemberID       string                   `gorm:"column:member_id;unique"`
	FullName       string                   `gorm:"column:full_name"`
	NoRekening     string                   `gorm:"column:no_rekening"`
	GenVA          string                   `gorm:"column:gen_va"`
	Balance        decimal.Decimal          `gorm:"column:balance"`
	Currency       enums.Currency           `gorm:"column:currency"`
	Symbol         enums.SymbolCurrency     `gorm:"column:symbol"`
	Status         enums.KWalletStatus      `gorm:"column:status"`
	IsActive       bool                     `gorm:"column:is_active"`
	VirtualAccount []*VirtualAccountKWallet `gorm:"foreignKey:KWalletId;references:ID"`
	BaseField
}

func (*KWallet) TableName() string {
	return "k_wallet"
}

//func (k *KWallet) TopupKWallet(amount decimal.Decimal) {
//	k.Balance = k.Balance.Add(amount)
//}

//func (k *KWallet) PaymentKWallet(amount decimal.Decimal) {
//	log.Debug().Interface("balance", k.Balance).Interface("amount", amount).Msg("before sub() payment k-wallet")
//	k.Balance = k.Balance.Sub(amount)
//	log.Debug().Interface("balance", k.Balance).Interface("amount", amount).Msg("after sub() payment k-wallet")
//}

func (k *KWallet) SubBalance(add decimal.Decimal) {
	k.Balance = k.Balance.Sub(add)
}

func (k *KWallet) AddBalance(add decimal.Decimal) {
	k.Balance = k.Balance.Add(add)
}

func AllowedFilterColumnKWallet() map[string]FilterColumn {
	return map[string]FilterColumn{
		"member_id": {
			Operator: []string{"eq", "like"},
			Variant:  "string",
			Table:    "k_wallet",
		},
		"full_name": {
			Operator: []string{"eq", "like"},
			Variant:  "string",
			Table:    "k_wallet",
		},
		"no_rekening": {
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "k_wallet",
		},
		"gen_va": {
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "k_wallet",
		},
		"currency": {
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "k_wallet",
		},
		"is_active": {
			Operator: []string{"eq"},
			Variant:  "boolean",
			Table:    "k_wallet",
		},
	}
}
