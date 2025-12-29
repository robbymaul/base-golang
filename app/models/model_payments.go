package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"paymentserviceklink/app/enums"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Payments struct {
	Id                   int64               `gorm:"column:id;primaryKey;autoIncrement"`
	TransactionId        string              `gorm:"column:transaction_id"`
	OrderId              string              `gorm:"column:order_id"`
	PlatformId           int64               `gorm:"column:platform_id;foreignKey"`
	PaymentMethodId      int64               `gorm:"column:payment_method_id;foreignKey"`
	AggregatorId         *int64              `gorm:"column:aggregator_id;foreignKey"`
	Amount               decimal.Decimal     `gorm:"column:amount"`
	FeeAmount            decimal.Decimal     `gorm:"column:fee_amount"`
	TotalAmount          decimal.Decimal     `gorm:"column:total_amount"`
	Currency             enums.Currency      `gorm:"column:currency"`
	Status               enums.PaymentStatus `gorm:"column:status"`
	CustomerId           string              `gorm:"column:customer_id"`
	CustomerName         string              `gorm:"column:customer_name"`
	CustomerEmail        string              `gorm:"column:customer_email"`
	CustomerPhone        string              `gorm:"column:customer_phone"`
	ReferenceId          string              `gorm:"column:reference_id"`
	ReferenceType        string              `gorm:"column:reference_type"`
	GatewayTransactionId string              `gorm:"column:gateway_transaction_id"`
	GatewayReference     string              `gorm:"column:gateway_reference"`
	GatewayResponse      datatypes.JSON      `gorm:"column:gateway_response"`
	CallbackUrl          string              `gorm:"column:callback_url"`
	ReturnUrl            string              `gorm:"column:return_url"`
	ExpiredAt            *time.Time          `gorm:"column:expired_at"`
	ExpiredTime          string              `gorm:"column:expired_time"`
	PaidAt               *time.Time          `gorm:"column:paid_at"`
	NotificationCallback *bool               `gorm:"column:notification_callback"`
	Platform             *Platforms          `gorm:"foreignKey:PlatformId;references:Id"`
	Channel              *Channel            `gorm:"foreignKey:PaymentMethodId;references:Id"`
	Aggregator           *Aggregator         `gorm:"foreignKey:AggregatorId;references:Id"`
	BaseField
}

type GatewayResponse struct {
	RedirectUrl string
}

func (p *Payments) TableName() string {
	return "payments"
}

// Implement driver.Valuer
func (p *GatewayResponse) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Implement sql.Scanner
func (p *GatewayResponse) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, p)
}

func AllowedFilterColumnPayment() map[string]FilterColumn {
	return map[string]FilterColumn{
		"order_id": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "payments",
		},
		"total_amount": FilterColumn{
			Operator: []string{"eq", "ne", "lt", "lte", "gt", "gte"},
			Variant:  "number",
			Table:    "payments",
		},
		"status": FilterColumn{
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "payments",
		},
		"customer_id": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "payments",
		},
		"customer_name": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "payments",
		},
		"customer_email": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "payments",
		},
		"customer_phone": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "payments",
		},
	}
}

func (p *Payments) SetFixPayment(totalAmount decimal.Decimal, feeAmount decimal.Decimal, expireTime string) {
	p.TotalAmount = totalAmount
	p.FeeAmount = feeAmount
	p.ExpiredTime = expireTime
	expiredAt, _ := time.Parse(time.DateTime, expireTime)
	p.ExpiredAt = &expiredAt
}
