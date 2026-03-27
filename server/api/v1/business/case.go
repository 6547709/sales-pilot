package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type CaseApi struct{}

func (b *CaseApi) CreateCase(c *gin.Context) {
	var caseModel business.Case
	err := c.ShouldBindJSON(&caseModel)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = CaseServiceApp.CreateCase(&caseModel)
	if err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(caseModel, "创建成功", c)
}

func (b *CaseApi) UpdateCase(c *gin.Context) {
	var caseModel business.Case
	err := c.ShouldBindJSON(&caseModel)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = CaseServiceApp.UpdateCase(&caseModel)
	if err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(caseModel, "更新成功", c)
}

func (b *CaseApi) DeleteCase(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = CaseServiceApp.DeleteCase(reqId.Uint())
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (b *CaseApi) GetCase(c *gin.Context) {
	id := c.Param("id")
	caseModel, err := CaseServiceApp.GetCase(parseUint(id))
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(caseModel, "获取成功", c)
}

func (b *CaseApi) GetCasesByProductID(c *gin.Context) {
	productID := c.Query("product_id")
	cases, err := CaseServiceApp.GetCasesByProductID(parseUint(productID))
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(cases, "获取成功", c)
}
