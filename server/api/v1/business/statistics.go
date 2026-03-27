package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/business"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type StatisticsApi struct{}

// StatisticsOverview 业务统计数据
type StatisticsOverview struct {
	ProductCount    int64 `json:"product_count"`
	ScriptCount     int64 `json:"script_count"`
	CaseCount       int64 `json:"case_count"`
	CategoryCount   int64 `json:"category_count"`
	LayerCount      int64 `json:"layer_count"`
}

// GetOverview 获取业务统计数据
func (s *StatisticsApi) GetOverview(c *gin.Context) {
	var stats StatisticsOverview

	global.GVA_DB.Model(&business.Product{}).Count(&stats.ProductCount)
	global.GVA_DB.Model(&business.SalesScript{}).Count(&stats.ScriptCount)
	global.GVA_DB.Model(&business.Case{}).Count(&stats.CaseCount)
	global.GVA_DB.Model(&business.SolutionCategory{}).Count(&stats.CategoryCount)
	global.GVA_DB.Model(&business.TopologyLayer{}).Count(&stats.LayerCount)

	response.OkWithDetailed(stats, "获取成功", c)
}
