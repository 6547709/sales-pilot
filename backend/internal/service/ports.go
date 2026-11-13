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
	ListLocalPage(ctx context.Context, p paging.Request) ([]models.User, int64, error)
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
}
