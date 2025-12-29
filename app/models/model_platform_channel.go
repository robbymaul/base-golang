package models

type PlatformChannel struct {
	Id         int64 `gorm:"column:id;primaryKey;autoIncrement"`
	PlatformId int64 `gorm:"column:platform_id;foreignKey"`
	ChannelId  int64 `gorm:"column:channel_id;foreignKey"`
}

func (*PlatformChannel) TableName() string {
	return "platform_channel"
}
