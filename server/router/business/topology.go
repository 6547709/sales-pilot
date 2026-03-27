package business

import (
	"github.com/gin-gonic/gin"
)

type TopologyRouter struct{}

func (s *TopologyRouter) InitTopologyRouter(Router *gin.RouterGroup) {
	layerRouter := Router.Group("topology/layer")
	{
		layerRouter.POST("", TopologyApiApp.CreateLayer)
		layerRouter.PUT("", TopologyApiApp.UpdateLayer)
		layerRouter.DELETE("", TopologyApiApp.DeleteLayer)
		layerRouter.GET("/all", TopologyApiApp.GetLayers)
	}

	categoryRouter := Router.Group("topology/category")
	{
		categoryRouter.POST("", TopologyApiApp.CreateCategory)
		categoryRouter.PUT("", TopologyApiApp.UpdateCategory)
		categoryRouter.DELETE("", TopologyApiApp.DeleteCategory)
		categoryRouter.GET("/by-layer", TopologyApiApp.GetCategoriesByLayer)
	}

	vendorRouter := Router.Group("topology/vendor")
	{
		vendorRouter.POST("", TopologyApiApp.CreateVendor)
		vendorRouter.PUT("", TopologyApiApp.UpdateVendor)
		vendorRouter.DELETE("", TopologyApiApp.DeleteVendor)
		vendorRouter.GET("/by-category", TopologyApiApp.GetVendorsByCategory)
	}
}
