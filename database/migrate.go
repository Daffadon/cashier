package database

import (
	"log"
	"tiga-putra-cashier-be/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&entity.Product{})
	if err != nil {
		panic(err)
	}
	log.Println("Migration has been processed")
}
