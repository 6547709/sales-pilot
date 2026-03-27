package business

import (
	"github.com/gin-gonic/gin"
)

type ProductRouter struct{}

func (s *ProductRouter) InitProductRouter(Router *gin.RouterGroup) {
	productRouter := Router.Group("product")
	{
		productRouter.POST("", productApi.CreateProduct)
		productRouter.PUT("", productApi.UpdateProduct)
		productRouter.DELETE("", productApi.DeleteProduct)
		productRouter.GET("/:id", productApi.GetProduct)
		productRouter.POST("/list", productApi.GetProductList)
		productRouter.POST("/public/list", productApi.GetPublicProductList)
	}
}
