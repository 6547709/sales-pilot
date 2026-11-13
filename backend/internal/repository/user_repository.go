package repository

import (
	"context"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) ListLocal(ctx context.Context) ([]models.User, error) {
	var list []models.User
	err := r.db.WithContext(ctx).Where("auth_provider = ?", "local").Order("id asc").Find(&list).Error
	return list, err
}

func (r *UserRepository) ListAll(ctx context.Context) ([]models.User, error) {
	var list []models.User
	err := r.db.WithContext(ctx).Order("id asc").Find(&list).Error
	return list, err
}

func (r *UserRepository) FindByUsernameAndProvider(ctx context.Context, username, provider string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).Where("username = ? AND auth_provider = ?", username, provider).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepository) Save(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("last_login", time.Now()).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *UserRepository) CountActiveLocalAdmins(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("auth_provider = ? AND role = ? AND is_active = ?", "local", "admin", true).
		Count(&n).Error
	return n, err
}
