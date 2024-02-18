package orm

import (
	"time"

	"gorm.io/gorm"
)

type Interview struct {
	gorm.Model

	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      StausType `gorm:"type:ENUM('To Do', 'In Progress', 'Done')"`
	UserID      float64
	User        User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type StausType string

// const available value for enum
const (
	ToDo       StausType = "To Do"
	InProgress StausType = "In Progress"
	Done       StausType = "Done"
)
