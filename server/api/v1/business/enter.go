package business

import "github.com/flipped-aurora/gin-vue-admin/server/service/business"

type ApiGroup struct {
	ProductApi
	SalesScriptApi
	CaseApi
	TopologyApi
	SolutionCategoryApi
	APIKeyApi
}

var (
	ProductApiApp             = new(ProductApi)
	SalesScriptApiApp         = new(SalesScriptApi)
	CaseApiApp                = new(CaseApi)
	TopologyApiApp            = new(TopologyApi)
	SolutionCategoryApiApp    = new(SolutionCategoryApi)
	APIKeyApiApp              = new(APIKeyApi)
)

var (
	ProductServiceApp             = business.ProductServiceApp
	SalesScriptServiceApp         = business.SalesScriptServiceApp
	CaseServiceApp                = business.CaseServiceApp
	TopologyServiceApp            = business.TopologyServiceApp
	SolutionCategoryServiceApp    = business.SolutionCategoryServiceApp
	APIKeyServiceApp              = business.APIKeyServiceApp
)
