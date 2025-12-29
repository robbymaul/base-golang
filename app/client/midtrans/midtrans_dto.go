package midtrans

import (
	"paymentserviceklink/app/enums"

	"github.com/shopspring/decimal"
)

// PaymentRequest represents the payment request to Midtrans
type PaymentRequest struct {
	OrderID       string              `json:"order_id"`
	Amount        decimal.Decimal     `json:"amount"`
	PaymentType   string              `json:"payment_type"` // e.g., "bank_transfer", "credit_card"
	Method        enums.PaymentMethod `json:"method"`
	Channel       enums.Channel       `json:"channel"`
	Description   string              `json:"description"`
	CustomerName  string              `json:"customer_name"`
	CustomerEmail string              `json:"customer_email"`
	CustomerPhone string              `json:"customer_phone"`
	CallbackURL   string              `json:"callback_url"`
}

// ChargeResponse represents the response from Midtrans charge API
type ChargeResponse struct {
	StatusCode        string     `json:"status_code"`
	StatusMessage     string     `json:"status_message"`
	TransactionID     string     `json:"transaction_id"`
	OrderID           string     `json:"order_id"`
	GrossAmount       string     `json:"gross_amount"`
	PaymentType       string     `json:"payment_type"`
	TransactionTime   string     `json:"transaction_time"`
	TransactionStatus string     `json:"transaction_status"`
	FraudStatus       string     `json:"fraud_status"`
	Actions           []Action   `json:"actions,omitempty"`
	VaNumbers         []VaNumber `json:"va_numbers,omitempty"`
	PaymentCode       string     `json:"payment_code,omitempty"`
	Store             string     `json:"store,omitempty"`
	MerchantID        string     `json:"merchant_id"`
	Currency          string     `json:"currency"`
}

// StatusResponse represents the response from Midtrans status API
type StatusResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
}

// Action represents available actions for the payment
type Action struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

// VaNumber represents virtual account number information
type VaNumber struct {
	Bank     string `json:"bank"`
	VaNumber string `json:"va_number"`
}

// Notification represents the payment notification from Midtrans
type Notification struct {
	TransactionTime        string `json:"transaction_time"`
	TransactionStatus      string `json:"transaction_status"`
	TransactionID          string `json:"transaction_id"`
	StatusMessage          string `json:"status_message"`
	StatusCode             string `json:"status_code"`
	SignatureKey           string `json:"signature_key"`
	SettlementTime         string `json:"settlement_time"`
	PaymentType            string `json:"payment_type"`
	OrderID                string `json:"order_id"`
	MerchantID             string `json:"merchant_id"`
	MaskedCard             string `json:"masked_card"`
	GrossAmount            string `json:"gross_amount"`
	FraudStatus            string `json:"fraud_status"`
	Eci                    string `json:"eci"`
	Currency               string `json:"currency"`
	ChannelResponseMessage string `json:"channel_response_message"`
	ChannelResponseCode    string `json:"channel_response_code"`
	CardType               string `json:"card_type"`
	Bank                   string `json:"bank"`
	ApprovalCode           string `json:"approval_code"`
}

type CheckStatusPaymentRequest struct {
	TransactionId string `json:"transactionId"`
	OrderId       string `json:"orderId"`
}

type CheckStatusPaymentResponse struct {
	StatusCode               string `json:"status_code"`
	StatusMessage            string `json:"status_message"`
	TransactionId            string `json:"transaction_id"`
	MaskedCard               string `json:"masked_card"`
	OrderId                  string `json:"order_id"`
	PaymentType              string `json:"payment_type"`
	TransactionTime          string `json:"transaction_time"`
	TransactionStatus        string `json:"transaction_status"`
	FraudStatus              string `json:"fraud_status"`
	ApprovalCode             string `json:"approval_code"`
	SignatureKey             string `json:"signature_key"`
	Bank                     string `json:"bank"`
	GrossAmount              string `json:"gross_amount"`
	ChannelResponseCode      string `json:"channel_response_code"`
	ChannelResponseMessage   string `json:"channel_response_message"`
	CardType                 string `json:"card_type"`
	PaymentOptionType        string `json:"payment_option_type"`
	ShopeepayReferenceNumber string `json:"shopeepay_reference_number"`
	ReferenceId              string `json:"reference_id"`
}
