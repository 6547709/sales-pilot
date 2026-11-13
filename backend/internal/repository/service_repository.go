package repository

import (
	"context"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func dedupeIDs(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func (r *ServiceRepository) List(ctx context.Context, kind string, activeOnly bool) ([]models.Service, error) {
	q := r.db.WithContext(ctx).Model(&models.Service{})
	if kind != "" {
		q = q.Where("kind = ?", kind)
	}
	if activeOnly {
		q = q.Where("is_active = ?", true)
	}
	var list []models.Service
	err := q.Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ServiceRepository) FindByID(ctx context.Context, id uint) (*models.Service, error) {
	var s models.Service
	if err := r.db.WithContext(ctx).First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ServiceRepository) FindBySlug(ctx context.Context, slug string) (*models.Service, error) {
	var s models.Service
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ServiceRepository) Create(ctx context.Context, s *models.Service) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *ServiceRepository) Save(ctx context.Context, s *models.Service) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *ServiceRepository) UpdateFields(ctx context.Context, id uint, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Service{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ServiceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("service_id = ?", id).Delete(&models.ServiceCategory{}).Error; err != nil {
			return err
		}
		if err := tx.Where("service_id = ?", id).Delete(&models.ScenarioService{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Service{}, id).Error
	})
}

func (r *ServiceRepository) SetCategories(ctx context.Context, serviceID uint, categoryIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("service_id = ?", serviceID).Delete(&models.ServiceCategory{}).Error; err != nil {
			return err
		}
		for _, cid := range dedupeIDs(categoryIDs) {
			if err := tx.Create(&models.ServiceCategory{ServiceID: serviceID, CategoryID: cid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ServiceRepository) CategoryIDs(ctx context.Context, serviceID uint) ([]uint, error) {
	var ids []uint
	err := r.db.WithContext(ctx).Model(&models.ServiceCategory{}).
		Where("service_id = ?", serviceID).Pluck("category_id", &ids).Error
	return ids, err
}

func (r *ServiceRepository) ListByCategoryID(ctx context.Context, categoryID uint) ([]models.Service, error) {
	var ids []uint
	if err := r.db.WithContext(ctx).Model(&models.ServiceCategory{}).
		Where("category_id = ?", categoryID).Pluck("service_id", &ids).Error; err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []models.Service{}, nil
	}
	var list []models.Service
	err := r.db.WithContext(ctx).
		Where("id IN ? AND is_active = ?", ids, true).
		Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *ServiceRepository) ListByCategoryIDs(ctx context.Context, ids []uint) (map[uint][]models.Service, error) {
	out := make(map[uint][]models.Service, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	type row struct {
		CategoryID uint
		ServiceID  uint
	}
	var links []row
	if err := r.db.WithContext(ctx).Model(&models.ServiceCategory{}).
		Where("category_id IN ?", ids).Find(&links).Error; err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return out, nil
	}
	svcIDs := make([]uint, 0, len(links))
	for _, l := range links {
		svcIDs = append(svcIDs, l.ServiceID)
	}
	var services []models.Service
	if err := r.db.WithContext(ctx).
		Where("id IN ? AND is_active = ?", dedupeIDs(svcIDs), true).
		Order("sort_order ASC, id ASC").Find(&services).Error; err != nil {
		return nil, err
	}
	byID := make(map[uint]models.Service, len(services))
	for _, s := range services {
		byID[s.ID] = s
	}
	for _, l := range links {
		if s, ok := byID[l.ServiceID]; ok {
			out[l.CategoryID] = append(out[l.CategoryID], s)
		}
	}
	return out, nil
}

func (r *ServiceRepository) DeleteCategoryLinks(ctx context.Context, categoryID uint) error {
	return r.db.WithContext(ctx).Where("category_id = ?", categoryID).Delete(&models.ServiceCategory{}).Error
}

func (r *ServiceRepository) ListByIDs(ctx context.Context, ids []uint) ([]models.Service, error) {
	if len(ids) == 0 {
		return []models.Service{}, nil
	}
	var list []models.Service
	err := r.db.WithContext(ctx).Where("id IN ?", dedupeIDs(ids)).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}
