package service

import (
	"context"
	"errors"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/config"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

type LoginOptions struct {
	LocalEnabled bool `json:"local_enabled"`
	LdapEnabled  bool `json:"ldap_enabled"`
	OidcEnabled  bool `json:"oidc_enabled"`
}

type AuthService struct {
	cfg      config.Config
	users    UserRepo
	settings SettingsRepo
}

func NewAuthService(cfg config.Config, users UserRepo, settings SettingsRepo) *AuthService {
	return &AuthService{cfg: cfg, users: users, settings: settings}
}

func (s *AuthService) LoginLocal(ctx context.Context, username, password string) (*models.User, string, error) {
	st, err := s.settings.Auth(ctx)
	if err != nil {
		return nil, "", apperr.Internal(err)
	}
	if !st.LocalLoginEnabled {
		return nil, "", apperr.Forbidden("local_login_disabled", "本地账号登录已关闭，请使用企业 SSO 或 LDAP")
	}
	u, err := s.users.FindByUsernameAndProvider(ctx, username, "local")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", apperr.Unauthorized("invalid_credentials", "用户名或密码错误")
		}
		return nil, "", apperr.Internal(err)
	}
	if !u.IsActive || !auth.CheckPassword(u.PasswordHash, password) {
		return nil, "", apperr.Unauthorized("invalid_credentials", "用户名或密码错误")
	}
	_ = s.users.UpdateLastLogin(ctx, u.ID)
	token, err := auth.IssueToken(s.cfg.JWTSecret, u.ID, u.Username, u.Role)
	if err != nil {
		return nil, "", apperr.Internal(err)
	}
	return u, token, nil
}

func (s *AuthService) LoginLDAP(ctx context.Context, username, password string) (*models.User, string, error) {
	st, err := s.settings.Auth(ctx)
	if err != nil {
		return nil, "", apperr.Internal(err)
	}
	if !st.LdapEnabled {
		return nil, "", apperr.Invalid("ldap_disabled", "LDAP 未启用")
	}
	lc, err := auth.ParseLdapConfig([]byte(st.LdapConfig))
	if err != nil {
		return nil, "", apperr.Invalid("ldap_config_invalid", "LDAP 配置无效")
	}
	if err := auth.TryLDAP(lc, username, password); err != nil {
		return nil, "", apperr.Unauthorized("ldap_auth_failed", "LDAP 认证失败")
	}
	u, err := s.users.FindByUsernameAndProvider(ctx, username, "ldap")
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", apperr.Internal(err)
		}
		u = &models.User{
			Username: username, Email: username + "@ldap.local",
			AuthProvider: "ldap", Role: "user", IsActive: true,
		}
		if err := s.users.Create(ctx, u); err != nil {
			return nil, "", apperr.Internal(err)
		}
	}
	if !u.IsActive {
		return nil, "", apperr.Forbidden("user_disabled", "账号已禁用")
	}
	_ = s.users.UpdateLastLogin(ctx, u.ID)
	token, err := auth.IssueToken(s.cfg.JWTSecret, u.ID, u.Username, u.Role)
	if err != nil {
		return nil, "", apperr.Internal(err)
	}
	return u, token, nil
}

func (s *AuthService) LoginOptions(ctx context.Context) (LoginOptions, error) {
	st, err := s.settings.Auth(ctx)
	if err != nil {
		return LoginOptions{LocalEnabled: true}, nil
	}
	return LoginOptions{
		LocalEnabled: st.LocalLoginEnabled,
		LdapEnabled:  st.LdapEnabled,
		OidcEnabled:  st.OidcEnabled,
	}, nil
}

func (s *AuthService) Me(ctx context.Context, userID uint) (*models.User, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("user_not_found", "用户不存在")
		}
		return nil, apperr.Internal(err)
	}
	return u, nil
}

// AuthSettings 供 OIDC 流程读取配置
func (s *AuthService) AuthSettings(ctx context.Context) (*models.AuthSettings, error) {
	st, err := s.settings.Auth(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return st, nil
}

// LoginOIDC 以 OIDC 身份查找或创建用户并签发令牌
func (s *AuthService) LoginOIDC(ctx context.Context, username, email string) (string, error) {
	u, err := s.users.FindByUsernameAndProvider(ctx, username, "oidc")
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperr.Internal(err)
		}
		u = &models.User{
			Username: username, Email: email,
			AuthProvider: "oidc", Role: "user", IsActive: true,
		}
		if err := s.users.Create(ctx, u); err != nil {
			return "", apperr.Internal(err)
		}
	}
	if !u.IsActive {
		return "", apperr.Forbidden("user_disabled", "账号已禁用")
	}
	_ = s.users.UpdateLastLogin(ctx, u.ID)
	token, err := auth.IssueToken(s.cfg.JWTSecret, u.ID, u.Username, u.Role)
	if err != nil {
		return "", apperr.Internal(err)
	}
	return token, nil
}
