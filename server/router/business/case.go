package business

import (
	"github.com/gin-gonic/gin"
)

type CaseRouter struct{}

func (s *CaseRouter) InitCaseRouter(Router *gin.RouterGroup) {
	caseRouter := Router.Group("case")
	{
		caseRouter.POST("", caseApi.CreateCase)
		caseRouter.PUT("", caseApi.UpdateCase)
		caseRouter.DELETE("", caseApi.DeleteCase)
		caseRouter.GET("/:id", caseApi.GetCase)
		caseRouter.GET("/product", caseApi.GetCasesByProductID)
	}
}
