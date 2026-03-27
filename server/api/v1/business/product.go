package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type ProductApi struct{}

func (b *ProductApi) CreateProduct(c *gin.Context) {
	var product business.Product
	err := c.ShouldBindJSON(&product)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ProductServiceApp.CreateProduct(&product)
	if err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(product, "创建成功", c)
}

func (b *ProductApi) UpdateProduct(c *gin.Context) {
	var product business.Product
	err := c.ShouldBindJSON(&product)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ProductServiceApp.UpdateProduct(&product)
	if err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(product, "更新成功", c)
}

func (b *ProductApi) DeleteProduct(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ProductServiceApp.DeleteProduct(reqId.Uint())
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (b *ProductApi) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := ProductServiceApp.GetProduct(parseUint(id))
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	// 获取关联的销售话术和客户案例
	scripts, _ := ProductServiceApp.GetScriptsByProductID(parseUint(id))
	cases, _ := ProductServiceApp.GetCasesByProductID(parseUint(id))
	response.OkWithDetailed(gin.H{
		"product": product,
		"scripts": scripts,
		"cases":  cases,
	}, "获取成功", c)
}

func (b *ProductApi) GetProductList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := ProductServiceApp.GetProductList(&pageInfo)
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (b *ProductApi) GetPublicProductList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	vendorMarket := c.Query("vendor_market")
	list, total, err := ProductServiceApp.GetPublicProductList(&pageInfo, vendorMarket)
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetAllPublicProducts 获取所有公开产品（供首页使用）
func (b *ProductApi) GetAllPublicProducts(c *gin.Context) {
	list, err := ProductServiceApp.GetAllPublicProducts()
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetProductsByCategory 获取指定分类的产品（供首页使用）
func (b *ProductApi) GetProductsByCategory(c *gin.Context) {
	categoryID := parseUint(c.Query("solution_category_id"))
	vendorMarket := c.Query("vendor_market")
	list, err := ProductServiceApp.GetProductsByCategory(categoryID, vendorMarket)
	if err != nil {
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

func parseUint(s string) uint {
	var u uint
	for _, c := range s {
		if c >= '0' && c <= '9' {
			u = u*10 + uint(c-'0')
		}
	}
	return u
}
