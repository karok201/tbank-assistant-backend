package db

import (
	"log"
	"tbank-assistant-backend/pkg/common/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(url string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.Book{})
	db.AutoMigrate(&models.User{})

	return db
}
