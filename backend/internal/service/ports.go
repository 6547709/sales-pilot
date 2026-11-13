package service

import (
	"context"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/paging"
	"github.com/6547709/sales-pilot/backend/internal/repository"
)

// 端口接口定义在消费侧，repository 包提供实现。

type UserRepo interface {
	ListLocal(ctx context.Context) ([]models.User, error)
	ListLocalPage(ctx context.Context, f repository.LocalUserFilter, p paging.Request) ([]models.User, int64, error)
	ListAll(ctx context.Context) ([]models.User, error)
	FindByUsernameAndProvider(ctx context.Context, username, provider string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	Create(ctx context.Context, u *models.User) error
	Save(ctx context.Context, u *models.User) error
	UpdateLastLogin(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
	CountActiveLocalAdmins(ctx context.Context) (int64, error)
}

type APIKeyRepo interface {
	List(ctx context.Context) ([]models.AdminAPIKey, error)
	FindByHash(ctx context.Context, hash string) (*models.AdminAPIKey, error)
	Create(ctx context.Context, k *models.AdminAPIKey) error
	Revoke(ctx context.Context, id uint) error
	TouchLastUsed(ctx context.Context, id uint) error
}

type ProductRepo interface {
	ListPublic(ctx context.Context, solutionID *uint, keyword, market, manufacturer string) ([]models.Product, error)
	FindByID(ctx context.Context, id uint) (*models.Product, error)
	FindPublishedByID(ctx context.Context, id uint) (*models.Product, error)
	ListAdmin(ctx context.Context, f repository.AdminProductFilter, p paging.Request) ([]models.Product, int64, error)
	ListAdminAll(ctx context.Context, f repository.AdminProductFilter) ([]models.Product, error)
	Create(ctx context.Context, p *models.Product) error
	Reload(ctx context.Context, p *models.Product) error
	UpdateFields(ctx context.Context, id uint, fields map[string]any) error
	UpdateDraft(ctx context.Context, id uint, isDraft bool) error
	Delete(ctx context.Context, id uint) error
	ListScripts(ctx context.Context, productID uint) ([]models.SalesScript, error)
	CreateScript(ctx context.Context, s *models.SalesScript) error
	FindScript(ctx context.Context, id uint) (*models.SalesScript, error)
	UpdateScript(ctx context.Context, s *models.SalesScript) error
	DeleteScript(ctx context.Context, id uint) error
	ListCases(ctx context.Context, productID uint) ([]models.Case, error)
	CreateCase(ctx context.Context, c *models.Case) error
	FindCase(ctx context.Context, id uint) (*models.Case, error)
	UpdateCase(ctx context.Context, c *models.Case) error
	UpdateCasePublish(ctx context.Context, id uint, isPublic bool, industry string) error
	DeleteCase(ctx context.Context, id uint) error
	FindCategoryLabel(ctx context.Context, id uint) (string, error)
	CountByCategory(ctx context.Context, categoryID uint) (int64, error)
}

type TopologyRepo interface {
	Snapshot(ctx context.Context) (*repository.TopologySnapshot, error)
	Categories(ctx context.Context, activeOnly bool) ([]models.SolutionCategory, error)
	FindCategory(ctx context.Context, id uint) (*models.SolutionCategory, error)
	CreateCategory(ctx context.Context, c *models.SolutionCategory) error
	UpdateCategoryFields(ctx context.Context, id uint, fields map[string]any) error
	DeleteCategory(ctx context.Context, id uint) error
	CountCategoriesInLayer(ctx context.Context, layerID uint) (int64, error)
	CountProductsInCategory(ctx context.Context, categoryID uint) (int64, error)
	Layers(ctx context.Context) ([]models.TopologyLayer, error)
	FindLayer(ctx context.Context, id uint) (*models.TopologyLayer, error)
	CreateLayer(ctx context.Context, l *models.TopologyLayer) error
	UpdateLayerFields(ctx context.Context, id uint, fields map[string]any) error
	DeleteLayer(ctx context.Context, id uint) error
	CreateVendor(ctx context.Context, v *models.SolutionVendor) error
	FindVendor(ctx context.Context, id uint) (*models.SolutionVendor, error)
	UpdateVendorFields(ctx context.Context, id uint, fields map[string]any) error
	DeleteVendor(ctx context.Context, id uint) error
}

type SettingsRepo interface {
	Auth(ctx context.Context) (*models.AuthSettings, error)
	SaveAuth(ctx context.Context, st *models.AuthSettings) error
	System(ctx context.Context) (*models.SystemSettings, error)
	SaveSystem(ctx context.Context, st *models.SystemSettings) error
}

type BackupRepo interface {
	Dump(ctx context.Context) (*models.BackupPayload, error)
	Restore(ctx context.Context, payload *models.BackupPayload) error
	ExportContent(ctx context.Context) (*models.BackupContent, error)
	ExportConfig(ctx context.Context) (*models.BackupConfig, error)
	ExportUsers(ctx context.Context) (*models.BackupUsers, error)
	RestoreSections(ctx context.Context, sections map[string]bool, content *models.BackupContent, cfg *models.BackupConfig, users *models.BackupUsers) error
}

type ServiceRepo interface {
	List(ctx context.Context, kind string, activeOnly bool) ([]models.Service, error)
	FindByID(ctx context.Context, id uint) (*models.Service, error)
	FindBySlug(ctx context.Context, slug string) (*models.Service, error)
	Create(ctx context.Context, s *models.Service) error
	Save(ctx context.Context, s *models.Service) error
	UpdateFields(ctx context.Context, id uint, fields map[string]any) error
	Delete(ctx context.Context, id uint) error
	SetCategories(ctx context.Context, serviceID uint, categoryIDs []uint) error
	CategoryIDs(ctx context.Context, serviceID uint) ([]uint, error)
	ListByCategoryID(ctx context.Context, categoryID uint) ([]models.Service, error)
	ListByCategoryIDs(ctx context.Context, ids []uint) (map[uint][]models.Service, error)
	ListByIDs(ctx context.Context, ids []uint) ([]models.Service, error)
}

type ScenarioRepo interface {
	Groups(ctx context.Context) ([]models.ScenarioGroup, error)
	Scenarios(ctx context.Context) ([]models.Scenario, error)
	FindBySlug(ctx context.Context, slug string) (*models.Scenario, error)
	FindByID(ctx context.Context, id uint) (*models.Scenario, error)
	FindGroup(ctx context.Context, id uint) (*models.ScenarioGroup, error)
	CreateGroup(ctx context.Context, g *models.ScenarioGroup) error
	SaveGroup(ctx context.Context, g *models.ScenarioGroup) error
	UpdateGroupFields(ctx context.Context, id uint, fields map[string]any) error
	DeleteGroup(ctx context.Context, id uint) error
	CountGroups(ctx context.Context) (int64, error)
	CountScenariosInGroup(ctx context.Context, groupID uint) (int64, error)
	Create(ctx context.Context, s *models.Scenario) error
	Save(ctx context.Context, s *models.Scenario) error
	UpdateFields(ctx context.Context, id uint, fields map[string]any) error
	Delete(ctx context.Context, id uint) error
	ProductIDs(ctx context.Context, scenarioID uint) ([]uint, error)
	ServiceIDs(ctx context.Context, scenarioID uint) ([]uint, error)
	SetLinks(ctx context.Context, scenarioID uint, productIDs, serviceIDs []uint) error
	CountProducts(ctx context.Context) (map[uint]int64, error)
	CountServices(ctx context.Context) (map[uint]int64, error)
	ListProductsByIDs(ctx context.Context, ids []uint) ([]models.Product, error)
	ListProductsByIDsAll(ctx context.Context, ids []uint) ([]models.Product, error)
	ListServicesByIDs(ctx context.Context, ids []uint) ([]models.Service, error)
	ListServicesByIDsAll(ctx context.Context, ids []uint) ([]models.Service, error)
	ListPublicCasesByProductIDs(ctx context.Context, ids []uint) ([]models.Case, error)
	ListCasesByProductIDs(ctx context.Context, ids []uint) ([]models.Case, error)
	ListByIDs(ctx context.Context, ids []uint) ([]models.Scenario, error)
	ProductNamesByScenario(ctx context.Context) (map[uint][]string, error)
	ServiceNamesByScenario(ctx context.Context) (map[uint][]string, error)
	ScenarioIDsForProduct(ctx context.Context, productID uint) ([]uint, error)
	ScenarioIDsForService(ctx context.Context, serviceID uint) ([]uint, error)
	SetScenariosForProduct(ctx context.Context, productID uint, scenarioIDs []uint) error
}
