package httpapi

import (
	"encoding/json"
	"io"

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

func (h *Handlers) adminExportBackupZip(c *gin.Context) {
	data, err := h.Backup.ExportZip(c.Request.Context())
	if err != nil {
		Fail(c, err)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="sales-pilot-backup.zip"`)
	c.Data(200, "application/zip", data)
}

func (h *Handlers) adminImportBackupZip(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(64 << 20); err != nil {
		Fail(c, apperr.Invalid("invalid_body", "请上传 zip 文件"))
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		Fail(c, apperr.Invalid("invalid_body", "请上传 zip 文件"))
		return
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, 64<<20))
	if err != nil {
		Fail(c, apperr.Internal(err))
		return
	}
	sections := c.Request.FormValue("sections")
	confirm := c.Request.FormValue("confirm")
	report, err := h.Backup.ImportZip(c.Request.Context(), data, []string{sections}, confirm)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, report)
}
