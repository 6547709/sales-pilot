package repository

import (
	"context"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Auth(ctx context.Context) (*models.AuthSettings, error) {
	var st models.AuthSettings
	if err := r.db.WithContext(ctx).FirstOrCreate(&st, models.AuthSettings{ID: 1}).Error; err != nil {
		return nil, err
	}
	return &st, nil
}

func (r *SettingsRepository) SaveAuth(ctx context.Context, st *models.AuthSettings) error {
	return r.db.WithContext(ctx).Save(st).Error
}

func (r *SettingsRepository) System(ctx context.Context) (*models.SystemSettings, error) {
	var st models.SystemSettings
	if err := r.db.WithContext(ctx).FirstOrCreate(&st, models.SystemSettings{ID: 1}).Error; err != nil {
		return nil, err
	}
	return &st, nil
}

func (r *SettingsRepository) SaveSystem(ctx context.Context, st *models.SystemSettings) error {
	return r.db.WithContext(ctx).Save(st).Error
}
