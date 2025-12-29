package models

type PaymentItems struct {
	Id              int64   `gorm:"column:id;primaryKey;autoIncrement"`
	PaymentId       int64   `gorm:"column:payment_id;foreignKey"`
	ItemId          string  `gorm:"column:item_id"`
	ItemName        string  `gorm:"column:item_name"`
	ItemDescription string  `gorm:"column:item_description"`
	Quantity        int64   `gorm:"column:quantity"`
	UnitPrice       float32 `gorm:"column:unit_price"`
	TotalPrice      float32 `gorm:"column:total_price"`
	BasePrice       float32 `gorm:"column:base_price"`
}

func (paymentItems *PaymentItems) TableName() string {
	return "payment_items"
}
