package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type SolutionCategoryApi struct{}

func (b *SolutionCategoryApi) Create(c *gin.Context) {
	var cat business.SolutionCategory
	err := c.ShouldBindJSON(&cat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = SolutionCategoryServiceApp.Create(&cat)
	if err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(cat, "创建成功", c)
}

func (b *SolutionCategoryApi) Update(c *gin.Context) {
	var cat business.SolutionCategory
	err := c.ShouldBindJSON(&cat)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = SolutionCategoryServiceApp.Update(&cat)
	if err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(cat, "更新成功", c)
}

func (b *SolutionCategoryApi) Delete(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = SolutionCategoryServiceApp.Delete(reqId.Uint())
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (b *SolutionCategoryApi) Get(c *gin.Context) {
	id := c.Param("id")
	cat, err := SolutionCategoryServiceApp.Get(parseUint(id))
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(cat, "获取成功", c)
}

func (b *SolutionCategoryApi) GetAll(c *gin.Context) {
	cats, err := SolutionCategoryServiceApp.GetAll()
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(cats, "获取成功", c)
}
