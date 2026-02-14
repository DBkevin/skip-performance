package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string         `gorm:"type:varchar(32);not null;uniqueIndex" json:"username"`
	Password   string         `gorm:"type:varchar(255);not null" json:"-"`
	EmployeeID *uint          `gorm:"index" json:"employee_id,omitempty"`
	Role       string         `gorm:"type:varchar(20);not null" json:"role"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  *time.Time     `json:"created_at,omitempty"`
	UpdatedAt  *time.Time     `json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Employee *Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
