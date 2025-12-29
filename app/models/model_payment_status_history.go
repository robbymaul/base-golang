package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"paymentserviceklink/app/enums"
)

type PaymentStatusHistory struct {
	Id        int64               `gorm:"column:id;primaryKey;autoIncrement"`
	PaymentId int64               `gorm:"column:payment_id;foreignKey"`
	Status    enums.PaymentStatus `gorm:"column:status"`
	Notes     string              `gorm:"column:notes"`
	CreatedBy CreatedBy           `gorm:"column:created_by"`
	BaseField
}

type CreatedBy struct {
	ID       string
	Name     string
	Role     string
	Platform string
}

func (paymentStatusHistory *PaymentStatusHistory) TableName() string {
	return "payment_status_history"
}

// Implement driver.Valuer
func (p *CreatedBy) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Implement sql.Scanner
func (p *CreatedBy) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, p)
}
