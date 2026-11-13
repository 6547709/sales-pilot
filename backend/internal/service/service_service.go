package service

import (
	"context"
	"errors"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

var serviceKinds = map[string]bool{
	"v2v": true, "backup": true, "dr": true, "migration": true, "upgrade": true, "ops": true,
}

type ServiceService struct {
	repo ServiceRepo
}

func NewServiceService(repo ServiceRepo) *ServiceService {
	return &ServiceService{repo: repo}
}

func (s *ServiceService) List(ctx context.Context, kind string, activeOnly bool) ([]models.Service, error) {
	list, err := s.repo.List(ctx, kind, activeOnly)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return list, nil
}

func (s *ServiceService) GetBySlug(ctx context.Context, slug string) (*models.Service, error) {
	svc, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("service_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	if !svc.IsActive {
		return nil, apperr.NotFound("service_not_found", "未找到")
	}
	return svc, nil
}

func (s *ServiceService) GetAdmin(ctx context.Context, id uint) (*models.Service, error) {
	svc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("service_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	return svc, nil
}

func (s *ServiceService) CategoryIDs(ctx context.Context, id uint) ([]uint, error) {
	ids, err := s.repo.CategoryIDs(ctx, id)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return ids, nil
}

func validateService(svc *models.Service) error {
	if strings.TrimSpace(svc.Slug) == "" || strings.TrimSpace(svc.Name) == "" {
		return apperr.Invalid("invalid_body", "slug 与 name 必填")
	}
	if !serviceKinds[strings.TrimSpace(svc.Kind)] {
		return apperr.Invalid("invalid_service_kind", "服务类型不合法")
	}
	return nil
}

func (s *ServiceService) Create(ctx context.Context, svc *models.Service, categoryIDs []uint) (*models.Service, error) {
	if err := validateService(svc); err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, svc); err != nil {
		return nil, apperr.Conflict("service_slug_taken", "创建失败，slug 可能已存在")
	}
	if categoryIDs != nil {
		if err := s.repo.SetCategories(ctx, svc.ID, categoryIDs); err != nil {
			return nil, apperr.Internal(err)
		}
	}
	return svc, nil
}

func (s *ServiceService) Update(ctx context.Context, id uint, svc *models.Service, categoryIDs []uint) (*models.Service, error) {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("service_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	if err := validateService(svc); err != nil {
		return nil, err
	}
	svc.ID = id
	if err := s.repo.Save(ctx, svc); err != nil {
		return nil, apperr.Internal(err)
	}
	if categoryIDs != nil {
		if err := s.repo.SetCategories(ctx, id, categoryIDs); err != nil {
			return nil, apperr.Internal(err)
		}
	}
	return svc, nil
}

func (s *ServiceService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}
