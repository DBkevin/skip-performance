package models

import (
	"time"
	"gorm.io/gorm"
)

type VisitItem struct {
	ID                    uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	VisitID               uint           `gorm:"not null;index:idx_visit_id" json:"visit_id"`
	ProjectID             uint           `gorm:"not null;index:idx_project_id" json:"project_id"`
	Amount                float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	MainDoctorID          uint           `gorm:"not null;index:idx_main_doctor_id" json:"main_doctor_id"`
	CoDoctor1ID           *uint          `gorm:"index:idx_co_doctor1_id" json:"co_doctor1_id,omitempty"`
	CoRatio1              float64        `gorm:"type:decimal(3,2);default:0" json:"co_ratio1"`
	CoDoctor2ID           *uint          `gorm:"index:idx_co_doctor2_id" json:"co_doctor2_id,omitempty"`
	CoRatio2              float64        `gorm:"type:decimal(3,2);default:0" json:"co_ratio2"`
	Nurse1ID              *uint          `gorm:"index:idx_nurse1_id" json:"nurse1_id,omitempty"`
	Nurse2ID              *uint          `gorm:"index:idx_nurse2_id" json:"nurse2_id,omitempty"`
	MainDoctorPerformance float64        `gorm:"type:decimal(10,2);default:0" json:"main_doctor_performance"`
	CoDoctor1Performance  float64        `gorm:"type:decimal(10,2);default:0" json:"co_doctor1_performance"`
	CoDoctor2Performance  float64        `gorm:"type:decimal(10,2);default:0" json:"co_doctor2_performance"`
	Nurse1Performance     float64        `gorm:"type:decimal(10,2);default:0" json:"nurse1_performance"`
	Nurse2Performance     float64        `gorm:"type:decimal(10,2);default:0" json:"nurse2_performance"`
	Remark                *string        `gorm:"type:text" json:"remark,omitempty"`
	CreatedAt             *time.Time     `json:"created_at,omitempty"`
	UpdatedAt             *time.Time     `json:"updated_at,omitempty"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Visit        Visit    `gorm:"foreignKey:VisitID" json:"visit,omitempty"`
	Project      Project  `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	MainDoctor   Employee `gorm:"foreignKey:MainDoctorID" json:"main_doctor,omitempty"`
	CoDoctor1    *Employee `gorm:"foreignKey:CoDoctor1ID" json:"co_doctor1,omitempty"`
	CoDoctor2    *Employee `gorm:"foreignKey:CoDoctor2ID" json:"co_doctor2,omitempty"`
	Nurse1       *Employee `gorm:"foreignKey:Nurse1ID" json:"nurse1,omitempty"`
	Nurse2       *Employee `gorm:"foreignKey:Nurse2ID" json:"nurse2,omitempty"`
}

func (VisitItem) TableName() string {
	return "visit_items"
}
