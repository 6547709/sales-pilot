package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type SalesScriptApi struct{}

func (b *SalesScriptApi) CreateSalesScript(c *gin.Context) {
	var script business.SalesScript
	err := c.ShouldBindJSON(&script)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = SalesScriptServiceApp.CreateSalesScript(&script)
	if err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(script, "创建成功", c)
}

func (b *SalesScriptApi) UpdateSalesScript(c *gin.Context) {
	var script business.SalesScript
	err := c.ShouldBindJSON(&script)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = SalesScriptServiceApp.UpdateSalesScript(&script)
	if err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(script, "更新成功", c)
}

func (b *SalesScriptApi) DeleteSalesScript(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = SalesScriptServiceApp.DeleteSalesScript(reqId.Uint())
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (b *SalesScriptApi) GetSalesScript(c *gin.Context) {
	id := c.Param("id")
	script, err := SalesScriptServiceApp.GetSalesScript(parseUint(id))
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(script, "获取成功", c)
}

func (b *SalesScriptApi) GetSalesScriptsByProductID(c *gin.Context) {
	productID := c.Query("product_id")
	scripts, err := SalesScriptServiceApp.GetSalesScriptsByProductID(parseUint(productID))
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(scripts, "获取成功", c)
}
