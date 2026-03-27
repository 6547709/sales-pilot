package business

import (
	"github.com/gin-gonic/gin"
)

type ProductRouter struct{}

func (s *ProductRouter) InitProductRouter(privateRouter, publicRouter *gin.RouterGroup) {
	// 公开接口 - 首页产品列表
	productPublicRouter := publicRouter.Group("product")
	{
		productPublicRouter.POST("/public/list", productApi.GetPublicProductList)
		productPublicRouter.GET("/public/all", productApi.GetAllPublicProducts)      // 供首页使用
		productPublicRouter.GET("/public/by-category", productApi.GetProductsByCategory) // 按分类获取
	}

	// 私有接口 - 产品管理（需要鉴权）
	productRouter := privateRouter.Group("product")
	{
		productRouter.POST("", productApi.CreateProduct)
		productRouter.PUT("", productApi.UpdateProduct)
		productRouter.DELETE("", productApi.DeleteProduct)
		productRouter.GET("/:id", productApi.GetProduct)
		productRouter.POST("/list", productApi.GetProductList)
	}
}
