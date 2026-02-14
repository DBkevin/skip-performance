package models

import (
	"time"
	"gorm.io/gorm"
)

type Visit struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	VisitID       string         `gorm:"type:varchar(64);not null;uniqueIndex:uniq_visit_id" json:"visit_id"`
	CustomerID    uint           `gorm:"not null;index:idx_customer_id" json:"customer_id"`
	ConsultantID  *uint          `gorm:"index:idx_consultant_id" json:"consultant_id,omitempty"`
	VisitDate     time.Time      `gorm:"not null;index:idx_visit_date" json:"visit_date"`
	TotalAmount   float64        `gorm:"type:decimal(10,2);default:0" json:"total_amount"`
	Remark        *string        `gorm:"type:text" json:"remark,omitempty"`
	CreatedAt     *time.Time     `json:"created_at,omitempty"`
	UpdatedAt     *time.Time     `json:"updated_at,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Customer   Customer   `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Consultant *Employee  `gorm:"foreignKey:ConsultantID" json:"consultant,omitempty"`
	Items      []VisitItem `gorm:"foreignKey:VisitID;references:ID" json:"items,omitempty"`
}

func (Visit) TableName() string {
	return "visits"
}
