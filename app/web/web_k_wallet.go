package web

import (
	"paymentserviceklink/app/enums"

	"github.com/shopspring/decimal"
)

type VirtualAccountKWallet struct {
	ID             int64         `json:"id,omitempty"`
	Bank           enums.Channel `json:"bank,omitempty"`
	BankCode       string        `json:"bankCode,omitempty"`
	VirtualAccount string        `json:"virtualAccount,omitempty"`
}

type KWalletResponse struct {
	MemberID       string                   `json:"memberId,omitempty"`
	FullName       string                   `json:"fullName,omitempty"`
	NoRekening     string                   `json:"noRekening,omitempty"`
	GenVa          string                   `json:"genVa,omitempty"`
	Balance        int64                    `json:"balance"`
	Currency       enums.Currency           `json:"currency,omitempty"`
	Symbol         enums.SymbolCurrency     `json:"symbol,omitempty"`
	Status         enums.KWalletStatus      `json:"status,omitempty"`
	IsActive       bool                     `json:"isActive"`
	CreatedAt      string                   `json:"createdAt,omitempty"`
	VirtualAccount []*VirtualAccountKWallet `json:"virtualAccount,omitempty"`
}

type KWalletTransaction struct {
	ID                       int64                  `json:"id"`
	KWalletTypeTransactionID int64                  `json:"kWalletTypeTransactionId"`
	PaymentID                string                 `json:"paymentId"`
	Title                    string                 `json:"title"`
	PaymentCode              string                 `json:"paymentCode"`
	TransactionCode          string                 `json:"transactionCode"`
	TransactionType          string                 `json:"transactionType"`
	Direction                enums.KWalletDirection `json:"direction"`
	CounterpartyName         string                 `json:"counterpartyName"`
	CounterpartyBank         enums.Channel          `json:"counterpartyBank"`
	PaymentChannel           string                 `json:"paymentChannel"`
	Description              string                 `json:"description"`
	Balance                  int64                  `json:"balance"`
	Debit                    int64                  `json:"debit"`
	Credit                   int64                  `json:"credit"`
	Amount                   int64                  `json:"amount"`
	Currency                 enums.Currency         `json:"currency"`
	Symbol                   enums.SymbolCurrency   `json:"symbol"`
	Date                     string                 `json:"date"`
	Time                     string                 `json:"time"`
	DateTime                 string                 `json:"dateTime"`
}

type TopupTransaction struct {
	ID          int64                `json:"id"`
	KWalletID   int64                `json:"kWalletId"`
	MemberID    string               `json:"memberId"`
	ChannelId   int64                `json:"channelId"`
	Aggregator  enums.AggregatorName `json:"aggregator"`
	Merchant    string               `json:"merchant"`
	Amount      decimal.Decimal      `json:"amount"`
	FeeAdmin    decimal.Decimal      `json:"feeAdmin"`
	Currency    enums.Currency       `json:"currency"`
	Symbol      enums.SymbolCurrency `json:"symbol"`
	ReferenceID string               `json:"referenceId"`
	Status      enums.PaymentStatus  `json:"status"`
	CompletedAt string               `json:"completedAt"`
	Description string               `json:"description"`
	KWallet     *KWalletResponse     `json:"kWallet"`
}

type CreateKWalletRequest struct {
	MemberID      string         `json:"memberId"`
	CustomerPhone string         `json:"customerPhone"`
	CustomerName  string         `json:"customerName"`
	CustomerEmail string         `json:"customerEmail"`
	Currency      enums.Currency `json:"currency"`
}

type GetKWalletRequest struct {
	MemberId string `json:"memberId"`
}

type GetListKWalletTransactionRequest struct {
	MemberId   string `json:"memberId"`
	NoRekening string `json:"noRekening"`
	FromDate   string `json:"fromDate"`
	ToDate     string `json:"toDate"`
	Month      int64  `json:"month"`
}

type GetVirtualAccountKWalletRequest struct {
	MemberId   string `json:"memberId"`
	NoRekening string `json:"noRekening"`
}

type CreateTopupKWalletRequest struct {
	MemberId   string `json:"memberId"`
	NoRekening string `json:"noRekening"`
	//Amount     int64           `json:"amount"`
	Channel *ChannelPayment `json:"channel"`
}
