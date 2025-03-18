package cmd

import (
	"os"
	"tiga-putra-cashier-be/database"

	"gorm.io/gorm"
)

func Command(db *gorm.DB) {
	migrateUp := false
	migrateDown := false

	for _, arg := range os.Args[1:] {
		if arg == "migrate-up" {
			migrateUp = true
		}
		if arg == "migrate-down" {
			migrateDown = true
		}
	}
	if migrateUp {
		if err := database.MigrateUp(db); err != nil {
			os.Exit(1)
		}
		os.Exit(0)

	}
	if migrateDown {
		if err := database.MigrateDown(db); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
}
