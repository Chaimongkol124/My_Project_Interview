package orm

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(20);uniqueIndex "`
	Password  string `gorm:"type:varchar(100) "`
	FirstName string `gorm:"type:varchar(20)"`
	LastName  string `gorm:"type:varchar(20)"`
	Email     string `gorm:"type:varchar(20);uniqueIndex"`
	Avatar    string `gorm:"type:varchar(500)"`
}
