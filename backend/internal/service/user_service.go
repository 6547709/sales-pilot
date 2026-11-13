package service

import (
	"context"
	"errors"

	"github.com/6547709/sales-pilot/backend/internal/paging"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/repository"
	"gorm.io/gorm"
)

type CreateLocalUserInput struct {
	Username string
	Password string
	Email    string
	Role     string
}

type UpdateLocalUserInput struct {
	Email    *string
	Role     *string
	IsActive *bool
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ListLocal(ctx context.Context) ([]models.User, error) {
	list, err := s.repo.ListLocal(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return list, nil
}

func (s *UserService) CreateLocal(ctx context.Context, in CreateLocalUserInput) (*models.User, error) {
	if in.Username == "" || len(in.Password) < 6 {
		return nil, apperr.Invalid("invalid_body", "用户名、密码必填，密码至少 6 位")
	}
	role := in.Role
	if role != "admin" && role != "user" {
		role = "user"
	}
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	email := in.Email
	if email == "" {
		email = in.Username + "@local"
	}
	u := &models.User{
		Username: in.Username, Email: email, PasswordHash: hash,
		AuthProvider: "local", Role: role, IsActive: true,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, apperr.Conflict("username_taken", "创建失败，用户名可能已存在")
	}
	return u, nil
}

func (s *UserService) UpdateLocal(ctx context.Context, id uint, in UpdateLocalUserInput) (*models.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("user_not_found", "未找到")
		}
		return nil, apperr.Internal(err)
	}
	if in.Email != nil {
		u.Email = *in.Email
	}
	if in.Role != nil {
		if *in.Role != "admin" && *in.Role != "user" {
			return nil, apperr.Invalid("invalid_role", "角色只能是 admin 或 user")
		}
		if u.AuthProvider == "local" && u.Role == "admin" && u.IsActive && *in.Role != "admin" {
			n, err := s.repo.CountActiveLocalAdmins(ctx)
			if err != nil {
				return nil, apperr.Internal(err)
			}
			if n <= 1 {
				return nil, apperr.Conflict("last_admin", "至少保留一名可用管理员")
			}
		}
		u.Role = *in.Role
	}
	if in.IsActive != nil {
		if !*in.IsActive && u.AuthProvider == "local" && u.Role == "admin" && u.IsActive {
			n, err := s.repo.CountActiveLocalAdmins(ctx)
			if err != nil {
				return nil, apperr.Internal(err)
			}
			if n <= 1 {
				return nil, apperr.Conflict("last_admin", "至少保留一名可用管理员")
			}
		}
		u.IsActive = *in.IsActive
	}
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, apperr.Internal(err)
	}
	return u, nil
}

func (s *UserService) ResetPassword(ctx context.Context, id uint, password string) error {
	if len(password) < 6 {
		return apperr.Invalid("invalid_password", "密码至少 6 位")
	}
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.NotFound("user_not_found", "未找到")
		}
		return apperr.Internal(err)
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return apperr.Internal(err)
	}
	u.PasswordHash = hash
	if err := s.repo.Save(ctx, u); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *UserService) DeleteLocal(ctx context.Context, id uint) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.NotFound("user_not_found", "未找到")
		}
		return apperr.Internal(err)
	}
	if u.AuthProvider == "local" && u.Role == "admin" && u.IsActive {
		n, err := s.repo.CountActiveLocalAdmins(ctx)
		if err != nil {
			return apperr.Internal(err)
		}
		if n <= 1 {
			return apperr.Conflict("last_admin", "至少保留一名可用管理员")
		}
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *UserService) ListLocalPage(ctx context.Context, f repository.LocalUserFilter, p paging.Request) ([]models.User, int64, error) {
	list, total, err := s.repo.ListLocalPage(ctx, f, p)
	if err != nil {
		return nil, 0, apperr.Internal(err)
	}
	return list, total, nil
}
