package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

var contentTables = []string{
	"scenario_services",
	"scenario_products",
	"scenarios",
	"scenario_groups",
	"service_categories",
	"services",
	"cases",
	"sales_scripts",
	"products",
	"solution_vendors",
	"solution_categories",
	"topology_layers",
}

var contentSequenceTables = []string{
	"scenarios",
	"scenario_groups",
	"services",
	"cases",
	"sales_scripts",
	"products",
	"solution_vendors",
	"solution_categories",
	"topology_layers",
}

var configTables = []string{"auth_settings", "system_settings"}
var userTables = []string{"admin_api_keys", "users"}

func truncateSet(tx *gorm.DB, tables []string) error {
	sql := "TRUNCATE TABLE " + strings.Join(tables, ", ") + " RESTART IDENTITY CASCADE"
	return tx.Exec(sql).Error
}

func resetSequences(tx *gorm.DB, tables []string) error {
	for _, tbl := range tables {
		q := fmt.Sprintf(
			`SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE((SELECT MAX(id) FROM %s), 1))`,
			tbl, tbl,
		)
		if err := tx.Exec(q).Error; err != nil {
			return fmt.Errorf("重置序列 %s 失败: %w", tbl, err)
		}
	}
	return nil
}

func (r *BackupRepository) ExportContent(ctx context.Context) (*models.BackupContent, error) {
	db := r.db.WithContext(ctx)
	out := &models.BackupContent{}
	steps := []struct {
		order string
		dest  any
	}{
		{"id", &out.TopologyLayers},
		{"id", &out.SolutionCategories},
		{"id", &out.SolutionVendors},
		{"id", &out.Products},
		{"id", &out.SalesScripts},
		{"id", &out.Cases},
		{"id", &out.Services},
		{"service_id, category_id", &out.ServiceCategories},
		{"id", &out.ScenarioGroups},
		{"id", &out.Scenarios},
		{"scenario_id, product_id", &out.ScenarioProducts},
		{"scenario_id, service_id", &out.ScenarioServices},
	}
	for _, step := range steps {
		if err := db.Order(step.order).Find(step.dest).Error; err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (r *BackupRepository) ExportConfig(ctx context.Context) (*models.BackupConfig, error) {
	db := r.db.WithContext(ctx)
	out := &models.BackupConfig{}
	var auth models.AuthSettings
	if err := db.First(&auth, 1).Error; err == nil {
		out.AuthSettings = &auth
	}
	var system models.SystemSettings
	if err := db.First(&system, 1).Error; err == nil {
		out.SystemSettings = &system
	}
	return out, nil
}

func (r *BackupRepository) ExportUsers(ctx context.Context) (*models.BackupUsers, error) {
	db := r.db.WithContext(ctx)
	var dbUsers []models.User
	if err := db.Order("id").Find(&dbUsers).Error; err != nil {
		return nil, err
	}
	out := &models.BackupUsers{Users: make([]models.BackupUserRow, 0, len(dbUsers))}
	for _, u := range dbUsers {
		out.Users = append(out.Users, models.BackupUserRow{
			ID: u.ID, Username: u.Username, Email: u.Email, PasswordHash: u.PasswordHash,
			AuthProvider: u.AuthProvider, Role: u.Role, IsActive: u.IsActive,
			LastLogin: u.LastLogin, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
		})
	}
	var keys []models.AdminAPIKey
	if err := db.Order("id").Find(&keys).Error; err != nil {
		return nil, err
	}
	out.APIKeys = make([]models.BackupAPIKeyRow, 0, len(keys))
	for _, k := range keys {
		out.APIKeys = append(out.APIKeys, models.BackupAPIKeyRow{
			ID: k.ID, UserID: k.UserID, Name: k.Name, KeyHash: k.KeyHash, Prefix: k.Prefix,
			ExpiresAt: k.ExpiresAt, LastUsedAt: k.LastUsedAt, IsActive: k.IsActive,
			CreatedAt: k.CreatedAt, UpdatedAt: k.UpdatedAt,
		})
	}
	return out, nil
}

// RestoreSections 按选择的分区覆盖恢复；未选中的分区保持不动。
func (r *BackupRepository) RestoreSections(
	ctx context.Context,
	sections map[string]bool,
	content *models.BackupContent,
	cfg *models.BackupConfig,
	users *models.BackupUsers,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sess := tx.Session(&gorm.Session{FullSaveAssociations: false})

		if sections["content"] && content != nil {
			if err := truncateSet(tx, contentTables); err != nil {
				return fmt.Errorf("清空内容数据失败: %w", err)
			}
			for i := range content.TopologyLayers {
				if err := sess.Create(&content.TopologyLayers[i]).Error; err != nil {
					return fmt.Errorf("写入分层失败: %w", err)
				}
			}
			for i := range content.SolutionCategories {
				content.SolutionCategories[i].Vendors = nil
				if err := sess.Create(&content.SolutionCategories[i]).Error; err != nil {
					return fmt.Errorf("写入解决方案失败: %w", err)
				}
			}
			for i := range content.SolutionVendors {
				if err := sess.Create(&content.SolutionVendors[i]).Error; err != nil {
					return fmt.Errorf("写入厂商失败: %w", err)
				}
			}
			for i := range content.Products {
				content.Products[i].SolutionCategory = nil
				if err := sess.Create(&content.Products[i]).Error; err != nil {
					return fmt.Errorf("写入产品失败: %w", err)
				}
			}
			for i := range content.SalesScripts {
				if err := sess.Create(&content.SalesScripts[i]).Error; err != nil {
					return fmt.Errorf("写入话术失败: %w", err)
				}
			}
			for i := range content.Cases {
				if err := sess.Create(&content.Cases[i]).Error; err != nil {
					return fmt.Errorf("写入案例失败: %w", err)
				}
			}
			for i := range content.Services {
				if err := sess.Create(&content.Services[i]).Error; err != nil {
					return fmt.Errorf("写入服务失败: %w", err)
				}
			}
			for i := range content.ServiceCategories {
				if err := sess.Create(&content.ServiceCategories[i]).Error; err != nil {
					return fmt.Errorf("写入服务分类关联失败: %w", err)
				}
			}
			for i := range content.ScenarioGroups {
				if err := sess.Create(&content.ScenarioGroups[i]).Error; err != nil {
					return fmt.Errorf("写入场景组失败: %w", err)
				}
			}
			for i := range content.Scenarios {
				if err := sess.Create(&content.Scenarios[i]).Error; err != nil {
					return fmt.Errorf("写入场景失败: %w", err)
				}
			}
			for i := range content.ScenarioProducts {
				if err := sess.Create(&content.ScenarioProducts[i]).Error; err != nil {
					return fmt.Errorf("写入场景产品关联失败: %w", err)
				}
			}
			for i := range content.ScenarioServices {
				if err := sess.Create(&content.ScenarioServices[i]).Error; err != nil {
					return fmt.Errorf("写入场景服务关联失败: %w", err)
				}
			}
			if err := resetSequences(tx, contentSequenceTables); err != nil {
				return err
			}
		}

		if sections["config"] && cfg != nil {
			if err := truncateSet(tx, configTables); err != nil {
				return fmt.Errorf("清空配置失败: %w", err)
			}
			auth := models.AuthSettings{ID: 1}
			if cfg.AuthSettings != nil {
				auth = *cfg.AuthSettings
				if auth.ID == 0 {
					auth.ID = 1
				}
			}
			if err := sess.Create(&auth).Error; err != nil {
				return fmt.Errorf("写入认证配置失败: %w", err)
			}
			system := models.SystemSettings{ID: 1}
			if cfg.SystemSettings != nil {
				system = *cfg.SystemSettings
				if system.ID == 0 {
					system.ID = 1
				}
			}
			if err := sess.Create(&system).Error; err != nil {
				return fmt.Errorf("写入系统配置失败: %w", err)
			}
			if err := resetSequences(tx, configTables); err != nil {
				return err
			}
		}

		if sections["users"] && users != nil {
			if err := truncateSet(tx, userTables); err != nil {
				return fmt.Errorf("清空用户失败: %w", err)
			}
			for _, row := range users.Users {
				u := models.User{
					ID: row.ID, Username: row.Username, Email: row.Email, PasswordHash: row.PasswordHash,
					AuthProvider: row.AuthProvider, Role: row.Role, IsActive: row.IsActive,
					LastLogin: row.LastLogin, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
				}
				if u.AuthProvider == "" {
					u.AuthProvider = "local"
				}
				if err := sess.Create(&u).Error; err != nil {
					return fmt.Errorf("写入用户失败: %w", err)
				}
			}
			for _, row := range users.APIKeys {
				k := models.AdminAPIKey{
					ID: row.ID, UserID: row.UserID, Name: row.Name, KeyHash: row.KeyHash,
					Prefix: row.Prefix, ExpiresAt: row.ExpiresAt, LastUsedAt: row.LastUsedAt,
					IsActive: row.IsActive, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
				}
				if err := sess.Create(&k).Error; err != nil {
					return fmt.Errorf("写入 API 密钥失败: %w", err)
				}
			}
			if err := resetSequences(tx, userTables); err != nil {
				return err
			}
		}
		return nil
	})
}
