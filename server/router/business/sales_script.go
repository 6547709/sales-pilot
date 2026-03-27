package business

import (
	"github.com/gin-gonic/gin"
)

type SalesScriptRouter struct{}

func (s *SalesScriptRouter) InitSalesScriptRouter(Router *gin.RouterGroup) {
	scriptRouter := Router.Group("script")
	{
		scriptRouter.POST("", SalesScriptApiApp.CreateSalesScript)
		scriptRouter.PUT("", SalesScriptApiApp.UpdateSalesScript)
		scriptRouter.DELETE("", SalesScriptApiApp.DeleteSalesScript)
		scriptRouter.GET("/:id", SalesScriptApiApp.GetSalesScript)
		scriptRouter.GET("/product", SalesScriptApiApp.GetSalesScriptsByProductID)
	}
}
