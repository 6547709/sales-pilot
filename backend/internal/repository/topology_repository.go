package repository

import (
	"context"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

type TopologyRepository struct {
	db *gorm.DB
}

func NewTopologyRepository(db *gorm.DB) *TopologyRepository {
	return &TopologyRepository{db: db}
}

// TopologySnapshot 首页拓扑所需的全部只读数据
type TopologySnapshot struct {
	Layers   []models.TopologyLayer
	Security []models.SolutionCategory
	Ops      []models.SolutionCategory
	Central  map[uint][]models.SolutionCategory
}

func vendorPreload(db *gorm.DB) *gorm.DB {
	return db.Order("sort_order ASC, id ASC")
}

func (r *TopologyRepository) Snapshot(ctx context.Context) (*TopologySnapshot, error) {
	db := r.db.WithContext(ctx)
	var layers []models.TopologyLayer
	if err := db.Order("sort_order ASC, level DESC").Find(&layers).Error; err != nil {
		return nil, err
	}
	var sec []models.SolutionCategory
	if err := db.Where("column_type = ? AND is_active = ?", "security", true).
		Preload("Vendors", vendorPreload).
		Order("sort_order ASC, id ASC").Find(&sec).Error; err != nil {
		return nil, err
	}
	var ops []models.SolutionCategory
	if err := db.Where("column_type = ? AND is_active = ?", "ops", true).
		Preload("Vendors", vendorPreload).
		Order("sort_order ASC, id ASC").Find(&ops).Error; err != nil {
		return nil, err
	}
	central := make(map[uint][]models.SolutionCategory, len(layers))
	for _, l := range layers {
		var cats []models.SolutionCategory
		if err := db.Where("column_type = ? AND layer_id = ? AND is_active = ?", "central", l.ID, true).
			Preload("Vendors", vendorPreload).
			Order("sort_order ASC, id ASC").Find(&cats).Error; err != nil {
			return nil, err
		}
		central[l.ID] = cats
	}
	return &TopologySnapshot{Layers: layers, Security: sec, Ops: ops, Central: central}, nil
}

func (r *TopologyRepository) Categories(ctx context.Context, activeOnly bool) ([]models.SolutionCategory, error) {
	q := r.db.WithContext(ctx).Model(&models.SolutionCategory{}).
		Order("column_type ASC, sort_order ASC, id ASC")
	if activeOnly {
		q = q.Where("is_active = ?", true)
	}
	var list []models.SolutionCategory
	err := q.Find(&list).Error
	return list, err
}

func (r *TopologyRepository) FindCategory(ctx context.Context, id uint) (*models.SolutionCategory, error) {
	var c models.SolutionCategory
	if err := r.db.WithContext(ctx).Preload("Vendors").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *TopologyRepository) CreateCategory(ctx context.Context, c *models.SolutionCategory) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *TopologyRepository) UpdateCategoryFields(ctx context.Context, id uint, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.SolutionCategory{}).Where("id = ?", id).Updates(fields).Error
}

func (r *TopologyRepository) DeleteCategory(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("category_id = ?", id).Delete(&models.SolutionVendor{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.SolutionCategory{}, id).Error
	})
}

func (r *TopologyRepository) CountCategoriesInLayer(ctx context.Context, layerID uint) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.SolutionCategory{}).Where("layer_id = ?", layerID).Count(&n).Error
	return n, err
}

func (r *TopologyRepository) Layers(ctx context.Context) ([]models.TopologyLayer, error) {
	var list []models.TopologyLayer
	err := r.db.WithContext(ctx).Order("sort_order ASC, level DESC").Find(&list).Error
	return list, err
}

func (r *TopologyRepository) FindLayer(ctx context.Context, id uint) (*models.TopologyLayer, error) {
	var l models.TopologyLayer
	if err := r.db.WithContext(ctx).First(&l, id).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *TopologyRepository) CreateLayer(ctx context.Context, l *models.TopologyLayer) error {
	return r.db.WithContext(ctx).Create(l).Error
}

func (r *TopologyRepository) UpdateLayerFields(ctx context.Context, id uint, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.TopologyLayer{}).Where("id = ?", id).Updates(fields).Error
}

func (r *TopologyRepository) DeleteLayer(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TopologyLayer{}, id).Error
}

func (r *TopologyRepository) CreateVendor(ctx context.Context, v *models.SolutionVendor) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *TopologyRepository) FindVendor(ctx context.Context, id uint) (*models.SolutionVendor, error) {
	var v models.SolutionVendor
	if err := r.db.WithContext(ctx).First(&v, id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *TopologyRepository) UpdateVendorFields(ctx context.Context, id uint, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.SolutionVendor{}).Where("id = ?", id).Updates(fields).Error
}

func (r *TopologyRepository) DeleteVendor(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.SolutionVendor{}, id).Error
}

func (r *TopologyRepository) CountProductsInCategory(ctx context.Context, categoryID uint) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("solution_category_id = ?", categoryID).Count(&n).Error
	return n, err
}
