package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminUsers struct {
	Id                  int64       `gorm:"column:id;primaryKey;autoIncrement;"`
	UUID                uuid.UUID   `gorm:"column:uuid;index;"`
	Username            string      `gorm:"column:username"`
	Email               string      `gorm:"column:email"`
	Password            string      `gorm:"column:password_hash"`
	FullName            string      `gorm:"column:full_name"`
	Phone               string      `gorm:"column:phone"`
	AvatarUrl           string      `gorm:"column:avatar_url"`
	RoleId              int64       `gorm:"column:role_id"`
	IsActive            bool        `gorm:"column:is_active"`
	IsVerified          bool        `gorm:"column:is_verified"`
	LastLoginAt         *time.Time  `gorm:"column:last_login_at"`
	LastLoginIp         string      `gorm:"column:last_login_ip"`
	FailedLoginAttempts int64       `gorm:"column:failed_login_attempts"`
	LockedUntil         *time.Time  `gorm:"column:locked_until"`
	PasswordChangedAt   *time.Time  `gorm:"column:password_changed_at"`
	TwoFactorEnabled    bool        `gorm:"column:two_factor_enabled"`
	TwoFactorSecret     string      `gorm:"column:two_factor_secret"`
	AdminRole           *AdminRoles `gorm:"foreignKey:RoleId;references:Id"`
	BaseField
}

func (adminUsers *AdminUsers) TableName() string {
	return "admin_users"
}

func AllowedFilterColumnAdminUser() map[string]FilterColumn {
	return map[string]FilterColumn{
		"username": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "admin_users",
		},
		"email": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "admin_users",
		},
		"full_name": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "admin_users",
		},
		"phone": FilterColumn{
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "admin_users",
		},
		"is_active": FilterColumn{
			Operator: []string{"eq"},
			Variant:  "boolean",
			Table:    "admin_users",
		},
		"is_verified": FilterColumn{
			Operator: []string{"eq"},
			Variant:  "boolean",
			Table:    "admin_users",
		},
		"last_login_at": FilterColumn{
			Operator: []string{"eq", "ne", "lt", "lte", "gt", "gte"},
			Variant:  "time",
			Table:    "admin_users",
		},
		"code": FilterColumn{
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "admin_roles",
		},
	}
}

//func AllowedFilterColumnAdminUser() map[string][]string {
//	return map[string][]string{
//		"username":         []string{"eq", "like", "notLike"},
//		"email":            []string{"eq", "like", "notLike"},
//		"full_name":        []string{"eq", "like", "notLike"},
//		"phone":            []string{"eq", "like", "notLike"},
//		"is_active":        []string{"eq"},
//		"is_verified":      []string{"eq"},
//		"last_login_at":    []string{"eq", "ne", "lt", "lte", "gt", "gte"},
//		"admin_roles.code": []string{"eq"},
//	}
//}
