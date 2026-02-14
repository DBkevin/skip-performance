package models

import (
	"time"
	"gorm.io/gorm"
)

type Project struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string         `gorm:"type:varchar(100);not null;uniqueIndex:uniq_name" json:"name"`
	Category      *string        `gorm:"type:varchar(50);index:idx_category" json:"category,omitempty"`
	StandardPrice *float64       `gorm:"type:decimal(10,2)" json:"standard_price,omitempty"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	Remark        *string        `gorm:"type:text" json:"remark,omitempty"`
	CreatedAt     *time.Time     `json:"created_at,omitempty"`
	UpdatedAt     *time.Time     `json:"updated_at,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Project) TableName() string {
	return "projects"
}
