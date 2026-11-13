package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/models"
)

type fakeBackupRepo struct {
	sections map[string]bool
	calls    int
}

func (f *fakeBackupRepo) Dump(ctx context.Context) (*models.BackupPayload, error) {
	return &models.BackupPayload{Version: 1}, nil
}
func (f *fakeBackupRepo) Restore(ctx context.Context, p *models.BackupPayload) error { return nil }
func (f *fakeBackupRepo) ExportContent(ctx context.Context) (*models.BackupContent, error) {
	return &models.BackupContent{Products: []models.Product{{ID: 1}}}, nil
}
func (f *fakeBackupRepo) ExportConfig(ctx context.Context) (*models.BackupConfig, error) {
	return &models.BackupConfig{}, nil
}
func (f *fakeBackupRepo) ExportUsers(ctx context.Context) (*models.BackupUsers, error) {
	return &models.BackupUsers{Users: []models.BackupUserRow{{Username: "admin"}}}, nil
}
func (f *fakeBackupRepo) RestoreSections(ctx context.Context, sections map[string]bool, content *models.BackupContent, cfg *models.BackupConfig, users *models.BackupUsers) error {
	f.calls++
	f.sections = sections
	return nil
}

func buildZip(t *testing.T, files map[string]string) []byte {
	t.Helper()
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)
	for name, body := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestParseSections(t *testing.T) {
	set, ordered, err := parseSections([]string{"content, users", "content"})
	if err != nil {
		t.Fatal(err)
	}
	if !set["content"] || !set["users"] || set["config"] {
		t.Fatalf("set=%v", set)
	}
	if len(ordered) != 2 {
		t.Fatalf("ordered=%v", ordered)
	}
	if _, _, err := parseSections([]string{"nope"}); apperr.CodeOf(err) != "invalid_section" {
		t.Fatalf("err=%v", err)
	}
	if _, _, err := parseSections(nil); apperr.CodeOf(err) != "invalid_section" {
		t.Fatalf("err=%v", err)
	}
}

func TestImportZipSelectiveSections(t *testing.T) {
	data := buildZip(t, map[string]string{
		"content.json": `{"products":[{"id":1,"name":"A"}]}`,
		"config.json":  `{}`,
		"users.json":   `{"users":[{"username":"admin","auth_provider":"local","role":"admin","is_active":true}]}`,
	})
	repo := &fakeBackupRepo{}
	svc := NewBackupService(repo, nil)
	report, err := svc.ImportZip(context.Background(), data, []string{"content,users"}, BackupConfirmToken)
	if err != nil {
		t.Fatal(err)
	}
	if repo.calls != 1 || !repo.sections["content"] || !repo.sections["users"] || repo.sections["config"] {
		t.Fatalf("sections=%v", repo.sections)
	}
	if report.Counts["products"] != 1 || report.Counts["users"] != 1 {
		t.Fatalf("report=%+v", report)
	}
}

func TestImportZipRejectsUsersWithoutAdmin(t *testing.T) {
	data := buildZip(t, map[string]string{
		"users.json": `{"users":[{"username":"u","auth_provider":"local","role":"user","is_active":true}]}`,
	})
	svc := NewBackupService(&fakeBackupRepo{}, nil)
	_, err := svc.ImportZip(context.Background(), data, []string{"users"}, BackupConfirmToken)
	if apperr.CodeOf(err) != "no_admin" {
		t.Fatalf("err=%v", err)
	}
}

func TestExportZipContainsThreeFiles(t *testing.T) {
	svc := NewBackupService(&fakeBackupRepo{}, nil)
	data, err := svc.ExportZip(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, f := range zr.File {
		names[f.Name] = true
	}
	for _, want := range []string{"content.json", "config.json", "users.json"} {
		if !names[want] {
			t.Fatalf("missing %s: %v", want, names)
		}
	}
	var content models.BackupContent
	_ = json.Unmarshal([]byte(`{}`), &content)
}
