package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type TopologyApi struct{}

// Layer APIs
func (b *TopologyApi) CreateLayer(c *gin.Context) {
	var layer business.TopologyLayer
	err := c.ShouldBindJSON(&layer)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = TopologyServiceApp.CreateLayer(&layer)
	if err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(layer, "创建成功", c)
}

func (b *TopologyApi) UpdateLayer(c *gin.Context) {
	var layer business.TopologyLayer
	err := c.ShouldBindJSON(&layer)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = TopologyServiceApp.UpdateLayer(&layer)
	if err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(layer, "更新成功", c)
}

func (b *TopologyApi) DeleteLayer(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = TopologyServiceApp.DeleteLayer(reqId.Uint())
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (b *TopologyApi) GetLayers(c *gin.Context) {
	layers, err := TopologyServiceApp.GetAllLayers()
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(layers, "获取成功", c)
}

// Category APIs
func (b *TopologyApi) CreateCategory(c *gin.Context) {
	var cat business.TopologyCategory
	err := c.ShouldBindJSON(&cat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = TopologyServiceApp.CreateCategory(&cat)
	if err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(cat, "创建成功", c)
}

func (b *TopologyApi) UpdateCategory(c *gin.Context) {
	var cat business.TopologyCategory
	err := c.ShouldBindJSON(&cat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = TopologyServiceApp.UpdateCategory(&cat)
	if err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(cat, "更新成功", c)
}

func (b *TopologyApi) DeleteCategory(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = TopologyServiceApp.DeleteCategory(reqId.Uint())
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (b *TopologyApi) GetCategoriesByLayer(c *gin.Context) {
	layerID := c.Query("layer_id")
	categories, err := TopologyServiceApp.GetCategoriesByLayerID(parseUint(layerID))
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(categories, "获取成功", c)
}

// Vendor APIs
func (b *TopologyApi) CreateVendor(c *gin.Context) {
	var vendor business.TopologyVendor
	err := c.ShouldBindJSON(&vendor)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = TopologyServiceApp.CreateVendor(&vendor)
	if err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(vendor, "创建成功", c)
}

func (b *TopologyApi) UpdateVendor(c *gin.Context) {
	var vendor business.TopologyVendor
	err := c.ShouldBindJSON(&vendor)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = TopologyServiceApp.UpdateVendor(&vendor)
	if err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(vendor, "更新成功", c)
}

func (b *TopologyApi) DeleteVendor(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = TopologyServiceApp.DeleteVendor(reqId.Uint())
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (b *TopologyApi) GetVendorsByCategory(c *gin.Context) {
	categoryID := c.Query("category_id")
	vendors, err := TopologyServiceApp.GetVendorsByCategoryID(parseUint(categoryID))
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(vendors, "获取成功", c)
}
