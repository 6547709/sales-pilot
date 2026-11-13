package service

import (
	"context"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/datatypes"
)

type AuthSettingsPatch struct {
	LocalLoginEnabled *bool
	LdapEnabled       *bool
	LdapConfig        *datatypes.JSON
	OidcEnabled       *bool
	OidcConfig        *datatypes.JSON
}

type SystemSettingsPatch struct {
	FooterText   *string
	FilingNumber *string
}

type SettingsService struct {
	repo SettingsRepo
}

func NewSettingsService(repo SettingsRepo) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) Auth(ctx context.Context) (*models.AuthSettings, error) {
	st, err := s.repo.Auth(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return st, nil
}

func (s *SettingsService) PatchAuth(ctx context.Context, in AuthSettingsPatch) (*models.AuthSettings, error) {
	st, err := s.repo.Auth(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if in.LocalLoginEnabled != nil {
		st.LocalLoginEnabled = *in.LocalLoginEnabled
	}
	if in.LdapEnabled != nil {
		st.LdapEnabled = *in.LdapEnabled
	}
	if in.LdapConfig != nil {
		st.LdapConfig = *in.LdapConfig
	}
	if in.OidcEnabled != nil {
		st.OidcEnabled = *in.OidcEnabled
	}
	if in.OidcConfig != nil {
		st.OidcConfig = *in.OidcConfig
	}
	if err := s.repo.SaveAuth(ctx, st); err != nil {
		return nil, apperr.Internal(err)
	}
	return st, nil
}

func (s *SettingsService) System(ctx context.Context) (*models.SystemSettings, error) {
	st, err := s.repo.System(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return st, nil
}

func (s *SettingsService) PatchSystem(ctx context.Context, in SystemSettingsPatch) (*models.SystemSettings, error) {
	st, err := s.repo.System(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if in.FooterText != nil {
		st.FooterText = *in.FooterText
	}
	if in.FilingNumber != nil {
		st.FilingNumber = *in.FilingNumber
	}
	if err := s.repo.SaveSystem(ctx, st); err != nil {
		return nil, apperr.Internal(err)
	}
	return st, nil
}

func (s *SettingsService) Config(ctx context.Context) (string, string) {
	st, err := s.repo.System(ctx)
	if err != nil {
		return "", ""
	}
	return st.FooterText, st.FilingNumber
}
