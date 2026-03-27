package business

import (
	"github.com/gin-gonic/gin"
)

type TopologyRouter struct{}

func (s *TopologyRouter) InitTopologyRouter(privateRouter, publicRouter *gin.RouterGroup) {
	// 公开接口 - 首页全景图
	topologyPublicRouter := publicRouter.Group("topology")
	{
		topologyPublicRouter.GET("/full", topologyApi.GetFullTopology) // 供首页全景图使用
	}

	// 私有接口 - 拓扑管理（需要鉴权）
	layerRouter := privateRouter.Group("topology/layer")
	{
		layerRouter.POST("", topologyApi.CreateLayer)
		layerRouter.PUT("", topologyApi.UpdateLayer)
		layerRouter.DELETE("", topologyApi.DeleteLayer)
		layerRouter.GET("/all", topologyApi.GetLayers)
	}

	categoryRouter := privateRouter.Group("topology/category")
	{
		categoryRouter.POST("", topologyApi.CreateCategory)
		categoryRouter.PUT("", topologyApi.UpdateCategory)
		categoryRouter.DELETE("", topologyApi.DeleteCategory)
		categoryRouter.GET("/by-layer", topologyApi.GetCategoriesByLayer)
	}

	vendorRouter := privateRouter.Group("topology/vendor")
	{
		vendorRouter.POST("", topologyApi.CreateVendor)
		vendorRouter.PUT("", topologyApi.UpdateVendor)
		vendorRouter.DELETE("", topologyApi.DeleteVendor)
		vendorRouter.GET("/by-category", topologyApi.GetVendorsByCategory)
	}
}
