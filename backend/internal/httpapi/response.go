package httpapi

import (
	"log"
	"net/http"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/gin-gonic/gin"
)

type problem struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func OK(c *gin.Context, v any) {
	c.JSON(http.StatusOK, v)
}

func Created(c *gin.Context, v any) {
	c.JSON(http.StatusCreated, v)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Fail 把领域错误映射为 {error, code}；未知错误统一 500，不泄露内部细节。
func Fail(c *gin.Context, err error) {
	e, ok := apperr.From(err)
	if !ok {
		log.Printf("unhandled error: %v", err)
		c.JSON(http.StatusInternalServerError, problem{Error: "服务器内部错误", Code: "internal_error"})
		return
	}
	status := statusFor(e.Kind)
	message := e.Message
	if status == http.StatusInternalServerError {
		log.Printf("internal error [%s]: %v", e.Code, e.Err)
		message = "服务器内部错误"
	}
	c.JSON(status, problem{Error: message, Code: e.Code})
}

func statusFor(k apperr.Kind) int {
	switch k {
	case apperr.KindInvalid:
		return http.StatusBadRequest
	case apperr.KindUnauthorized:
		return http.StatusUnauthorized
	case apperr.KindForbidden:
		return http.StatusForbidden
	case apperr.KindNotFound:
		return http.StatusNotFound
	case apperr.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
