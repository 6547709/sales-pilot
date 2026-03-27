package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/middleware"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) adminListAPIKeys(c *gin.Context) {
	var list []models.AdminAPIKey
	if err := s.DB.Order("id desc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

type createAPIKeyReq struct {
	Name      string  `json:"name" binding:"required"`
	ExpiresAt *string `json:"expires_at"` // 可选，RFC3339；省略表示永不过期
}

func (s *Server) adminCreateAPIKey(c *gin.Context) {
	cl, ok := middleware.Claims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	var req createAPIKeyReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供 name"})
		return
	}
	var exp *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "expires_at 须为 RFC3339，如 2030-12-31T23:59:59Z"})
			return
		}
		if !t.After(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "过期时间须晚于当前时间"})
			return
		}
		exp = &t
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成密钥失败"})
		return
	}
	fullKey := "sp_" + hex.EncodeToString(raw)
	hash := middleware.HashAPIKeySecret(fullKey)
	prefix := fullKey
	if len(prefix) > 20 {
		prefix = prefix[:20] + "…"
	}

	k := models.AdminAPIKey{
		UserID: cl.UserID, Name: req.Name, KeyHash: hash, Prefix: prefix,
		ExpiresAt: exp, IsActive: true,
	}
	if err := s.DB.Create(&k).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id": k.ID, "name": k.Name, "prefix": k.Prefix,
		"expires_at": k.ExpiresAt,
		"key":        fullKey,
		"warning":    "完整密钥仅显示一次，请立即保存；丢失请作废后重新创建",
	})
}

func (s *Server) adminRevokeAPIKey(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := s.DB.Model(&models.AdminAPIKey{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
