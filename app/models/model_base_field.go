package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseField struct {
	CreatedAt *time.Time     `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt *time.Time     `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;default:NULL" json:"-"`
}

type FilterColumn struct {
	Operator   []string
	Variant    string
	IsRequired bool
	JointTable string
	Table      string
}
