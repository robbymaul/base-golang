package models

import "time"

type AdminSessions struct {
	Id           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	AdminUserId  int64     `gorm:"column:admin_user_id;foreignKey"`
	SessionToken string    `gorm:"column:session_token"`
	RefreshToken string    `gorm:"column:refresh_token"`
	IdAddress    string    `gorm:"column:id_address"`
	UserAgent    string    `gorm:"column:user_agent"`
	ExpiresAt    time.Time `gorm:"column:expires_at"`
	IsActive     bool      `gorm:"column:is_active"`
	LastUsedAt   time.Time `gorm:"column:last_used_at"`
	BaseField
}

func (adminSessions *AdminSessions) TableName() string {
	return "admin_sessions"
}
