package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/config"
	"github.com/6547709/sales-pilot/backend/internal/models"
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
