package service

import (
	"context"
	"errors"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/paging"
	"github.com/6547709/sales-pilot/backend/internal/repository"
	"gorm.io/gorm"
)

type PublicProductFilter struct {
	SolutionCategoryID *uint
	Keyword            string
	Market             string
	Manufacturer       string
}

type ProductService struct {
	repo ProductRepo
}

func NewProductService(repo ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func NormalizeVendorMarket(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "domestic", "foreign":
		return strings.ToLower(strings.TrimSpace(v))
	default:
		return "all"
	}
}

func (s *ProductService) ListPublic(ctx context.Context, f PublicProductFilter) ([]models.Product, error) {
	list, err := s.repo.ListPublic(ctx, f.SolutionCategoryID, f.Keyword, f.Market, f.Manufacturer)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return list, nil
}

func (s *ProductService) GetPublic(ctx context.Context, id uint) (*models.Product, error) {
	p, err := s.repo.FindPublishedByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("product_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	return p, nil
}

func (s *ProductService) ListAdmin(ctx context.Context, f repository.AdminProductFilter, p paging.Request) ([]models.Product, int64, error) {
	list, total, err := s.repo.ListAdmin(ctx, f, p)
	if err != nil {
		return nil, 0, apperr.Internal(err)
	}
	return list, total, nil
}

func (s *ProductService) ListAdminAll(ctx context.Context, f repository.AdminProductFilter) ([]models.Product, error) {
	list, err := s.repo.ListAdminAll(ctx, f)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return list, nil
}

func (s *ProductService) GetAdmin(ctx context.Context, id uint) (*models.Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("product_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	return p, nil
}

func (s *ProductService) syncCategory(ctx context.Context, p *models.Product) {
	if p.SolutionCategoryID == nil || *p.SolutionCategoryID == 0 {
		return
	}
	if label, err := s.repo.FindCategoryLabel(ctx, *p.SolutionCategoryID); err == nil {
		p.Category = label
	}
}

func (s *ProductService) Create(ctx context.Context, p *models.Product) (*models.Product, error) {
	p.VendorMarket = NormalizeVendorMarket(p.VendorMarket)
	s.syncCategory(ctx, p)
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, apperr.Internal(err)
	}
	if err := s.repo.Reload(ctx, p); err != nil {
		return nil, apperr.Internal(err)
	}
	return p, nil
}

func (s *ProductService) Update(ctx context.Context, id uint, p *models.Product) (*models.Product, error) {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("product_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	p.ID = id
	p.VendorMarket = NormalizeVendorMarket(p.VendorMarket)
	s.syncCategory(ctx, p)
	fields := map[string]any{
		"name":                   p.Name,
		"category":               p.Category,
		"solution_category_id":   p.SolutionCategoryID,
		"vendor_market":          p.VendorMarket,
		"manufacturer_name":      p.ManufacturerName,
		"manufacturer_logo":      p.ManufacturerLogo,
		"sales_contact_name":     p.SalesContactName,
		"sales_contact_phone":    p.SalesContactPhone,
		"sales_contact_email":    p.SalesContactEmail,
		"presales_contact_name":  p.PresalesContactName,
		"presales_contact_phone": p.PresalesContactPhone,
		"presales_contact_email": p.PresalesContactEmail,
		"description":            p.Description,
		"highlights":             p.Highlights,
		"target_personas":        p.TargetPersonas,
		"trigger_events":         p.TriggerEvents,
		"discovery_questions":    p.DiscoveryQuestions,
		"competitor_analysis":    p.CompetitorAnalysis,
		"roi_metrics":            p.ROIMetrics,
		"is_draft":               p.IsDraft,
	}
	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		return nil, apperr.Internal(err)
	}
	if err := s.repo.Reload(ctx, p); err != nil {
		return nil, apperr.Internal(err)
	}
	return p, nil
}

func (s *ProductService) SetDraft(ctx context.Context, id uint, isDraft bool) (*models.Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("product_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	if err := s.repo.UpdateDraft(ctx, id, isDraft); err != nil {
		return nil, apperr.Internal(err)
	}
	if err := s.repo.Reload(ctx, p); err != nil {
		return nil, apperr.Internal(err)
	}
	return p, nil
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *ProductService) ListScriptsPublic(ctx context.Context, productID uint) ([]models.SalesScript, error) {
	if _, err := s.GetPublic(ctx, productID); err != nil {
		return nil, err
	}
	return s.ListScriptsAdmin(ctx, productID)
}

func (s *ProductService) ListScriptsAdmin(ctx context.Context, productID uint) ([]models.SalesScript, error) {
	list, err := s.repo.ListScripts(ctx, productID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return list, nil
}

func (s *ProductService) CreateScript(ctx context.Context, productID uint, sc *models.SalesScript) (*models.SalesScript, error) {
	sc.ProductID = productID
	if err := s.repo.CreateScript(ctx, sc); err != nil {
		return nil, apperr.Internal(err)
	}
	return sc, nil
}

func (s *ProductService) UpdateScript(ctx context.Context, id uint, sc *models.SalesScript) (*models.SalesScript, error) {
	if _, err := s.repo.FindScript(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("script_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	sc.ID = id
	if err := s.repo.UpdateScript(ctx, sc); err != nil {
		return nil, apperr.Internal(err)
	}
	return sc, nil
}

func (s *ProductService) DeleteScript(ctx context.Context, id uint) error {
	if err := s.repo.DeleteScript(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *ProductService) ListCasesPublic(ctx context.Context, productID uint) ([]models.Case, error) {
	if _, err := s.GetPublic(ctx, productID); err != nil {
		return nil, err
	}
	return s.ListCasesAdmin(ctx, productID)
}

func (s *ProductService) ListCasesAdmin(ctx context.Context, productID uint) ([]models.Case, error) {
	list, err := s.repo.ListCases(ctx, productID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return list, nil
}

func (s *ProductService) CreateCase(ctx context.Context, productID uint, c *models.Case) (*models.Case, error) {
	c.ProductID = productID
	if err := s.repo.CreateCase(ctx, c); err != nil {
		return nil, apperr.Internal(err)
	}
	return c, nil
}

func (s *ProductService) UpdateCase(ctx context.Context, id uint, c *models.Case) (*models.Case, error) {
	if _, err := s.repo.FindCase(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("case_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	c.ID = id
	if err := s.repo.UpdateCase(ctx, c); err != nil {
		return nil, apperr.Internal(err)
	}
	return c, nil
}

func (s *ProductService) DeleteCase(ctx context.Context, id uint) error {
	if err := s.repo.DeleteCase(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}
