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
	one models.Product
}

func (f *fakeProductRepo) FindPublishedByID(ctx context.Context, id uint) (*models.Product, error) {
	p := f.one
	return &p, nil
}

func (f *fakeProductRepo) FindByID(ctx context.Context, id uint) (*models.Product, error) {
	p := f.one
	return &p, nil
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

type fakeScenarioRepo struct {
	service.ScenarioRepo
	groups    []models.ScenarioGroup
	scenarios []models.Scenario
	pc, sc    map[uint]int64
}

func (f *fakeScenarioRepo) Groups(ctx context.Context) ([]models.ScenarioGroup, error) {
	return f.groups, nil
}

func (f *fakeScenarioRepo) Scenarios(ctx context.Context) ([]models.Scenario, error) {
	return f.scenarios, nil
}

func (f *fakeScenarioRepo) CountProducts(ctx context.Context) (map[uint]int64, error) {
	return f.pc, nil
}

func (f *fakeScenarioRepo) CountServices(ctx context.Context) (map[uint]int64, error) {
	return f.sc, nil
}

func (f *fakeScenarioRepo) ProductNamesByScenario(ctx context.Context) (map[uint][]string, error) {
	return map[uint][]string{11: {"产品 A"}}, nil
}

func (f *fakeScenarioRepo) ServiceNamesByScenario(ctx context.Context) (map[uint][]string, error) {
	return map[uint][]string{11: {"服务 B"}}, nil
}

func newTestRouter() (*httptest.ResponseRecorder, http.Handler, config.Config) {
	cfg := config.Config{JWTSecret: "test-secret", FrontendOrigin: "http://localhost:3000"}
	products := service.NewProductService(&fakeProductRepo{
		all: []models.Product{{ID: 1, Name: "A"}},
		one: models.Product{ID: 1, Name: "A", SalesContactName: "张三", ROIMetrics: "35%"},
	})
	apiKeys := service.NewAPIKeyService(&fakeKeyRepo{}, &fakeUserRepo{})
	scenarios := service.NewScenarioService(&fakeScenarioRepo{
		groups:    []models.ScenarioGroup{{ID: 1, Slug: "dp", Name: "数据保护"}},
		scenarios: []models.Scenario{{ID: 11, GroupID: 1, Slug: "db", Name: "核心数据库备份", IsActive: true}},
		pc:        map[uint]int64{11: 2}, sc: map[uint]int64{11: 3},
	})
	r := httpapi.NewRouter(httpapi.Deps{
		Cfg:       cfg,
		Products:  products,
		APIKeys:   apiKeys,
		Scenarios: scenarios,
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

func TestCustomerViewStripsInternalFields(t *testing.T) {
	w, r, cfg := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/1?view=customer", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken(t, cfg))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, key := range []string{"sales_contact_name", "roi_metrics"} {
		if strings.Contains(body, key) {
			t.Fatalf("customer view leaked %s: %s", key, body)
		}
	}
}

func TestScenariosReturnsGroups(t *testing.T) {
	w, r, cfg := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/scenarios", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken(t, cfg))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, key := range []string{`"groups"`, "核心数据库备份", `"product_count"`, `"service_count"`} {
		if !strings.Contains(body, key) {
			t.Fatalf("missing %s: %s", key, body)
		}
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
