package httpapi

import (
	"strconv"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/gin-gonic/gin"
)

// ParseID 解析路径参数；失败时写出 400 invalid_id 并返回 false。
func ParseID(c *gin.Context, name string) (uint, bool) {
	raw := strings.TrimSpace(c.Param(name))
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		Fail(c, apperr.Invalid("invalid_id", "无效的 ID"))
		return 0, false
	}
	return uint(id), true
}

// QueryInt 读取整数查询参数，缺失或非法时返回默认值。
func QueryInt(c *gin.Context, name string, def int) int {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return def
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return n
}
