package models

import "paymentserviceklink/app/enums"

type ChannelImage struct {
	ID            int64                `gorm:"column:id;primaryKey;autoIncrement"`
	ChannelID     int64                `gorm:"column:channel_id;foreignKey"`
	FileName      string               `gorm:"column:file_name"`
	SizeType      enums.ImageSizeType  `gorm:"column:size_type"`
	GeometricType enums.ImageGeometric `gorm:"column:geometric_type"`
	BaseField
}

func (*ChannelImage) TableName() string {
	return "channel_images"
}
