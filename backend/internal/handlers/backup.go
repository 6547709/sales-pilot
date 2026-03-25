package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 与前端约定：仅在确认框填入此令牌时才执行全量覆盖恢复
const backupConfirmToken = "RESTORE_FULL_V1"

type backupPayloadV1 struct {
	Version            int                      `json:"version"`
	ExportedAt         string                   `json:"exported_at"`
	TopologyLayers     []models.TopologyLayer   `json:"topology_layers"`
	SolutionCategories []models.SolutionCategory `json:"solution_categories"`
	SolutionVendors    []models.SolutionVendor  `json:"solution_vendors"`
	Users              []backupUserRow          `json:"users"`
	Products           []models.Product         `json:"products"`
	SalesScripts       []models.SalesScript     `json:"sales_scripts"`
	Cases              []models.Case            `json:"cases"`
	AuthSettings       *models.AuthSettings     `json:"auth_settings"`
}

type backupUserRow struct {
	ID           uint       `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	AuthProvider string     `json:"auth_provider"`
	Role         string     `json:"role"`
	IsActive     bool       `json:"is_active"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type backupImportReq struct {
	Confirm string          `json:"confirm" binding:"required"`
	Backup  backupPayloadV1 `json:"backup" binding:"required"`
}

// adminExportBackup 导出全库业务数据（含本地用户密码哈希，请妥善保管）
func (s *Server) adminExportBackup(c *gin.Context) {
	var layers []models.TopologyLayer
	_ = s.DB.Order("id").Find(&layers).Error
	var cats []models.SolutionCategory
	_ = s.DB.Order("id").Find(&cats).Error
	var vendors []models.SolutionVendor
	_ = s.DB.Order("id").Find(&vendors).Error
	var dbUsers []models.User
	_ = s.DB.Order("id").Find(&dbUsers).Error
	userRows := make([]backupUserRow, 0, len(dbUsers))
	for _, u := range dbUsers {
		userRows = append(userRows, backupUserRow{
			ID: u.ID, Username: u.Username, Email: u.Email, PasswordHash: u.PasswordHash,
			AuthProvider: u.AuthProvider, Role: u.Role, IsActive: u.IsActive,
			LastLogin: u.LastLogin, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
		})
	}
	var products []models.Product
	_ = s.DB.Order("id").Find(&products).Error
	var scripts []models.SalesScript
	_ = s.DB.Order("id").Find(&scripts).Error
	var cases []models.Case
	_ = s.DB.Order("id").Find(&cases).Error
	var st models.AuthSettings
	_ = s.DB.First(&st, 1).Error
	var stPtr *models.AuthSettings
	if st.ID != 0 {
		cp := st
		stPtr = &cp
	}

	payload := backupPayloadV1{
		Version:            1,
		ExportedAt:         time.Now().UTC().Format(time.RFC3339),
		TopologyLayers:     layers,
		SolutionCategories: cats,
		SolutionVendors:    vendors,
		Users:              userRows,
		Products:           products,
		SalesScripts:       scripts,
		Cases:              cases,
		AuthSettings:       stPtr,
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="sales-pilot-backup.json"`)
	enc := json.NewEncoder(c.Writer)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		return
	}
}

// adminImportBackup 全量覆盖恢复（事务 + TRUNCATE）
func (s *Server) adminImportBackup(c *gin.Context) {
	var req backupImportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON 格式错误或缺少 confirm / backup 字段"})
		return
	}
	if req.Confirm != backupConfirmToken {
		c.JSON(http.StatusBadRequest, gin.H{"error": "确认令牌不正确"})
		return
	}
	if req.Backup.Version != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的备份版本"})
		return
	}

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	truncateSQL := `
TRUNCATE TABLE
	cases,
	sales_scripts,
	products,
	solution_vendors,
	solution_categories,
	topology_layers,
	users,
	auth_settings
RESTART IDENTITY CASCADE`
	if err := tx.Exec(truncateSQL).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清空数据失败: " + err.Error()})
		return
	}

	sess := tx.Session(&gorm.Session{FullSaveAssociations: false})

	for i := range req.Backup.TopologyLayers {
		if err := sess.Create(&req.Backup.TopologyLayers[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "写入分层失败: " + err.Error()})
			return
		}
	}
	for i := range req.Backup.SolutionCategories {
		req.Backup.SolutionCategories[i].Vendors = nil
		if err := sess.Create(&req.Backup.SolutionCategories[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "写入解决方案失败: " + err.Error()})
			return
		}
	}
	for i := range req.Backup.SolutionVendors {
		if err := sess.Create(&req.Backup.SolutionVendors[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "写入厂商失败: " + err.Error()})
			return
		}
	}
	adminCount := 0
	for _, row := range req.Backup.Users {
		u := models.User{
			ID: row.ID, Username: row.Username, Email: row.Email, PasswordHash: row.PasswordHash,
			AuthProvider: row.AuthProvider, Role: row.Role, IsActive: row.IsActive,
			LastLogin: row.LastLogin, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		}
		if u.AuthProvider == "" {
			u.AuthProvider = "local"
		}
		if err := sess.Create(&u).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "写入用户失败: " + err.Error()})
			return
		}
		if u.AuthProvider == "local" && u.Role == "admin" && u.IsActive {
			adminCount++
		}
	}
	if adminCount == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "备份中必须至少包含一名可用的本地管理员账号"})
		return
	}
	for i := range req.Backup.Products {
		req.Backup.Products[i].SolutionCategory = nil
		if err := sess.Create(&req.Backup.Products[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "写入产品失败: " + err.Error()})
			return
		}
	}
	for i := range req.Backup.SalesScripts {
		if err := sess.Create(&req.Backup.SalesScripts[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "写入话术失败: " + err.Error()})
			return
		}
	}
	for i := range req.Backup.Cases {
		if err := sess.Create(&req.Backup.Cases[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "写入案例失败: " + err.Error()})
			return
		}
	}
	if req.Backup.AuthSettings != nil {
		st := *req.Backup.AuthSettings
		if st.ID == 0 {
			st.ID = 1
		}
		if err := sess.Create(&st).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "写入认证配置失败: " + err.Error()})
			return
		}
	} else {
		st := models.AuthSettings{ID: 1}
		if err := sess.Create(&st).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "初始化认证配置失败: " + err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, tbl := range []string{"users", "products", "topology_layers", "solution_categories", "solution_vendors", "sales_scripts", "cases", "auth_settings"} {
		q := fmt.Sprintf(
			`SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE((SELECT MAX(id) FROM %s), 1))`,
			tbl, tbl,
		)
		_ = s.DB.Exec(q).Error
	}
	s.SearchCache.Clear()
	c.JSON(http.StatusOK, gin.H{"ok": true, "message": "恢复完成，请重新登录"})
}
