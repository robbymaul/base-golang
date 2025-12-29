package models

type Platforms struct {
	Id              int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Code            string `gorm:"column:code"`
	Name            string `gorm:"column:name"`
	Description     string `gorm:"column:description"`
	ApiKey          string `gorm:"column:api_key"`
	SecretKey       string `gorm:"column:secret_key"`
	IsActive        bool   `gorm:"column:is_active"`
	NotificationURL string `gorm:"column:notification_url"`
	BaseField
}

func (platforms *Platforms) TableName() string {
	return "platforms"
}

func AllowedFilterColumnPlatform() map[string]FilterColumn {
	return map[string]FilterColumn{
		"name": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "platforms",
		},
		"is_active": FilterColumn{
			Operator: []string{"eq", "ne"},
			Variant:  "boolean",
			Table:    "platforms",
		},
	}
}
