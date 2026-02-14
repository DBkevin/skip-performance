package models

import (
	"time"
	"gorm.io/gorm"
)

type Employee struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string         `gorm:"type:varchar(32);not null" json:"name"`
	Role       string         `gorm:"type:varchar(20);not null;index:idx_role" json:"role"`
	Department *string        `gorm:"type:varchar(50)" json:"department,omitempty"`
	JobNumber  *string        `gorm:"type:varchar(32);uniqueIndex:uniq_job_number" json:"job_number,omitempty"`
	Phone      *string        `gorm:"type:varchar(20)" json:"phone,omitempty"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	Remark     *string        `gorm:"type:text" json:"remark,omitempty"`
	CreatedAt  *time.Time     `json:"created_at,omitempty"`
	UpdatedAt  *time.Time     `json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// Role constants
const (
	RoleDoctor     = "医生"
	RoleNurse      = "护士"
	RoleConsultant = "咨询师"
	RoleAdmin      = "管理员"
)

func (Employee) TableName() string {
	return "employees"
}

func (e Employee) IsAdmin() bool {
	return e.Role == RoleAdmin
}
