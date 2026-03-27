package business

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/business"
	"github.com/gin-gonic/gin"
)

type BackupApi struct{}

// ExportBackup 导出所有业务数据
func (b *BackupApi) ExportBackup(c *gin.Context) {
	data, err := business.BackupServiceApp.ExportAll()
	if err != nil {
		response.FailWithMessage("导出失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "导出成功", c)
}

// ImportBackup 导入业务数据
func (b *BackupApi) ImportBackup(c *gin.Context) {
	var data business.BackupData
	if err := c.ShouldBindJSON(&data); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}
	if err := business.BackupServiceApp.ImportAll(&data); err != nil {
		response.FailWithMessage("导入失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("导入成功", c)
}
