package middleware

import (
	"net/http"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

const ctxClaimsKey = "jwt_claims"

// JWTAuth 校验 Bearer JWT，写入 claims
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}
		raw := strings.TrimSpace(h[7:])
		claims, err := auth.ParseToken(secret, raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "令牌无效"})
			return
		}
		c.Set(ctxClaimsKey, claims)
		c.Next()
	}
}

// Claims 从上下文取 claims（需先经过 JWTAuth）
func Claims(c *gin.Context) (*auth.Claims, bool) {
	v, ok := c.Get(ctxClaimsKey)
	if !ok {
		return nil, false
	}
	cl, ok := v.(*auth.Claims)
	return cl, ok
}

// RequireAdmin 需 admin 角色（须排在 JWTAuth 之后）
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		cl, ok := Claims(c)
		if !ok || cl.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			return
		}
		c.Next()
	}
}
