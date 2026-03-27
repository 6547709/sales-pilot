package business

import (
	"github.com/gin-gonic/gin"
)

type TopologyRouter struct{}

func (s *TopologyRouter) InitTopologyRouter(Router *gin.RouterGroup) {
	layerRouter := Router.Group("topology/layer")
	{
		layerRouter.POST("", topologyApi.CreateLayer)
		layerRouter.PUT("", topologyApi.UpdateLayer)
		layerRouter.DELETE("", topologyApi.DeleteLayer)
		layerRouter.GET("/all", topologyApi.GetLayers)
	}

	categoryRouter := Router.Group("topology/category")
	{
		categoryRouter.POST("", topologyApi.CreateCategory)
		categoryRouter.PUT("", topologyApi.UpdateCategory)
		categoryRouter.DELETE("", topologyApi.DeleteCategory)
		categoryRouter.GET("/by-layer", topologyApi.GetCategoriesByLayer)
	}

	vendorRouter := Router.Group("topology/vendor")
	{
		vendorRouter.POST("", topologyApi.CreateVendor)
		vendorRouter.PUT("", topologyApi.UpdateVendor)
		vendorRouter.DELETE("", topologyApi.DeleteVendor)
		vendorRouter.GET("/by-category", topologyApi.GetVendorsByCategory)
	}
}
