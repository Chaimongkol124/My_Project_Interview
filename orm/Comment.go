package orm

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Massage     string `gorm:"type:varchar(500) "`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      float64
	User        User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	InterviewID float64
	Interview   Interview `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
