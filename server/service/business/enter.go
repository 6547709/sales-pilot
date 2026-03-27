package business

type ServiceGroup struct {
	ProductService
	SalesScriptService
	CaseService
	TopologyService
	SolutionCategoryService
	APIKeyService
	BackupService
}

var ServiceGroupApp = new(ServiceGroup)

var (
	ProductServiceApp           = new(ProductService)
	SalesScriptServiceApp       = new(SalesScriptService)
	CaseServiceApp              = new(CaseService)
	TopologyServiceApp           = new(TopologyService)
	SolutionCategoryServiceApp   = new(SolutionCategoryService)
	APIKeyServiceApp            = new(APIKeyService)
	BackupServiceApp            = new(BackupService)
)
