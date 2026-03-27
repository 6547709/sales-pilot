package business

import (
	"github.com/gin-gonic/gin"
)

type SolutionCategoryRouter struct{}

func (s *SolutionCategoryRouter) InitSolutionCategoryRouter(Router *gin.RouterGroup) {
	catRouter := Router.Group("solution-category")
	{
		catRouter.POST("", SolutionCategoryApiApp.Create)
		catRouter.PUT("", SolutionCategoryApiApp.Update)
		catRouter.DELETE("", SolutionCategoryApiApp.Delete)
		catRouter.GET("/:id", SolutionCategoryApiApp.Get)
		catRouter.GET("/all", SolutionCategoryApiApp.GetAll)
	}
}
