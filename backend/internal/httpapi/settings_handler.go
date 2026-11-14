package httpapi

import (
	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type authSettingsPatch struct {
	LocalLoginEnabled *bool           `json:"local_login_enabled"`
	LdapEnabled       *bool           `json:"ldap_enabled"`
	LdapConfig        *datatypes.JSON `json:"ldap_config"`
	OidcEnabled       *bool           `json:"oidc_enabled"`
	OidcConfig        *datatypes.JSON `json:"oidc_config"`
}

func (h *Handlers) getAuthSettings(c *gin.Context) {
	st, err := h.Settings.Auth(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, st)
}

func (h *Handlers) patchAuthSettings(c *gin.Context) {
	var body authSettingsPatch
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	st, err := h.Settings.PatchAuth(c.Request.Context(), service.AuthSettingsPatch{
		LocalLoginEnabled: body.LocalLoginEnabled,
		LdapEnabled:       body.LdapEnabled,
		LdapConfig:        body.LdapConfig,
		OidcEnabled:       body.OidcEnabled,
		OidcConfig:        body.OidcConfig,
	})
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, st)
}

type systemSettingsPatch struct {
	FooterText   *string `json:"footer_text"`
	FilingNumber *string `json:"filing_number"`
}

func (h *Handlers) getSystemSettings(c *gin.Context) {
	st, err := h.Settings.System(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, st)
}

func (h *Handlers) patchSystemSettings(c *gin.Context) {
	var body systemSettingsPatch
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	st, err := h.Settings.PatchSystem(c.Request.Context(), service.SystemSettingsPatch{
		FooterText:   body.FooterText,
		FilingNumber: body.FilingNumber,
	})
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, st)
}

func (h *Handlers) adminStats(c *gin.Context) {
	data, err := h.Stats.Overview(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, data)
}
