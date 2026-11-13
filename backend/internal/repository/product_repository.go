package repository

import (
	"context"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/paging"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

type AdminProductFilter struct {
	Q      string
	Status string // draft | published
	Market string // all | domestic | foreign
}

func applyVendorMarketFilter(db *gorm.DB, raw string) *gorm.DB {
	m := strings.TrimSpace(strings.ToLower(raw))
	if m != "domestic" && m != "foreign" {
		return db
	}
	return db.Where("LOWER(COALESCE(NULLIF(TRIM(vendor_market), ''), 'all')) IN ?", []string{"all", m})
}

func (r *ProductRepository) ListPublic(ctx context.Context, solutionID *uint, keyword, market, manufacturer string) ([]models.Product, error) {
	q := r.db.WithContext(ctx).Model(&models.Product{}).Where("is_draft = ?", false)
	q = applyVendorMarketFilter(q, market)
	if solutionID != nil && *solutionID != 0 {
		q = q.Where("solution_category_id = ?", *solutionID)
	}
	if mf := strings.TrimSpace(manufacturer); mf != "" {
		q = q.Where("LOWER(TRIM(manufacturer_name)) = LOWER(?)", mf)
	}
	if kw := strings.TrimSpace(strings.ToLower(keyword)); kw != "" {
		pat := "%" + kw + "%"
		q = q.Where("LOWER(name) LIKE ? OR LOWER(category) LIKE ? OR LOWER(description) LIKE ?", pat, pat, pat)
	}
	var list []models.Product
	err := q.Order("updated_at desc").Find(&list).Error
	return list, err
}

func (r *ProductRepository) FindByID(ctx context.Context, id uint) (*models.Product, error) {
	var p models.Product
	if err := r.db.WithContext(ctx).Preload("SolutionCategory").First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) FindPublishedByID(ctx context.Context, id uint) (*models.Product, error) {
	var p models.Product
	if err := r.db.WithContext(ctx).Preload("SolutionCategory").First(&p, id).Error; err != nil {
		return nil, err
	}
	if p.IsDraft {
		return nil, gorm.ErrRecordNotFound
	}
	return &p, nil
}

func (r *ProductRepository) adminQuery(ctx context.Context, f AdminProductFilter) *gorm.DB {
	q := r.db.WithContext(ctx).Model(&models.Product{})
	if kw := strings.TrimSpace(strings.ToLower(f.Q)); kw != "" {
		pat := "%" + kw + "%"
		q = q.Where("LOWER(name) LIKE ? OR LOWER(category) LIKE ?", pat, pat)
	}
	switch strings.ToLower(strings.TrimSpace(f.Status)) {
	case "draft":
		q = q.Where("is_draft = ?", true)
	case "published":
		q = q.Where("is_draft = ?", false)
	}
	return applyVendorMarketFilter(q, f.Market)
}

func orderClause(p paging.Request) string {
	sort := strings.TrimSpace(p.Sort)
	switch sort {
	case "id", "name", "updated_at", "created_at":
	default:
		return "updated_at desc"
	}
	if strings.EqualFold(p.Order, "asc") {
		return sort + " asc"
	}
	return sort + " desc"
}

func (r *ProductRepository) ListAdmin(ctx context.Context, f AdminProductFilter, p paging.Request) ([]models.Product, int64, error) {
	var total int64
	if err := r.adminQuery(ctx, f).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []models.Product
	err := r.adminQuery(ctx, f).
		Preload("SolutionCategory").
		Order(orderClause(p)).
		Limit(p.Limit()).
		Offset(p.Offset()).
		Find(&list).Error
	return list, total, err
}

func (r *ProductRepository) ListAdminAll(ctx context.Context, f AdminProductFilter) ([]models.Product, error) {
	var list []models.Product
	err := r.adminQuery(ctx, f).
		Preload("SolutionCategory").
		Order("updated_at desc").
		Find(&list).Error
	return list, err
}

func (r *ProductRepository) Create(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductRepository) Reload(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Preload("SolutionCategory").First(p, p.ID).Error
}

func (r *ProductRepository) UpdateFields(ctx context.Context, id uint, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ProductRepository) UpdateDraft(ctx context.Context, id uint, isDraft bool) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Update("is_draft", isDraft).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", id).Delete(&models.SalesScript{}).Error; err != nil {
			return err
		}
		if err := tx.Where("product_id = ?", id).Delete(&models.Case{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Product{}, id).Error
	})
}

func (r *ProductRepository) ListScripts(ctx context.Context, productID uint) ([]models.SalesScript, error) {
	var list []models.SalesScript
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("id asc").Find(&list).Error
	return list, err
}

func (r *ProductRepository) CreateScript(ctx context.Context, s *models.SalesScript) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *ProductRepository) FindScript(ctx context.Context, id uint) (*models.SalesScript, error) {
	var s models.SalesScript
	if err := r.db.WithContext(ctx).First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ProductRepository) UpdateScript(ctx context.Context, s *models.SalesScript) error {
	return r.db.WithContext(ctx).Model(&models.SalesScript{}).Where("id = ?", s.ID).
		Updates(map[string]any{"scenario": s.Scenario, "content": s.Content}).Error
}

func (r *ProductRepository) DeleteScript(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.SalesScript{}, id).Error
}

func (r *ProductRepository) ListCases(ctx context.Context, productID uint) ([]models.Case, error) {
	var list []models.Case
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("id asc").Find(&list).Error
	return list, err
}

func (r *ProductRepository) CreateCase(ctx context.Context, c *models.Case) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *ProductRepository) FindCase(ctx context.Context, id uint) (*models.Case, error) {
	var c models.Case
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ProductRepository) UpdateCase(ctx context.Context, c *models.Case) error {
	return r.db.WithContext(ctx).Model(&models.Case{}).Where("id = ?", c.ID).
		Updates(map[string]any{
			"client_name":     c.ClientName,
			"pain_points":     c.PainPoints,
			"solution":        c.Solution,
			"value_delivered": c.ValueDelivered,
		}).Error
}

func (r *ProductRepository) DeleteCase(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Case{}, id).Error
}

func (r *ProductRepository) FindCategoryLabel(ctx context.Context, id uint) (string, error) {
	var cat models.SolutionCategory
	if err := r.db.WithContext(ctx).Select("id", "label").First(&cat, id).Error; err != nil {
		return "", err
	}
	return cat.Label, nil
}

func (r *ProductRepository) CountByCategory(ctx context.Context, categoryID uint) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("solution_category_id = ?", categoryID).Count(&n).Error
	return n, err
}
