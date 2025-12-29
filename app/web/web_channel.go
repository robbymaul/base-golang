package web

import (
	"paymentserviceklink/app/enums"
	"time"
)

type CreateChannelRequest struct {
	//Aggregator      *AggregatorResponse         `json:"aggregator"`
	//Code            enums.CodePaymentMethod `json:"code"`
	Name            string              `json:"name"`
	PaymentMethod   enums.PaymentMethod `json:"paymentMethod"`
	TransactionType string              `json:"transactionType"`
	Currency        enums.Currency      `json:"currency"`
	FeeType         enums.FeeType       `json:"feeType"`
	FeeFixed        int64               `json:"feeFixed"`
	FeePercentage   float32             `json:"feePercentage"`
	//IsEspay         bool                    `json:"isEspay"`
	ProductName string          `json:"productName"`
	ProductCode string          `json:"productCode"`
	BankName    enums.Channel   `json:"bankName"`
	BankCode    string          `json:"bankCode"`
	Instruction string          `json:"instruction"`
	Image       []*ImageRequest `json:"images"`
}

type GetListChannelRequest struct {
	MemberId string         `json:"memberId"`
	Currency enums.Currency `json:"currency"`
}

type DetailChannelResponse struct {
	Id            int64               `json:"id,omitempty"`
	Code          string              `json:"code,omitempty"`
	Name          string              `json:"name,omitempty"`
	PaymentMethod enums.PaymentMethod `json:"paymentMethod,omitempty"`
	//TransactionType string              `json:"transactionType"`
	//Provider        enums.ProviderPaymentMethod `json:"provider"`
	Currency      enums.Currency `json:"currency,omitempty"`
	FeeType       enums.FeeType  `json:"feeType,omitempty"`
	FeeFixed      int64          `json:"feeFixed"`
	FeePercentage float32        `json:"feePercentage"`
	IsActive      bool           `json:"isActive"`
	// penanda awal k-wallet
	MemberID   string               `json:"memberId,omitempty"`
	FullName   string               `json:"fullName,omitempty"`
	NoRekening string               `json:"noRekening,omitempty"`
	GenVa      string               `json:"genVa,omitempty"`
	Balance    int64                `json:"balance"`
	Symbol     enums.SymbolCurrency `json:"symbol,omitempty"`
	Status     enums.KWalletStatus  `json:"status,omitempty"`
	//Currency       enums.Currency           `json:"currency"`
	//IsActive       bool                     `json:"isActive"`
	// penanda akhir k-wallet
	ProductName string           `json:"productName,omitempty"`
	ProductCode string           `json:"productCode,omitempty"`
	BankName    enums.Channel    `json:"bankName,omitempty"`
	BankCode    string           `json:"bankCode,omitempty"`
	Instruction string           `json:"instruction,omitempty"`
	CreatedAt   *time.Time       `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time       `json:"updatedAt,omitempty"`
	Images      []*ImageResponse `json:"images"`
}

type ChannelResponse struct {
	Id            int64               `json:"id"`
	Name          string              `json:"name"`
	MethodGroup   enums.PaymentMethod `json:"methodGroup"`
	PaymentMethod enums.PaymentMethod `json:"paymentMethod"`
	Currency      enums.Currency      `json:"currency"`
	//Symbol        enums.SymbolCurrency `json:"symbol"`
	FeeType       enums.FeeType    `json:"feeType"`
	FeeFixed      int64            `json:"feeFixed"`
	FeePercentage float32          `json:"feePercentage"`
	Bank          enums.Channel    `json:"bank"`
	Instruction   string           `json:"instruction"`
	Images        []*ImageResponse `json:"images"`
}
