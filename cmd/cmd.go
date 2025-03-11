package cmd

import (
	"os"
	"tiga-putra-cashier-be/database"

	"gorm.io/gorm"
)

func Command(db *gorm.DB) {
	migrate := false

	for _, arg := range os.Args[1:] {
		if arg == "migrate" {
			migrate = true
		}
	}
	if migrate {
		database.Migrate(db)
	}
}
