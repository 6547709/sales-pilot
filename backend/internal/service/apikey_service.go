package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

type APIKeyService struct {
	repo  APIKeyRepo
	users UserRepo
}

func NewAPIKeyService(repo APIKeyRepo, users UserRepo) *APIKeyService {
	return &APIKeyService{repo: repo, users: users}
}

func (s *APIKeyService) List(ctx context.Context) ([]models.AdminAPIKey, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return list, nil
}

func (s *APIKeyService) Create(ctx context.Context, userID uint, name string, expiresAt *time.Time) (string, *models.AdminAPIKey, error) {
	if name == "" {
		return "", nil, apperr.Invalid("invalid_name", "请提供 name")
	}
	if expiresAt != nil && !expiresAt.After(time.Now()) {
		return "", nil, apperr.Invalid("invalid_expiry", "过期时间须晚于当前时间")
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", nil, apperr.Internal(err)
	}
	fullKey := "sp_" + hex.EncodeToString(raw)
	prefix := fullKey
	if len(prefix) > 20 {
		prefix = prefix[:20] + "…"
	}
	k := &models.AdminAPIKey{
		UserID: userID, Name: name, KeyHash: auth.HashAPIKey(fullKey), Prefix: prefix,
		ExpiresAt: expiresAt, IsActive: true,
	}
	if err := s.repo.Create(ctx, k); err != nil {
		return "", nil, apperr.Internal(err)
	}
	return fullKey, k, nil
}

func (s *APIKeyService) Revoke(ctx context.Context, id uint) error {
	if err := s.repo.Revoke(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *APIKeyService) Authenticate(ctx context.Context, raw string) (*models.User, error) {
	k, err := s.repo.FindByHash(ctx, auth.HashAPIKey(raw))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.Unauthorized("api_key_invalid", "API Key 无效或已停用")
		}
		return nil, apperr.Internal(err)
	}
	if !k.IsActive {
		return nil, apperr.Unauthorized("api_key_invalid", "API Key 无效或已停用")
	}
	if k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt) {
		return nil, apperr.Unauthorized("api_key_expired", "API Key 已过期")
	}
	u, err := s.users.FindByID(ctx, k.UserID)
	if err != nil || !u.IsActive || u.Role != "admin" {
		return nil, apperr.Forbidden("api_key_forbidden", "密钥关联用户不可用或非管理员")
	}
	_ = s.repo.TouchLastUsed(ctx, k.ID)
	return u, nil
}
