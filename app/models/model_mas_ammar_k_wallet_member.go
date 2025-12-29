package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type MasAmmarKWalletMember struct {
	RecId    int       `gorm:"column:rec_id;primaryKey"`
	IdMember string    `gorm:"column:id_member"`
	BankCode string    `gorm:"column:bank_code"`
	NumberVa string    `gorm:"column:number_va"`
	Nama     string    `gorm:"column:nama"`
	Hp       string    `gorm:"column:hp"`
	Email    string    `gorm:"column:email"`
	DateAdd  time.Time `gorm:"column:date_add"`
	DateUpd  time.Time `gorm:"column:date_upd"`
	Status   int       `gorm:"status"`
	UserId   int       `gorm:"user_id"`
}

func (*MasAmmarKWalletMember) TableName() string {
	return "mas_ammar_k_wallet_member"
}

type MasAmmarKWalletMemberSaldo struct {
	IdMember   string          `gorm:"column:id_member;primaryKey"`
	LastSaldo  decimal.Decimal `gorm:"column:last_saldo"`
	LastUpdate time.Time       `gorm:"column:last_update"`
	UserId     string          `gorm:"column:user_id"`
	CommCode   string          `gorm:"column:comm_code"`
}

func (*MasAmmarKWalletMemberSaldo) TableName() string {
	return "mas_ammar_k_wallet_member_saldo"
}

type MasAmmarKWalletGenVa struct {
	RecId    int       `gorm:"column:rec_id;primaryKey"`
	MemberId string    `gorm:"column:member_id"`
	UserId   string    `gorm:"user_id"`
	DatedAdd time.Time `gorm:"column:date_add"`
	DateUpd  time.Time `gorm:"column:date_upd"`
	Status   int       `gorm:"status"`
}

func (*MasAmmarKWalletGenVa) TableName() string {
	return "mas_ammar_k_wallet_gen_va"
}

type MasAmmarKWalletGenVaDetail struct {
	RecId    int    `gorm:"column:rec_id;primaryKey"`
	IdGenVa  int    `gorm:"column:id_gen_va"`
	BankCode string `gorm:"column:bank_code"`
	VaNumber string `gorm:"column:va_number"`
}

func (*MasAmmarKWalletGenVaDetail) TableName() string {
	return "mas_ammar_k_wallet_gen_va_detail"
}
