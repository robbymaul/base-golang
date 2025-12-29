package web

import (
	"paymentserviceklink/app/enums"
	"time"
)

type CreatePaymentRequest struct {
	ApiKey        string     `json:"apiKey"`
	SecretKey     string     `json:"secretKey"`
	Payment       []*Payment `json:"payments"`
	CustomerId    string     `json:"customerId"`
	CustomerName  string     `json:"customerName"`
	CustomerEmail string     `json:"customerEmail"`
	CustomerPhone string     `json:"customerPhone"`
	ReferenceId   string     `json:"referenceId"`
	ReferenceType string     `json:"referenceType"`
	ReturnUrl     string     `json:"returnUrl"`
}

type Payment struct {
	OrderId string            `json:"orderId"`
	Channel []*ChannelPayment `json:"channel"`
	Product []*PaymentProduct `json:"products"`
	Amount  int64             `json:"amount"`
}

type PaymentProduct struct {
}

type ChannelPayment struct {
	Id            int64                       `json:"id"`
	Code          string                      `json:"code"`
	Name          string                      `json:"name"`
	PaymentMethod enums.PaymentMethod         `json:"paymentMethod"`
	Provider      enums.ProviderPaymentMethod `json:"provider"`
	Currency      enums.Currency              `json:"currency"`
	FeeType       enums.FeeType               `json:"feeType"`
	FeeAmount     int64                       `json:"feeAmount"`
	Amount        int64                       `json:"amount"`
	MemberID      string                      `json:"memberId,omitempty"`
	FullName      string                      `json:"fullName,omitempty"`
	NoRekening    string                      `json:"noRekening,omitempty"`
	GenVa         string                      `json:"genVa,omitempty"`
	Balance       int64                       `json:"balance,omitempty"`
	Symbol        enums.SymbolCurrency        `json:"symbol,omitempty"`
	Status        enums.KWalletStatus         `json:"status,omitempty"`
	ProductName   string                      `json:"productName"`
	ProductCode   string                      `json:"productCode"`
	BankName      enums.Channel               `json:"bankName"`
}

type DetailPaymentResponse struct {
	Id                   int64                   `json:"id,omitempty"`
	TransactionId        string                  `json:"transactionId"`
	OrderId              string                  `json:"orderId"`
	PlatformId           int64                   `json:"platformId"`
	PaymentMethodId      int64                   `json:"paymentMethodId"`
	Amount               int64                   `json:"amount"`
	FeeAmount            int64                   `json:"feeAmount"`
	TotalAmount          int64                   `json:"totalAmount"`
	Currency             enums.Currency          `json:"currency"`
	Status               enums.PaymentStatus     `json:"status"`
	CustomerId           string                  `json:"customerId"`
	CustomerName         string                  `json:"customerName"`
	CustomerEmail        string                  `json:"customerEmail"`
	CustomerPhone        string                  `json:"customerPhone"`
	ReferenceId          string                  `json:"referenceId"`
	ReferenceType        string                  `json:"referenceType"`
	GatewayTransactionId string                  `json:"gatewayTransactionId"`
	GatewayReference     string                  `json:"gatewayReference"`
	GatewayResponse      any                     `json:"gatewayResponse"`
	CallbackUrl          string                  `json:"callbackUrl"`
	ReturnUrl            string                  `json:"returnUrl"`
	ExpiredAt            *time.Time              `json:"expiredAt"`
	PaidAt               *time.Time              `json:"paidAt"`
	CreatedAt            *time.Time              `json:"createdAt"`
	UpdatedAt            *time.Time              `json:"updatedAt"`
	Platform             *DetailPlatformResponse `json:"platform"`
	Channel              *DetailChannelResponse  `json:"channel"`
	Aggregator           *AggregatorResponse     `json:"aggregator"`
}

type GatewayResponse struct {
	RedirectUrl string `json:"redirectUrl"`
}

type GetDetailPaymentRequest struct {
	ApiKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
	OrderId   string `json:"orderId"`
}

type CheckStatusPaymentRequest struct {
	ApiKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
	OrderId   string `json:"orderId"`
}

type PaymentResponse struct {
	Id            int64               `json:"id"`
	TransactionId string              `json:"transactionId"`
	OrderId       string              `json:"orderId"`
	Status        enums.PaymentStatus `json:"status"`
	Amount        int64               `json:"amount"`
	FeeAdmin      int64               `json:"feeAdmin"`
	TotalAmount   int64               `json:"totalAmount"`
	Currency      enums.Currency      `json:"currency"`
	PaymentMethod enums.PaymentMethod `json:"paymentMethod"`
	PaymentType   enums.PaymentType   `json:"responseType"`
	PaymentDetail PaymentDetail       `json:"paymentDetail"`
	Customer      Customer            `json:"customer"`
	CreatedAt     *time.Time          `json:"createdAt"`
	UpdatedAt     *time.Time          `json:"updatedAt"`
}

type PaymentDetail struct {
	Bank            enums.Channel `json:"bank"`
	Url             []*Actions    `json:"url"`
	VaNumber        string        `json:"vaNumber"`
	BillKey         string        `json:"billKey"`
	BIllCode        string        `json:"billCode"`
	TransactionTime string        `json:"transactionTime"`
	ExpireTime      string        `json:"expireTime"`
	Instruction     string        `json:"instruction"`
}

type Actions struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Url    string `json:"url"`
}

type Customer struct {
	MemberId string `json:"memberId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type CheckStatusPaymentResponse struct {
	Id            int64               `json:"id"`
	TransactionId string              `json:"transactionId"`
	OrderId       string              `json:"orderId"`
	Status        enums.PaymentStatus `json:"status"`
	Amount        int64               `json:"amount"`
	Currency      enums.Currency      `json:"currency"`
}

type PaymentCallback struct {
	TransactionId string              `json:"transactionId"`
	OrderId       string              `json:"orderId"`
	Status        enums.PaymentStatus `json:"status"`
	Amount        int64               `json:"amount"`
	Currency      enums.Currency      `json:"currency"`
}
