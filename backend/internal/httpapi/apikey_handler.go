package httpapi

import (
	"time"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/httpapi/middleware"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) adminListAPIKeys(c *gin.Context) {
	list, err := h.APIKeys.List(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, list)
}

func (h *Handlers) adminCreateAPIKey(c *gin.Context) {
	cl, ok := middleware.Claims(c)
	if !ok {
		Fail(c, apperr.Unauthorized("unauthorized", "未授权"))
		return
	}
	var body struct {
		Name      string  `json:"name"`
		ExpiresAt *string `json:"expires_at"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "请提供 name"))
		return
	}
	var exp *time.Time
	if body.ExpiresAt != nil && *body.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, *body.ExpiresAt)
		if err != nil {
			Fail(c, apperr.Invalid("invalid_expiry", "expires_at 须为 RFC3339，如 2030-12-31T23:59:59Z"))
			return
		}
		exp = &t
	}
	fullKey, record, err := h.APIKeys.Create(c.Request.Context(), cl.UserID, body.Name, exp)
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, gin.H{
		"id":         record.ID,
		"name":       record.Name,
		"prefix":     record.Prefix,
		"expires_at": record.ExpiresAt,
		"key":        fullKey,
		"warning":    "完整密钥仅显示一次，请立即保存；丢失请作废后重新创建",
	})
}

func (h *Handlers) adminRevokeAPIKey(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}
	if err := h.APIKeys.Revoke(c.Request.Context(), id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}
