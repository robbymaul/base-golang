package web

import (
	"paymentserviceklink/app/enums"
	"time"
)

type AdminLoginRequest struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

type AdminSessionResponse struct {
	Token     string              `json:"token"`
	Role      enums.CodeAdminRole `json:"role"`
	ExpiredAt int64               `json:"expiredAt"`
}

type DetailAdminResponse struct {
	Id                string             `json:"id,omitempty"`
	Username          string             `json:"username,omitempty"`
	Email             string             `json:"email,omitempty"`
	FullName          string             `json:"fullName,omitempty"`
	Phone             string             `json:"phone,omitempty"`
	Avatar            string             `json:"avatar,omitempty"`
	RoleId            int64              `json:"roleId,omitempty"`
	IsActive          bool               `json:"isActive"`
	IsVerified        bool               `json:"isVerified"`
	LastLoginAt       *time.Time         `json:"lastLoginAt,omitempty"`
	LastLoginIp       string             `json:"lastLoginIp,omitempty"`
	LockedUntil       *time.Time         `json:"lockedUntil,omitempty"`
	PasswordChangedAt *time.Time         `json:"passwordChangedAt,omitempty"`
	TwoFactorEnabled  bool               `json:"twoFactorEnabled,omitempty"`
	TwoFactorSecret   string             `json:"twoFactorSecret,omitempty"`
	CreatedAt         *time.Time         `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time         `json:"updatedAt,omitempty"`
	Role              *AdminRoleResponse `json:"role,omitempty"`
}

type AdminRoleResponse struct {
	Id          int64                `json:"id,omitempty"`
	Code        enums.CodeAdminRole  `json:"code,omitempty"`
	Name        enums.NameAdminRoles `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	IsActive    bool                 `json:"isActive"`
	CreatedAt   *time.Time           `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time           `json:"updatedAt,omitempty"`
}

type CreateAdminRolesRequest struct {
	Code        enums.CodeAdminRole  `json:"code"`
	Name        enums.NameAdminRoles `json:"name"`
	Description string               `json:"description"`
}

type CreateAdminUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FullName  string `json:"fullName"`
	Phone     string `json:"phone"`
	AvatarUrl string `json:"avatarUrl"`
	RoleId    int64  `json:"roleId"`
}
