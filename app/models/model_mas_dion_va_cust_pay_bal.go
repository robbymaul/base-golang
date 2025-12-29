package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type MasDionVaCustPayBal struct {
	Id          int             `gorm:"column:id;primaryKey"`
	Trcd        string          `gorm:"column:trcd"`
	Trdt        time.Time       `gorm:"column:trdt"`
	Novac       string          `gorm:"column:novac"`
	Dfno        string          `gorm:"column:dfno"`
	Fullnm      string          `gorm:"column:fullnm"`
	Type        string          `gorm:"column:type"`
	Refno       string          `gorm:"column:refno"`
	Amount      decimal.Decimal `gorm:"column:amount"`
	Status      string          `gorm:"column:status"`
	Custtype    string          `gorm:"column:custtype"`
	Description string          `gorm:"column:description"`
	Remarks     string          `gorm:"column:remarks"`
	Migration   bool            `gorm:"column:migration"`
}

func (*MasDionVaCustPayBal) TableName() string {
	return "mas_dion_va_cust_pay_bal"
}
