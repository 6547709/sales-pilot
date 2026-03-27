package business

import (
	"github.com/gin-gonic/gin"
)

type APIKeyRouter struct{}

func (s *APIKeyRouter) InitAPIKeyRouter(Router *gin.RouterGroup) {
	apiKeyRouter := Router.Group("api-key")
	{
		apiKeyRouter.POST("", APIKeyApiApp.Create)
		apiKeyRouter.GET("/list", APIKeyApiApp.GetAPIKeys)
		apiKeyRouter.POST("/revoke", APIKeyApiApp.Revoke)
		apiKeyRouter.DELETE("", APIKeyApiApp.Delete)
		apiKeyRouter.POST("/validate", APIKeyApiApp.Validate)
	}
}
