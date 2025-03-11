package router

import (
	"tiga-putra-cashier-be/controller"
	"tiga-putra-cashier-be/router/product"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.Engine, pc controller.ProductController) *gin.Engine {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))
	r.Static("/assets/image", "./assets/image")
	v1 := r.Group("/v1")
	{
		product.ProductRouter(v1, pc)
	}
	return r
}
