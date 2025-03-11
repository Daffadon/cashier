package product

import (
	"tiga-putra-cashier-be/controller"

	"github.com/gin-gonic/gin"
)

func ProductRouter(router *gin.RouterGroup, pc controller.ProductController) {
	productRoutes := router.Group("/product")
	{
		productRoutes.GET("", pc.GetProduct)
		productRoutes.GET("/:barcode_id", pc.GetProductDetail) //get product detail
		productRoutes.GET("/search", pc.SearchProduct)
		productRoutes.POST("", pc.AddProduct)
		productRoutes.PATCH("/:barcode_id", pc.UpdateProduct)
		productRoutes.DELETE("/:barcode_id", pc.DeleteProduct)
	}
}
