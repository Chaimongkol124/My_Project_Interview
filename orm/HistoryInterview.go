package orm

import "gorm.io/gorm"

type HistoryInterview struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(30) "`
	Description string    `gorm:"type:varchar(500) "`
	Status      StausType `gorm:"type:ENUM('To Do', 'In Progress', 'Done')"`
	InterviewID float64
	Interview   Interview `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
