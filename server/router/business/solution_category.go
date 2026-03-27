package business

import (
	"github.com/gin-gonic/gin"
)

type SolutionCategoryRouter struct{}

func (s *SolutionCategoryRouter) InitSolutionCategoryRouter(Router *gin.RouterGroup) {
	catRouter := Router.Group("solution-category")
	{
		catRouter.POST("", solutionCategoryApi.Create)
		catRouter.PUT("", solutionCategoryApi.Update)
		catRouter.DELETE("", solutionCategoryApi.Delete)
		catRouter.GET("/:id", solutionCategoryApi.Get)
		catRouter.GET("/all", solutionCategoryApi.GetAll)
	}
}
