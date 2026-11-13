package httpapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/gin-gonic/gin"
)

func newCtx() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestFailMapsNotFound(t *testing.T) {
	c, w := newCtx()
	Fail(c, apperr.NotFound("product_not_found", "产品不存在"))
	if w.Code != http.StatusNotFound {
		t.Fatalf("status=%d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"code":"product_not_found"`) {
		t.Fatalf("body=%s", w.Body.String())
	}
}

func TestFailMapsConflict(t *testing.T) {
	c, w := newCtx()
	Fail(c, apperr.Conflict("scenario_group_limit", "场景组最多 6 个"))
	if w.Code != http.StatusConflict {
		t.Fatalf("status=%d", w.Code)
	}
}

func TestFailHidesInternalDetail(t *testing.T) {
	c, w := newCtx()
	Fail(c, errors.New(`pq: relation "users" does not exist`))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", w.Code)
	}
	if strings.Contains(w.Body.String(), "pq:") {
		t.Fatalf("leaked: %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "服务器内部错误") {
		t.Fatalf("body=%s", w.Body.String())
	}
}

func TestParseIDRejectsInvalid(t *testing.T) {
	c, w := newCtx()
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	if _, ok := ParseID(c, "id"); ok {
		t.Fatal("expected false")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "invalid_id") {
		t.Fatalf("body=%s", w.Body.String())
	}
}

func TestParseIDAcceptsDigits(t *testing.T) {
	c, _ := newCtx()
	c.Params = gin.Params{{Key: "id", Value: "42"}}
	id, ok := ParseID(c, "id")
	if !ok || id != 42 {
		t.Fatalf("id=%d ok=%v", id, ok)
	}
}
