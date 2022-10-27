package db

import (
	"fantasy/database/entities"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitConnection() {
	db, err := gorm.Open("mysql", GetDSN())
	if err != nil {
		log.Fatalf("Database Connection Failed. %v", err)
	}

	db.AutoMigrate(&entities.Player{})
	DB = db
}

func GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(127.0.0.1:3306)/fantasy?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
	)
}
