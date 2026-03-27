package business

import (
	"github.com/gin-gonic/gin"
)

type ProductRouter struct{}

func (s *ProductRouter) InitProductRouter(Router *gin.RouterGroup) {
	productRouter := Router.Group("product")
	{
		productRouter.POST("", ProductApiApp.CreateProduct)
		productRouter.PUT("", ProductApiApp.UpdateProduct)
		productRouter.DELETE("", ProductApiApp.DeleteProduct)
		productRouter.GET("/:id", ProductApiApp.GetProduct)
		productRouter.POST("/list", ProductApiApp.GetProductList)
		productRouter.POST("/public/list", ProductApiApp.GetPublicProductList)
	}
}
