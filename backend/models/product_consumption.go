package models

import (
	"time"
	"gorm.io/gorm"
)

type ProductConsumption struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	VisitItemID  uint           `gorm:"not null;index:idx_visit_item_id" json:"visit_item_id"`
	ProductName  string         `gorm:"type:varchar(100);not null" json:"product_name"`
	Quantity     int            `gorm:"not null" json:"quantity"`
	UnitPrice    *float64       `gorm:"type:decimal(10,2)" json:"unit_price,omitempty"`
	Remark       *string        `gorm:"type:text" json:"remark,omitempty"`
	CreatedAt    *time.Time     `json:"created_at,omitempty"`
	UpdatedAt    *time.Time     `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	VisitItem VisitItem `gorm:"foreignKey:VisitItemID" json:"visit_item,omitempty"`
}

func (ProductConsumption) TableName() string {
	return "product_consumption"
}
