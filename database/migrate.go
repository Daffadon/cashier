package database

import (
	"log"
	"tiga-putra-cashier-be/entity"

	"gorm.io/gorm"
)

func MigrateUp(db *gorm.DB) error {
	err := db.AutoMigrate(&entity.Product{})
	if err != nil {
		log.Println("Migration has been processed")
		return err
	}
	return nil
}

func MigrateDown(db *gorm.DB) error {
	err := db.Migrator().DropTable(&entity.Product{})
	if err != nil {
		log.Println("Migration has been rolled back")
		return err
	}
	return nil
}
