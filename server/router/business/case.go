package business

import (
	"github.com/gin-gonic/gin"
)

type CaseRouter struct{}

func (s *CaseRouter) InitCaseRouter(Router *gin.RouterGroup) {
	caseRouter := Router.Group("case")
	{
		caseRouter.POST("", CaseApiApp.CreateCase)
		caseRouter.PUT("", CaseApiApp.UpdateCase)
		caseRouter.DELETE("", CaseApiApp.DeleteCase)
		caseRouter.GET("/:id", CaseApiApp.GetCase)
		caseRouter.GET("/product", CaseApiApp.GetCasesByProductID)
	}
}
