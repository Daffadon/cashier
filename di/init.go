package di

import (
	"log"
	"tiga-putra-cashier-be/controller"
	"tiga-putra-cashier-be/database"
	"tiga-putra-cashier-be/repository"
	"tiga-putra-cashier-be/service"
	"tiga-putra-cashier-be/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	if err := container.Provide(database.InitDB); err != nil {
		log.Fatalf("Failed to provide database: %v", err)
	}

	if err := container.Provide(utils.FileInit); err != nil {
		log.Fatalf("Failed to provide file utils: %v", err)
	}

	if err := container.Provide(repository.NewProductRepository); err != nil {
		log.Fatalf("Failed to provide product repository: %v", err)
	}
	if err := container.Provide(service.NewProductService); err != nil {
		log.Fatalf("Failed to provide product service: %v", err)
	}

	if err := container.Provide(controller.NewProductController); err != nil {
		log.Fatalf("Failed to provide product controller: %v", err)
	}

	if err := container.Provide(gin.Default); err != nil {
		log.Fatalf("Failed to provide gin default instance: %v", err)
	}
	return container
}
