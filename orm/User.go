package orm

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"uniqueIndex"`
	Password  string
	FirstName string
	LastName  string
	Email     string `gorm:"uniqueIndex"`
	Avatar    string
}
