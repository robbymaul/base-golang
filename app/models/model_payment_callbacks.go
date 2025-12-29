package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/datatypes"
	"paymentserviceklink/app/enums"
)

type PaymentCallbacks struct {
	Id           int64               `gorm:"column:id;primaryKey;autoIncrement"`
	PaymentId    int64               `gorm:"column:payment_id;foreignKey"`
	GatewayName  string              `gorm:"column:gateway_name"`
	CallbackData datatypes.JSON      `gorm:"column:callback_data"`
	ResponseData datatypes.JSON      `gorm:"column:response_data"`
	Status       enums.PaymentStatus `gorm:"column:status"`
	BaseField
}

func (paymentCallbacks *PaymentCallbacks) TableName() string {
	return "payment_callbacks"
}

type CallbackData struct {
	Name          string
	Email         string
	Phone         string
	AmountPaid    string
	TxnStatus     enums.TxnStatusSenangpay
	TxnMessage    string
	OrderId       string
	TransactionId string
	HashedValue   string
}

// Implement driver.Valuer
func (p *CallbackData) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Implement sql.Scanner
func (p *CallbackData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, p)
}
