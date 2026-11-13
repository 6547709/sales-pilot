package httpapi

import (
	"context"
	"net/http"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/httpapi/middleware"
	"github.com/gin-gonic/gin"
)

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) postLocalLogin(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	u, token, err := h.Auth.LoginLocal(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"access_token": token, "token_type": "Bearer", "user": u})
}

func (h *Handlers) postLDAPLogin(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		Fail(c, apperr.Invalid("invalid_body", "参数错误"))
		return
	}
	u, token, err := h.Auth.LoginLDAP(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"access_token": token, "token_type": "Bearer", "user": u})
}

func (h *Handlers) getLoginOptions(c *gin.Context) {
	opts, err := h.Auth.LoginOptions(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, opts)
}

func (h *Handlers) getMe(c *gin.Context) {
	cl, ok := middleware.Claims(c)
	if !ok {
		Fail(c, apperr.Unauthorized("unauthorized", "未授权"))
		return
	}
	u, err := h.Auth.Me(c.Request.Context(), cl.UserID)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, u)
}

func (h *Handlers) getOidcStart(c *gin.Context) {
	st, err := h.Auth.AuthSettings(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	if !st.OidcEnabled {
		Fail(c, apperr.Invalid("oidc_disabled", "OIDC 未启用"))
		return
	}
	oc, err := auth.ParseOidcConfig([]byte(st.OidcConfig))
	if err != nil {
		Fail(c, apperr.Invalid("oidc_config_invalid", "OIDC 配置无效"))
		return
	}
	ctx := context.Background()
	oauthCfg, _, err := auth.OidcEndpoints(ctx, oc)
	if err != nil {
		Fail(c, apperr.Internal(err))
		return
	}
	state, err := auth.RandomState()
	if err != nil {
		Fail(c, apperr.Internal(err))
		return
	}
	h.OidcStates.Put(state)
	c.Redirect(http.StatusFound, oauthCfg.AuthCodeURL(state))
}

func (h *Handlers) getOidcCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	if !h.OidcStates.Take(state) {
		c.String(http.StatusBadRequest, "state 无效或过期")
		return
	}
	st, err := h.Auth.AuthSettings(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "配置读取失败")
		return
	}
	oc, err := auth.ParseOidcConfig([]byte(st.OidcConfig))
	if err != nil {
		c.String(http.StatusInternalServerError, "OIDC 配置无效")
		return
	}
	ctx := context.Background()
	oauthCfg, provider, err := auth.OidcEndpoints(ctx, oc)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	ui, err := auth.ExchangeAndVerify(ctx, oauthCfg, provider, code)
	if err != nil {
		c.String(http.StatusUnauthorized, "OIDC 校验失败: "+err.Error())
		return
	}
	username := ui.PreferredName
	if username == "" {
		username = ui.Sub
	}
	token, err := h.Auth.LoginOIDC(c.Request.Context(), username, ui.Email)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "登录失败"
		if e, ok := apperr.From(err); ok {
			msg = e.Message
			if e.Kind == apperr.KindForbidden {
				status = http.StatusForbidden
			}
		}
		c.String(status, msg)
		return
	}
	c.Redirect(http.StatusFound, h.Cfg.FrontendOrigin+"/auth/callback?token="+token)
}

func (h *Handlers) getConfig(c *gin.Context) {
	footer, filing := h.Settings.Config(c.Request.Context())
	OK(c, gin.H{"footer_text": footer, "filing_number": filing})
}
