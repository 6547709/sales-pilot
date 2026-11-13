package httpapi

import (
	"encoding/json"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) adminExportBackup(c *gin.Context) {
	payload, err := h.Backup.Export(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="sales-pilot-backup.json"`)
	enc := json.NewEncoder(c.Writer)
	enc.SetIndent("", "  ")
	_ = enc.Encode(payload)
}

func (h *Handlers) adminImportBackup(c *gin.Context) {
	var req struct {
		Confirm string                `json:"confirm" binding:"required"`
		Backup  *models.BackupPayload `json:"backup" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Backup == nil {
		Fail(c, apperr.Invalid("invalid_body", "JSON 格式错误或缺少 confirm / backup 字段"))
		return
	}
	if err := h.Backup.Import(c.Request.Context(), req.Backup, req.Confirm); err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"ok": true, "message": "恢复完成，请重新登录"})
}
