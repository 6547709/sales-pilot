package repository

import (
	"context"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

type APIKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

func (r *APIKeyRepository) List(ctx context.Context) ([]models.AdminAPIKey, error) {
	var list []models.AdminAPIKey
	err := r.db.WithContext(ctx).Order("id desc").Find(&list).Error
	return list, err
}

func (r *APIKeyRepository) FindByHash(ctx context.Context, hash string) (*models.AdminAPIKey, error) {
	var k models.AdminAPIKey
	if err := r.db.WithContext(ctx).Where("key_hash = ?", hash).First(&k).Error; err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *APIKeyRepository) Create(ctx context.Context, k *models.AdminAPIKey) error {
	return r.db.WithContext(ctx).Create(k).Error
}

func (r *APIKeyRepository) Revoke(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.AdminAPIKey{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *APIKeyRepository) TouchLastUsed(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.AdminAPIKey{}).Where("id = ?", id).Update("last_used_at", time.Now()).Error
}
