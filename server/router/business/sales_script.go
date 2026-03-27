package business

import (
	"github.com/gin-gonic/gin"
)

type SalesScriptRouter struct{}

func (s *SalesScriptRouter) InitSalesScriptRouter(Router *gin.RouterGroup) {
	scriptRouter := Router.Group("script")
	{
		scriptRouter.POST("", salesScriptApi.CreateSalesScript)
		scriptRouter.PUT("", salesScriptApi.UpdateSalesScript)
		scriptRouter.DELETE("", salesScriptApi.DeleteSalesScript)
		scriptRouter.GET("/:id", salesScriptApi.GetSalesScript)
		scriptRouter.GET("/product", salesScriptApi.GetSalesScriptsByProductID)
	}
}
