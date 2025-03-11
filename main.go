package main

import (
	"log"
	"os"
	"tiga-putra-cashier-be/cmd"
	"tiga-putra-cashier-be/controller"
	"tiga-putra-cashier-be/database"
	"tiga-putra-cashier-be/di"
	"tiga-putra-cashier-be/router"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	container := di.BuildContainer()
	err := container.Invoke(func(
		r *gin.Engine,
		db *gorm.DB,
		pc controller.ProductController,
	) {
		defer database.CloseDB(db)
		if len(os.Args) > 1 {
			cmd.Command(db)
		}
		router.AppRouter(r, pc)
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		var serve string
		if os.Getenv("APP_ENV") == "development" {
			serve = "127.0.0.1:" + port
		} else {
			serve = ":" + port
		}

		if err := r.Run(serve); err != nil {
			log.Fatalf("error running server: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
}
