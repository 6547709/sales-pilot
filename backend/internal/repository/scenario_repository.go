package repository

import (
	"context"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

type ScenarioRepository struct {
	db *gorm.DB
}

func NewScenarioRepository(db *gorm.DB) *ScenarioRepository {
	return &ScenarioRepository{db: db}
}

func (r *ScenarioRepository) Groups(ctx context.Context) ([]models.ScenarioGroup, error) {
	var list []models.ScenarioGroup
	err := r.db.WithContext(ctx).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ScenarioRepository) Scenarios(ctx context.Context) ([]models.Scenario, error) {
	var list []models.Scenario
	err := r.db.WithContext(ctx).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ScenarioRepository) FindBySlug(ctx context.Context, slug string) (*models.Scenario, error) {
	var s models.Scenario
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ScenarioRepository) FindByID(ctx context.Context, id uint) (*models.Scenario, error) {
	var s models.Scenario
	if err := r.db.WithContext(ctx).First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ScenarioRepository) FindGroup(ctx context.Context, id uint) (*models.ScenarioGroup, error) {
	var g models.ScenarioGroup
	if err := r.db.WithContext(ctx).First(&g, id).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *ScenarioRepository) CreateGroup(ctx context.Context, g *models.ScenarioGroup) error {
	return r.db.WithContext(ctx).Create(g).Error
}

func (r *ScenarioRepository) SaveGroup(ctx context.Context, g *models.ScenarioGroup) error {
	return r.db.WithContext(ctx).Save(g).Error
}

func (r *ScenarioRepository) UpdateGroupFields(ctx context.Context, id uint, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.ScenarioGroup{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ScenarioRepository) UpdateFields(ctx context.Context, id uint, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Scenario{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ScenarioRepository) DeleteGroup(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ScenarioGroup{}, id).Error
}

func (r *ScenarioRepository) CountGroups(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.ScenarioGroup{}).Count(&n).Error
	return n, err
}

func (r *ScenarioRepository) CountScenariosInGroup(ctx context.Context, groupID uint) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Scenario{}).Where("group_id = ?", groupID).Count(&n).Error
	return n, err
}

func (r *ScenarioRepository) Create(ctx context.Context, s *models.Scenario) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *ScenarioRepository) Save(ctx context.Context, s *models.Scenario) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *ScenarioRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("scenario_id = ?", id).Delete(&models.ScenarioProduct{}).Error; err != nil {
			return err
		}
		if err := tx.Where("scenario_id = ?", id).Delete(&models.ScenarioService{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Scenario{}, id).Error
	})
}

func (r *ScenarioRepository) ProductIDs(ctx context.Context, scenarioID uint) ([]uint, error) {
	var ids []uint
	err := r.db.WithContext(ctx).Model(&models.ScenarioProduct{}).
		Where("scenario_id = ?", scenarioID).Pluck("product_id", &ids).Error
	return ids, err
}

func (r *ScenarioRepository) ServiceIDs(ctx context.Context, scenarioID uint) ([]uint, error) {
	var ids []uint
	err := r.db.WithContext(ctx).Model(&models.ScenarioService{}).
		Where("scenario_id = ?", scenarioID).Pluck("service_id", &ids).Error
	return ids, err
}

func (r *ScenarioRepository) SetLinks(ctx context.Context, scenarioID uint, productIDs, serviceIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("scenario_id = ?", scenarioID).Delete(&models.ScenarioProduct{}).Error; err != nil {
			return err
		}
		if err := tx.Where("scenario_id = ?", scenarioID).Delete(&models.ScenarioService{}).Error; err != nil {
			return err
		}
		for _, pid := range dedupeIDs(productIDs) {
			if err := tx.Create(&models.ScenarioProduct{ScenarioID: scenarioID, ProductID: pid}).Error; err != nil {
				return err
			}
		}
		for _, sid := range dedupeIDs(serviceIDs) {
			if err := tx.Create(&models.ScenarioService{ScenarioID: scenarioID, ServiceID: sid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

type countRow struct {
	ScenarioID uint
	N          int64
}

func (r *ScenarioRepository) CountProducts(ctx context.Context) (map[uint]int64, error) {
	var rows []countRow
	err := r.db.WithContext(ctx).Model(&models.ScenarioProduct{}).
		Select("scenario_id, count(*) as n").Group("scenario_id").Scan(&rows).Error
	out := make(map[uint]int64, len(rows))
	for _, row := range rows {
		out[row.ScenarioID] = row.N
	}
	return out, err
}

func (r *ScenarioRepository) CountServices(ctx context.Context) (map[uint]int64, error) {
	var rows []countRow
	err := r.db.WithContext(ctx).Model(&models.ScenarioService{}).
		Select("scenario_id, count(*) as n").Group("scenario_id").Scan(&rows).Error
	out := make(map[uint]int64, len(rows))
	for _, row := range rows {
		out[row.ScenarioID] = row.N
	}
	return out, err
}

func (r *ScenarioRepository) ListProductsByIDs(ctx context.Context, ids []uint) ([]models.Product, error) {
	if len(ids) == 0 {
		return []models.Product{}, nil
	}
	var list []models.Product
	err := r.db.WithContext(ctx).Where("id IN ? AND is_draft = ?", dedupeIDs(ids), false).
		Order("updated_at desc").Find(&list).Error
	return list, err
}

func (r *ScenarioRepository) ListServicesByIDs(ctx context.Context, ids []uint) ([]models.Service, error) {
	if len(ids) == 0 {
		return []models.Service{}, nil
	}
	var list []models.Service
	err := r.db.WithContext(ctx).Where("id IN ? AND is_active = ?", dedupeIDs(ids), true).
		Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ScenarioRepository) ListPublicCasesByProductIDs(ctx context.Context, ids []uint) ([]models.Case, error) {
	if len(ids) == 0 {
		return []models.Case{}, nil
	}
	var list []models.Case
	err := r.db.WithContext(ctx).
		Where("product_id IN ? AND is_public = ?", dedupeIDs(ids), true).
		Order("id asc").Find(&list).Error
	return list, err
}

func (r *ScenarioRepository) ListCasesByProductIDs(ctx context.Context, ids []uint) ([]models.Case, error) {
	if len(ids) == 0 {
		return []models.Case{}, nil
	}
	var list []models.Case
	err := r.db.WithContext(ctx).Where("product_id IN ?", dedupeIDs(ids)).Order("id asc").Find(&list).Error
	return list, err
}

func (r *ScenarioRepository) ListProductsByIDsAll(ctx context.Context, ids []uint) ([]models.Product, error) {
	if len(ids) == 0 {
		return []models.Product{}, nil
	}
	var list []models.Product
	err := r.db.WithContext(ctx).Where("id IN ?", dedupeIDs(ids)).Order("updated_at desc").Find(&list).Error
	return list, err
}

func (r *ScenarioRepository) ListServicesByIDsAll(ctx context.Context, ids []uint) ([]models.Service, error) {
	if len(ids) == 0 {
		return []models.Service{}, nil
	}
	var list []models.Service
	err := r.db.WithContext(ctx).Where("id IN ?", dedupeIDs(ids)).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ScenarioRepository) DeleteLinksForProduct(ctx context.Context, productID uint) error {
	return r.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&models.ScenarioProduct{}).Error
}
