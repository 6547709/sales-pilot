package business

import (
	"github.com/gin-gonic/gin"
)

type APIKeyRouter struct{}

func (s *APIKeyRouter) InitAPIKeyRouter(Router *gin.RouterGroup) {
	apiKeyRouter := Router.Group("api-key")
	{
		apiKeyRouter.POST("", apiKeyApi.Create)
		apiKeyRouter.GET("/list", apiKeyApi.GetAPIKeys)
		apiKeyRouter.POST("/revoke", apiKeyApi.Revoke)
		apiKeyRouter.DELETE("", apiKeyApi.Delete)
		apiKeyRouter.POST("/validate", apiKeyApi.Validate)
	}
}
