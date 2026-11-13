package httpapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/6547709/sales-pilot/backend/internal/auth"
	"github.com/6547709/sales-pilot/backend/internal/config"
	"github.com/6547709/sales-pilot/backend/internal/httpapi"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/paging"
	"github.com/6547709/sales-pilot/backend/internal/repository"
	"github.com/6547709/sales-pilot/backend/internal/service"
)

type fakeProductRepo struct {
	service.ProductRepo
	all []models.Product
}

func (f *fakeProductRepo) ListAdminAll(ctx context.Context, _ repository.AdminProductFilter) ([]models.Product, error) {
	return f.all, nil
}

func (f *fakeProductRepo) ListAdmin(ctx context.Context, _ repository.AdminProductFilter, p paging.Request) ([]models.Product, int64, error) {
	return f.all, int64(len(f.all)), nil
}

type fakeUserRepo struct {
	service.UserRepo
}

type fakeKeyRepo struct {
	service.APIKeyRepo
}

func newTestRouter() (*httptest.ResponseRecorder, http.Handler, config.Config) {
	cfg := config.Config{JWTSecret: "test-secret", FrontendOrigin: "http://localhost:3000"}
	products := service.NewProductService(&fakeProductRepo{all: []models.Product{{ID: 1, Name: "A"}}})
	apiKeys := service.NewAPIKeyService(&fakeKeyRepo{}, &fakeUserRepo{})
	r := httpapi.NewRouter(httpapi.Deps{
		Cfg:      cfg,
		Products: products,
		APIKeys:  apiKeys,
	})
	return httptest.NewRecorder(), r, cfg
}

func adminToken(t *testing.T, cfg config.Config) string {
	t.Helper()
	token, err := auth.IssueToken(cfg.JWTSecret, 1, "admin", "admin")
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestAdminProductsWithoutPageReturnsArray(t *testing.T) {
	w, r, cfg := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/products", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken(t, cfg))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if !strings.HasPrefix(strings.TrimSpace(w.Body.String()), "[") {
		t.Fatalf("body=%s", w.Body.String())
	}
}

func TestAdminProductsWithPageReturnsEnvelope(t *testing.T) {
	w, r, cfg := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/products?page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken(t, cfg))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, key := range []string{`"items"`, `"total"`, `"page"`, `"page_size"`} {
		if !strings.Contains(body, key) {
			t.Fatalf("missing %s: %s", key, body)
		}
	}
}

func TestUnauthorizedReturnsCode(t *testing.T) {
	w, r, _ := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/products", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"code"`) {
		t.Fatalf("body=%s", w.Body.String())
	}
}

func TestInvalidIDReturns400(t *testing.T) {
	w, r, cfg := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/products/abc", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken(t, cfg))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "invalid_id") {
		t.Fatalf("body=%s", w.Body.String())
	}
}
