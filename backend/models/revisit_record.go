package models

import (
	"time"
	"gorm.io/gorm"
)

type RevisitRecord struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	NurseID         uint           `gorm:"not null" json:"nurse_id"`
	Date            time.Time      `gorm:"type:date;not null;index:idx_date" json:"date"`
	ReceptionCount  int            `gorm:"default:0" json:"reception_count"`
	AddWechatCount  int            `gorm:"default:0" json:"add_wechat_count"`
	RevisitCount    int            `gorm:"default:0" json:"revisit_count"`
	Remark          *string        `gorm:"type:text" json:"remark,omitempty"`
	CreatedAt       *time.Time     `json:"created_at,omitempty"`
	UpdatedAt       *time.Time     `json:"updated_at,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Nurse Employee `gorm:"foreignKey:NurseID" json:"nurse,omitempty"`
}

func (RevisitRecord) TableName() string {
	return "revisit_records"
}
