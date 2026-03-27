package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router"
	"github.com/gin-gonic/gin"
)

// 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}

func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0]
	publicGroup := routers[1]

	holder(publicGroup, privateGroup)

	// 注册业务路由
	businessRouter := router.RouterGroupApp.Business
	businessRouter.InitProductRouter(privateGroup) // 产品管理（需要鉴权）
	businessRouter.InitSalesScriptRouter(privateGroup) // 销售话术（需要鉴权）
	businessRouter.InitCaseRouter(privateGroup) // 客户案例（需要鉴权）
	businessRouter.InitTopologyRouter(privateGroup) // 拓扑管理（需要鉴权）
	businessRouter.InitSolutionCategoryRouter(privateGroup) // 解决方案分类（需要鉴权）
	businessRouter.InitAPIKeyRouter(privateGroup) // API密钥（需要鉴权）
	businessRouter.InitStatisticsRouter(privateGroup) // 统计（需要鉴权）
	businessRouter.InitBackupRouter(privateGroup) // 备份（需要鉴权）
}
