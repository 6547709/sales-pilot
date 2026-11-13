package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/config"
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/repository"
	"github.com/6547709/sales-pilot/backend/internal/service"
)

type fakeUserRepo struct {
	service.UserRepo
	user   *models.User
	admins int64
	saved  *models.User
}

func (f *fakeUserRepo) FindByID(ctx context.Context, id uint) (*models.User, error) {
	if f.user == nil {
		return nil, errors.New("not found")
	}
	return f.user, nil
}

func (f *fakeUserRepo) CountActiveLocalAdmins(ctx context.Context) (int64, error) {
	return f.admins, nil
}

func (f *fakeUserRepo) Save(ctx context.Context, u *models.User) error {
	f.saved = u
	return nil
}

func (f *fakeUserRepo) Delete(ctx context.Context, id uint) error { return nil }

type fakeKeyRepo struct {
	service.APIKeyRepo
	key *models.AdminAPIKey
}

func (f *fakeKeyRepo) FindByHash(ctx context.Context, hash string) (*models.AdminAPIKey, error) {
	return f.key, nil
}

func (f *fakeKeyRepo) TouchLastUsed(ctx context.Context, id uint) error { return nil }

type fakeSettingsRepo struct {
	service.SettingsRepo
	st *models.AuthSettings
}

func (f *fakeSettingsRepo) Auth(ctx context.Context) (*models.AuthSettings, error) {
	return f.st, nil
}

type fakeProductRepo struct {
	service.ProductRepo
	saved *models.Product
	label string
}

func (f *fakeProductRepo) Create(ctx context.Context, p *models.Product) error {
	f.saved = p
	return nil
}

func (f *fakeProductRepo) Reload(ctx context.Context, p *models.Product) error { return nil }

func (f *fakeProductRepo) FindCategoryLabel(ctx context.Context, id uint) (string, error) {
	return f.label, nil
}

type fakeBackupRepo struct {
	service.BackupRepo
	restored   bool
	restoreErr error
}

func (f *fakeBackupRepo) Restore(ctx context.Context, p *models.BackupPayload) error {
	f.restored = true
	return f.restoreErr
}

func TestLocalLoginDisabled(t *testing.T) {
	svc := service.NewAuthService(config.Config{}, &fakeUserRepo{}, &fakeSettingsRepo{
		st: &models.AuthSettings{LocalLoginEnabled: false},
	})
	_, _, err := svc.LoginLocal(context.Background(), "u", "p")
	if code := apperr.CodeOf(err); code != "local_login_disabled" {
		t.Fatalf("code=%s err=%v", code, err)
	}
}

func TestAuthenticateRejectsExpired(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	svc := service.NewAPIKeyService(
		&fakeKeyRepo{key: &models.AdminAPIKey{IsActive: true, ExpiresAt: &past, UserID: 1}},
		&fakeUserRepo{user: &models.User{ID: 1, Role: "admin", IsActive: true}},
	)
	_, err := svc.Authenticate(context.Background(), "sp_x")
	if code := apperr.CodeOf(err); code != "api_key_expired" {
		t.Fatalf("code=%s err=%v", code, err)
	}
}

func TestDeleteLastAdminConflicts(t *testing.T) {
	svc := service.NewUserService(&fakeUserRepo{
		user:   &models.User{ID: 7, AuthProvider: "local", Role: "admin", IsActive: true},
		admins: 1,
	})
	err := svc.DeleteLocal(context.Background(), 7)
	if code := apperr.CodeOf(err); code != "last_admin" {
		t.Fatalf("code=%s err=%v", code, err)
	}
}

func TestUpdateLocalKeepsUnsetFields(t *testing.T) {
	role := "admin"
	repo := &fakeUserRepo{user: &models.User{ID: 7, Email: "old@x", Role: "user", IsActive: true, AuthProvider: "local"}}
	svc := service.NewUserService(repo)
	if _, err := svc.UpdateLocal(context.Background(), 7, service.UpdateLocalUserInput{Role: &role}); err != nil {
		t.Fatal(err)
	}
	if repo.user.Email != "old@x" {
		t.Fatalf("email=%q", repo.user.Email)
	}
	if repo.user.Role != "admin" {
		t.Fatalf("role=%q", repo.user.Role)
	}
}

func TestProductCreateNormalizesMarket(t *testing.T) {
	repo := &fakeProductRepo{label: "虚拟化"}
	svc := service.NewProductService(repo)
	id := uint(3)
	if _, err := svc.Create(context.Background(), &models.Product{
		Name: "A", VendorMarket: "  DOMESTIC ", SolutionCategoryID: &id,
	}); err != nil {
		t.Fatal(err)
	}
	if repo.saved.VendorMarket != "domestic" {
		t.Fatalf("market=%q", repo.saved.VendorMarket)
	}
	if repo.saved.Category != "虚拟化" {
		t.Fatalf("category=%q", repo.saved.Category)
	}
}

func TestBackupImportRollsBackOnFailure(t *testing.T) {
	repo := &fakeBackupRepo{restoreErr: errors.New("boom")}
	svc := service.NewBackupService(repo, nil)
	err := svc.Import(context.Background(), &models.BackupPayload{
		Version: 1,
		Users: []models.BackupUserRow{
			{Username: "admin", Role: "admin", AuthProvider: "local", IsActive: true},
		},
	}, service.BackupConfirmToken)
	if err == nil || !repo.restored {
		t.Fatalf("err=%v restored=%v", err, repo.restored)
	}
}

