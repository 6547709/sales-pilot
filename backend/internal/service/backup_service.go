package service

import (
	"context"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/cache"
	"github.com/6547709/sales-pilot/backend/internal/models"
)

// BackupConfirmToken 与前端约定：仅在确认框填入此令牌时才执行全量覆盖恢复
const BackupConfirmToken = "RESTORE_FULL_V1"

type BackupService struct {
	repo  BackupRepo
	cache *cache.SearchCache
}

func NewBackupService(repo BackupRepo, c *cache.SearchCache) *BackupService {
	return &BackupService{repo: repo, cache: c}
}

func (s *BackupService) Export(ctx context.Context) (*models.BackupPayload, error) {
	payload, err := s.repo.Dump(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return payload, nil
}

func (s *BackupService) Import(ctx context.Context, payload *models.BackupPayload, confirm string) error {
	if confirm != BackupConfirmToken {
		return apperr.Invalid("invalid_confirm", "确认令牌不正确")
	}
	if payload.Version != 1 {
		return apperr.Invalid("unsupported_backup_version", "不支持的备份版本")
	}
	admins := 0
	for _, u := range payload.Users {
		if u.AuthProvider == "local" && u.Role == "admin" && u.IsActive {
			admins++
		}
	}
	if admins == 0 {
		return apperr.Invalid("no_admin", "备份中必须至少包含一名可用的本地管理员账号")
	}
	if err := s.repo.Restore(ctx, payload); err != nil {
		return apperr.Internal(err)
	}
	if s.cache != nil {
		s.cache.Clear()
	}
	return nil
}
