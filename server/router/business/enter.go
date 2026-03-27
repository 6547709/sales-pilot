package business

import (
	businessApi "github.com/flipped-aurora/gin-vue-admin/server/api/v1/business"
)

type RouterGroup struct {
	ProductRouter
	SalesScriptRouter
	CaseRouter
	TopologyRouter
	SolutionCategoryRouter
	APIKeyRouter
	StatisticsRouter
	BackupRouter
}

var (
	productApi              = businessApi.ProductApiApp
	salesScriptApi          = businessApi.SalesScriptApiApp
	caseApi                 = businessApi.CaseApiApp
	topologyApi             = businessApi.TopologyApiApp
	solutionCategoryApi     = businessApi.SolutionCategoryApiApp
	apiKeyApi               = businessApi.APIKeyApiApp
	statisticsApi           = businessApi.StatisticsApiApp
	backupApi               = businessApi.BackupApiApp
)

var (
	ProductRouterApp           = new(ProductRouter)
	SalesScriptRouterApp       = new(SalesScriptRouter)
	CaseRouterApp              = new(CaseRouter)
	TopologyRouterApp           = new(TopologyRouter)
	SolutionCategoryRouterApp   = new(SolutionCategoryRouter)
	APIKeyRouterApp            = new(APIKeyRouter)
	StatisticsRouterApp        = new(StatisticsRouter)
	BackupRouterApp            = new(BackupRouter)
)