type fakeScenarioRepo struct {
	service.ScenarioRepo
	groups int64
}

func (f *fakeScenarioRepo) CountGroups(ctx context.Context) (int64, error) {
	return f.groups, nil
}

type fakeServiceRepo struct {
	service.ServiceRepo
	byCat map[uint][]models.Service
}

func (f *fakeServiceRepo) ListByCategoryIDs(ctx context.Context, ids []uint) (map[uint][]models.Service, error) {
	return f.byCat, nil
}

type fakeTopologyRepo struct {
	service.TopologyRepo
	snap *repository.TopologySnapshot
}

func (f *fakeTopologyRepo) Snapshot(ctx context.Context) (*repository.TopologySnapshot, error) {
	return f.snap, nil
}

func TestScenarioGroupLimit(t *testing.T) {
	svc := service.NewScenarioService(&fakeScenarioRepo{groups: 6})
	_, err := svc.CreateGroup(context.Background(), "新组", "new-group", "")
	if code := apperr.CodeOf(err); code != "scenario_group_limit" {
		t.Fatalf("code=%s err=%v", code, err)
	}
}

func TestServiceCreateRejectsBadKind(t *testing.T) {
	svc := service.NewServiceService(&fakeServiceRepo{})
	_, err := svc.Create(context.Background(), &models.Service{Slug: "s", Name: "S", Kind: "nope"}, nil)
	if code := apperr.CodeOf(err); code != "invalid_service_kind" {
		t.Fatalf("code=%s err=%v", code, err)
	}
}

func TestTopologyTreeAttachesServices(t *testing.T) {
	repo := &fakeTopologyRepo{snap: &repository.TopologySnapshot{
		Layers: nil,
		Security: []models.SolutionCategory{
			{ID: 3, Slug: "virtualization", Label: "虚拟化", ColumnType: "security", IsActive: true},
		},
		Ops:     []models.SolutionCategory{},
		Central: map[uint][]models.SolutionCategory{},
	}}
	svcs := &fakeServiceRepo{byCat: map[uint][]models.Service{
		3: {{ID: 9, Slug: "v2v", Name: "V2V 迁移服务", Summary: "迁移"}},
	}}
	tree, err := service.NewTopologyService(repo, svcs).Tree(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(tree.Security) != 1 || len(tree.Security[0].Services) != 1 {
		t.Fatalf("services not attached: %+v", tree.Security)
	}
}

func TestDemoteLastAdminConflicts(t *testing.T) {
	role := "user"
	repo := &fakeUserRepo{
		user:   &models.User{ID: 7, AuthProvider: "local", Role: "admin", IsActive: true},
		admins: 1,
	}
	svc := service.NewUserService(repo)
	_, err := svc.UpdateLocal(context.Background(), 7, service.UpdateLocalUserInput{Role: &role})
	if code := apperr.CodeOf(err); code != "last_admin" {
		t.Fatalf("code=%s err=%v", code, err)
	}
}

type fakeScenarioAdminRepo struct {
	service.ScenarioRepo
	scenario models.Scenario
	products []models.Product
	services []models.Service
}

func (f *fakeScenarioAdminRepo) FindByID(ctx context.Context, id uint) (*models.Scenario, error) {
	s := f.scenario
	return &s, nil
}

func (f *fakeScenarioAdminRepo) FindGroup(ctx context.Context, id uint) (*models.ScenarioGroup, error) {
	return &models.ScenarioGroup{ID: id, Name: "数据保护"}, nil
}

func (f *fakeScenarioAdminRepo) ProductIDs(ctx context.Context, id uint) ([]uint, error) {
	ids := make([]uint, 0, len(f.products))
	for _, p := range f.products {
		ids = append(ids, p.ID)
	}
	return ids, nil
}

func (f *fakeScenarioAdminRepo) ServiceIDs(ctx context.Context, id uint) ([]uint, error) {
	ids := make([]uint, 0, len(f.services))
	for _, s := range f.services {
		ids = append(ids, s.ID)
	}
	return ids, nil
}

func (f *fakeScenarioAdminRepo) ListProductsByIDsAll(ctx context.Context, ids []uint) ([]models.Product, error) {
	return f.products, nil
}

func (f *fakeScenarioAdminRepo) ListServicesByIDsAll(ctx context.Context, ids []uint) ([]models.Service, error) {
	return f.services, nil
}

func TestAdminScenarioDetailIncludesInactiveAndDraft(t *testing.T) {
	repo := &fakeScenarioAdminRepo{
		scenario: models.Scenario{ID: 11, GroupID: 1, Name: "核心数据库备份", IsActive: false},
		products: []models.Product{{ID: 1, Name: "草稿产品", IsDraft: true}},
		services: []models.Service{{ID: 2, Name: "停用服务", IsActive: false}},
	}
	svc := service.NewScenarioService(repo)
	detail, err := svc.GetAdminDetail(context.Background(), 11)
	if err != nil {
		t.Fatal(err)
	}
	if len(detail.Products) != 1 || len(detail.Services) != 1 {
		t.Fatalf("links filtered: %+v", detail)
	}
}
