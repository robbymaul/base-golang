package models

type PlatformConfiguration struct {
	Id              int64 `gorm:"column:id;primaryKey;autoIncrement"`
	ConfigurationId int64 `gorm:"column:configuration_id;foreignKey"`
	PlatformId      int64 `gorm:"column:platform_id;foreignKey"`
}

func (*PlatformConfiguration) TableName() string {
	return "platform_configuration"
}
