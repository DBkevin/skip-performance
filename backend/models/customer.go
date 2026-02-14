package models

import (
	"time"
	"gorm.io/gorm"
)

type Customer struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string         `gorm:"type:varchar(64);not null" json:"name"`
	Phone          string         `gorm:"type:varchar(20);not null;uniqueIndex:uniq_phone" json:"phone"`
	CustomerType   *string        `gorm:"type:varchar(20)" json:"customer_type,omitempty"`
	FirstVisitDate *time.Time     `gorm:"type:date" json:"first_visit_date,omitempty"`
	Remark         *string        `gorm:"type:text" json:"remark,omitempty"`
	CreatedAt      *time.Time     `json:"created_at,omitempty"`
	UpdatedAt      *time.Time     `json:"updated_at,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Customer) TableName() string {
	return "customers"
}
