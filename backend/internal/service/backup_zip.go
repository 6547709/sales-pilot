package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/apperr"
	"github.com/6547709/sales-pilot/backend/internal/models"
)

// BackupReport 导入结果
type BackupReport struct {
	Sections []string       `json:"sections"`
	Counts   map[string]int `json:"counts"`
}

var backupSections = map[string]bool{"content": true, "config": true, "users": true}

func writeJSONEntry(zw *zip.Writer, name string, value any) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}

// ExportZip 导出 zip：content.json / config.json / users.json
func (s *BackupService) ExportZip(ctx context.Context) ([]byte, error) {
	content, err := s.repo.ExportContent(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	cfg, err := s.repo.ExportConfig(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	users, err := s.repo.ExportUsers(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)
	if err := writeJSONEntry(zw, "content.json", content); err != nil {
		return nil, apperr.Internal(err)
	}
	if err := writeJSONEntry(zw, "config.json", cfg); err != nil {
		return nil, apperr.Internal(err)
	}
	if err := writeJSONEntry(zw, "users.json", users); err != nil {
		return nil, apperr.Internal(err)
	}
	if err := zw.Close(); err != nil {
		return nil, apperr.Internal(err)
	}
	return buf.Bytes(), nil
}

func parseSections(raw []string) (map[string]bool, []string, error) {
	set := make(map[string]bool, len(raw))
	ordered := make([]string, 0, len(raw))
	for _, item := range raw {
		for _, part := range strings.Split(item, ",") {
			name := strings.TrimSpace(part)
			if name == "" {
				continue
			}
			if !backupSections[name] {
				return nil, nil, apperr.Invalid("invalid_section", "未知的备份分区: "+name)
			}
			if !set[name] {
				set[name] = true
				ordered = append(ordered, name)
			}
		}
	}
	if len(set) == 0 {
		return nil, nil, apperr.Invalid("invalid_section", "请至少选择一个恢复分区")
	}
	return set, ordered, nil
}

func readZipJSON(zr *zip.Reader, name string, dest any) error {
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		return json.NewDecoder(rc).Decode(dest)
	}
	return fmt.Errorf("缺少 %s", name)
}

// ImportZip 按选择的分区覆盖恢复
func (s *BackupService) ImportZip(ctx context.Context, data []byte, sections []string, confirm string) (*BackupReport, error) {
	if confirm != BackupConfirmToken {
		return nil, apperr.Invalid("invalid_confirm", "确认令牌不正确")
	}
	set, ordered, err := parseSections(sections)
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, apperr.Invalid("invalid_zip", "备份文件不是有效的 zip")
	}
	var content models.BackupContent
	var cfg models.BackupConfig
	var users models.BackupUsers
	if set["content"] {
		if err := readZipJSON(zr, "content.json", &content); err != nil {
			return nil, apperr.Invalid("invalid_content", "zip 中缺少 content.json")
		}
	}
	if set["config"] {
		if err := readZipJSON(zr, "config.json", &cfg); err != nil {
			return nil, apperr.Invalid("invalid_config", "zip 中缺少 config.json")
		}
	}
	if set["users"] {
		if err := readZipJSON(zr, "users.json", &users); err != nil {
			return nil, apperr.Invalid("invalid_users", "zip 中缺少 users.json")
		}
		admins := 0
		for _, u := range users.Users {
			if u.AuthProvider == "local" && u.Role == "admin" && u.IsActive {
				admins++
			}
		}
		if admins == 0 {
			return nil, apperr.Invalid("no_admin", "用户备份中必须至少包含一名可用的本地管理员")
		}
	}
	if err := s.repo.RestoreSections(ctx, set, &content, &cfg, &users); err != nil {
		return nil, apperr.Internal(err)
	}
	if s.cache != nil {
		s.cache.Clear()
	}
	return &BackupReport{
		Sections: ordered,
		Counts: map[string]int{
			"products":  len(content.Products),
			"services":  len(content.Services),
			"scenarios": len(content.Scenarios),
			"users":     len(users.Users),
			"api_keys":  len(users.APIKeys),
		},
	}, nil
}
