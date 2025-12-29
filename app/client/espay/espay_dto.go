package espay

import (
	"paymentserviceklink/app/enums"
	"time"

	"github.com/rs/zerolog/log"
)

// PaymentRequest represents the payment request to espay
type PaymentRequest struct {
	RQUUID     string    `json:"rq_uuid"`
	RQDateTime time.Time `json:"rq_datetime"`
	OrderID    string    `json:"order_id"`
	Amount     string    `json:"amount"`
	FeeAmount  string    `json:"fee_amount"`
	//PaymentType   string  `json:"payment_type"` // e.g., "bank_transfer", "credit_card"
	CCY enums.Currency `json:"ccy"` // e.g., "IDR", "USD" kode mata uang transaksi
	//CommCode      string         `json:"comm_code"` // e.g. SGWYESSISHOP // MERCHANT CODE
	Method        enums.PaymentMethod `json:"method"`
	CustomerID    string              `json:"customer_id"`
	CustomerPhone string              `json:"customer_phone"`
	CustomerName  string              `json:"customer_name"`
	CustomerEmail string              `json:"customer_email"`
	//Description   string              `json:"description"`
	BankCode    string          `json:"bank_code"`
	ProductCode string          `json:"product_code"`
	ProductName string          `json:"product_name"`
	VaExpired   enums.VaExpired `json:"va_expired"`
	//MerchantId    string `json:"merchant_id"`
	//Signature string `json:"signature"`
	ReturnUrl string `json:"return_url"`
}

type CheckStatusPaymentRequest struct {
	PartnerServiceId string `json:"partnerServiceId"`
	CustomerNo       string `json:"customerNo"`
	VirtualAccountNo string `json:"virtualAccountNo"`
	InquiryRequestId string `json:"inquiryRequestId"`
	PaymentRequestId string `json:"paymentRequestId"`
	AdditionalInfo   struct {
		BillNo         string `json:"billNo"`
		IsPaymentNotif string `json:"isPaymentNotif"`
		PaymentRef     string `json:"paymentRef"`
		TerminalId     string `json:"terminalId"`
		TraceNumber    string `json:"traceNumber"`
	} `json:"additionalInfo"`
}

type CreditCardRequest struct {
	PartnerReferenceNo string             `json:"partnerReferenceNo"`
	MerchantId         string             `json:"merchantId"`
	SubMerchantId      string             `json:"subMerchantId"`
	Amount             CCAmount           `json:"amount"`
	UrlParam           CCUrlParam         `json:"urlParam"`
	ValidUpTo          time.Time          `json:"validUpTo"`
	PointOfInitiation  string             `json:"pointOfInitiation"`
	PayOptionDetails   CCPayOptionDetails `json:"payOptionDetails"`
	AdditionalInfo     CCAdditionalInfo   `json:"additionalInfo"`
}

type CCAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type CCUrlParam struct {
	Url        string `json:"url"`
	Type       string `json:"type"`
	IsDeeplink string `json:"isDeeplink"`
}

type CCPayOptionDetails struct {
	PayMethod   string        `json:"payMethod"`
	PayOption   string        `json:"payOption"`
	TransAmount CCTransAmount `json:"transAmount"`
	CCFeeAmount CCFeeAmount   `json:"feeAmount"`
}

type CCTransAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type CCFeeAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type CCAdditionalInfo struct {
	PayType       string `json:"payType"`
	UserId        string `json:"userId"`
	UserName      string `json:"userName"`
	UserEmail     string `json:"userEmail"`
	UserPhone     string `json:"userPhone"`
	BuyerId       string `json:"buyerId"`
	ProductCode   string `json:"productCode"`
	BalanceType   string `json:"balanceType"`
	BankCardToken string `json:"bankCardToken"`
}

type CCHeader map[string]string

func (c *CreditCardRequest) Headers(timestamp time.Time, signature string, externalId string, partnerId string) CCHeader {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"X-TIMESTAMP":   timestamp.Format(time.RFC3339),
		"X-SIGNATURE":   signature,
		"X-EXTERNAL-ID": externalId,
		"X-PARTNER-ID":  partnerId,
		"CHANNEL-ID":    "ESPAY",
	}
	log.Debug().Interface("headers", headers).Msg("headers credit card request")
	return headers
}

type QRISRequest struct {
	PartnerReferenceNo string `json:"partnerReferenceNo"`
	MerchantId         string `json:"merchantId"`
	Amount             struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	AdditionalInfo struct {
		ProductCode string `json:"productCode"`
	} `json:"additionalInfo"`
	ValidityPeriod time.Time `json:"validityPeriod"`
}

func (q *QRISRequest) Headers(timestamp time.Time, signature string, externalId string, partnerId string) map[string]string {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"X-TIMESTAMP":   timestamp.Format(time.RFC3339),
		"X-SIGNATURE":   signature,
		"X-EXTERNAL-ID": externalId,
		"X-PARTNER-ID":  partnerId,
		"CHANNEL-ID":    "ESPAY",
	}
	log.Debug().Interface("headers", headers).Msg("headers qris request")
	return headers
}

type PaymentLinkRequest struct {
}
