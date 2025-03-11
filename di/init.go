package di

import (
	"tiga-putra-cashier-be/controller"
	"tiga-putra-cashier-be/database"
	"tiga-putra-cashier-be/repository"
	"tiga-putra-cashier-be/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	container.Provide(database.InitDB)
	container.Provide(repository.NewProductRepository)
	container.Provide(service.NewProductService)
	container.Provide(controller.NewProductController)
	container.Provide(gin.Default)
	return container
}
