package middleware

import (
	"net/http"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/service"
	"github.com/gin-gonic/gin"
)

const ctxClaimsKey = "jwt_claims"

// Claims 从上下文取 claims（需先经过 AdminAuth/JWTAuth）
func Claims(c *gin.Context) (*auth.Claims, bool) {
	v, ok := c.Get(ctxClaimsKey)
	if !ok {
		return nil, false
	}
	cl, ok := v.(*auth.Claims)
	return cl, ok
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		cl, ok := Claims(c)
		if !ok || cl.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "需要管理员权限", "code": "forbidden"})
			return
		}
		c.Next()
	}
}

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权", "code": "unauthorized"})
			return
		}
		claims, err := auth.ParseToken(secret, strings.TrimSpace(h[7:]))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "令牌无效", "code": "invalid_token"})
			return
		}
		c.Set(ctxClaimsKey, claims)
		c.Next()
	}
}

// AdminAuth 优先 X-API-Key，否则 Authorization Bearer JWT
func AdminAuth(jwtSecret string, apiKeys *service.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if raw := strings.TrimSpace(c.GetHeader("X-API-Key")); raw != "" {
			u, err := apiKeys.Authenticate(c.Request.Context(), raw)
			if err != nil {
				abortWith(c, err)
				return
			}
			c.Set(ctxClaimsKey, &auth.Claims{UserID: u.ID, Username: u.Username, Role: u.Role})
			c.Next()
			return
		}
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权：请使用 JWT Bearer 或 X-API-Key", "code": "unauthorized"})
			return
		}
		claims, err := auth.ParseToken(jwtSecret, strings.TrimSpace(h[7:]))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "令牌无效", "code": "invalid_token"})
			return
		}
		c.Set(ctxClaimsKey, claims)
		c.Next()
	}
}

func abortWith(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	msg := "服务器内部错误"
	if e, ok := apperr.From(err); ok {
		msg = e.Message
		switch e.Kind {
		case apperr.KindInvalid:
			status = http.StatusBadRequest
		case apperr.KindUnauthorized:
			status = http.StatusUnauthorized
		case apperr.KindForbidden:
			status = http.StatusForbidden
		case apperr.KindNotFound:
			status = http.StatusNotFound
		case apperr.KindConflict:
			status = http.StatusConflict
		}
	}
	c.AbortWithStatusJSON(status, gin.H{"error": msg, "code": apperr.CodeOf(err)})
}
