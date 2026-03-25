package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HashAPIKeySecret 与校验时使用的摘要算法一致
func HashAPIKeySecret(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// AdminAuth 管理端：优先 X-API-Key，否则 Authorization Bearer JWT（须 admin）
func AdminAuth(jwtSecret string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawKey := strings.TrimSpace(c.GetHeader("X-API-Key"))
		if rawKey != "" {
			hash := HashAPIKeySecret(rawKey)
			var k models.AdminAPIKey
			if err := db.Where("key_hash = ? AND is_active = ?", hash, true).First(&k).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API Key 无效或已停用"})
				return
			}
			if k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API Key 已过期"})
				return
			}
			var u models.User
			if err := db.First(&u, k.UserID).Error; err != nil || !u.IsActive || u.Role != "admin" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "密钥关联用户不可用或非管理员"})
				return
			}
			now := time.Now()
			_ = db.Model(&k).Update("last_used_at", now).Error
			c.Set(ctxClaimsKey, &auth.Claims{
				UserID:   u.ID,
				Username: u.Username,
				Role:     u.Role,
			})
			c.Next()
			return
		}

		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权：请使用 JWT Bearer 或 X-API-Key"})
			return
		}
		token := strings.TrimSpace(h[7:])
		claims, err := auth.ParseToken(jwtSecret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "令牌无效"})
			return
		}
		c.Set(ctxClaimsKey, claims)
		c.Next()
	}
}
