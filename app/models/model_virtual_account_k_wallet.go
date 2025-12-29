package models

import "paymentserviceklink/app/enums"

type VirtualAccountKWallet struct {
	ID             int64         `gorm:"column:id;primaryKey;notNull;autoIncrement"`
	KWalletId      int64         `gorm:"column:k_wallet_id;foreignKey"`
	VirtualAccount string        `gorm:"column:virtual_account"`
	Bank           enums.Channel `gorm:"column:bank"`
	BankCode       string        `gorm:"bank_code"`
	BaseField
}

func (*VirtualAccountKWallet) TableName() string {
	return "virtual_account_k_wallet"
}

func AllowedFilterColumnVirtualAccountKWallet() map[string]FilterColumn {
	return map[string]FilterColumn{
		"bank": {
			Operator: []string{"eq", "like"},
			Variant:  "string",
			Table:    "virtual_account_k_wallet",
		},
		"member_id": {
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "k_wallet",
		},
		"no_rekening": {
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "k_wallet",
		},
	}
}
