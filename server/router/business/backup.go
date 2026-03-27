package business

import (
	"github.com/gin-gonic/gin"
)

type BackupRouter struct{}

func (s *BackupRouter) InitBackupRouter(Router *gin.RouterGroup) {
	backupRouter := Router.Group("backup")
	{
		backupRouter.GET("/export", backupApi.ExportBackup)
		backupRouter.POST("/import", backupApi.ImportBackup)
	}
}
