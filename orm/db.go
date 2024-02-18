package orm

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error

func InitDb() {
	dsn := os.Getenv("MYSQL_DNS")
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	Db.AutoMigrate(&User{})
	Db.AutoMigrate(&Interview{})
	Db.AutoMigrate(&Comment{})
	Db.AutoMigrate(&HistoryInterview{})
}

func (t StausType) Value() (driver.Value, error) {
	switch t {
	case ToDo, InProgress, Done: //valid case
		return string(t), nil
	}
	return nil, errors.New("Invalid status type value") //else is invalid
}

// Scan validate enum on read from data base
func (t *StausType) Scan(value interface{}) error {
	var pt StausType
	if value == nil {
		*t = ""
		return nil
	}
	st, ok := value.([]uint8) // if we declare db type as ENUM gorm will scan value as []uint8
	if !ok {
		return errors.New("Invalid data for staus type")
	}
	pt = StausType(string(st)) //convert type from string to ProductType

	switch pt {
	case ToDo, InProgress, Done: //valid case
		*t = pt
		return nil
	}
	return fmt.Errorf("Invalid staus type value :%s", st) //else is invalid
}
