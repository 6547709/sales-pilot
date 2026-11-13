package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

type BackupRepository struct {
	db *gorm.DB
}

func NewBackupRepository(db *gorm.DB) *BackupRepository {
	return &BackupRepository{db: db}
}

func (r *BackupRepository) Dump(ctx context.Context) (*models.BackupPayload, error) {
	db := r.db.WithContext(ctx)
	var layers []models.TopologyLayer
	if err := db.Order("id").Find(&layers).Error; err != nil {
		return nil, err
	}
	var cats []models.SolutionCategory
	if err := db.Order("id").Find(&cats).Error; err != nil {
		return nil, err
	}
	var vendors []models.SolutionVendor
	if err := db.Order("id").Find(&vendors).Error; err != nil {
		return nil, err
	}
	var dbUsers []models.User
	if err := db.Order("id").Find(&dbUsers).Error; err != nil {
		return nil, err
	}
	userRows := make([]models.BackupUserRow, 0, len(dbUsers))
	for _, u := range dbUsers {
		userRows = append(userRows, models.BackupUserRow{
			ID: u.ID, Username: u.Username, Email: u.Email, PasswordHash: u.PasswordHash,
			AuthProvider: u.AuthProvider, Role: u.Role, IsActive: u.IsActive,
			LastLogin: u.LastLogin, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
		})
	}
	var products []models.Product
	if err := db.Order("id").Find(&products).Error; err != nil {
		return nil, err
	}
	var scripts []models.SalesScript
	if err := db.Order("id").Find(&scripts).Error; err != nil {
		return nil, err
	}
	var cases []models.Case
	if err := db.Order("id").Find(&cases).Error; err != nil {
		return nil, err
	}
	var st models.AuthSettings
	var stPtr *models.AuthSettings
	if err := db.First(&st, 1).Error; err == nil && st.ID != 0 {
		cp := st
		stPtr = &cp
	}
	var services []models.Service
	if err := db.Order("id").Find(&services).Error; err != nil {
		return nil, err
	}
	var serviceCategories []models.ServiceCategory
	if err := db.Order("service_id, category_id").Find(&serviceCategories).Error; err != nil {
		return nil, err
	}
	var scenarioGroups []models.ScenarioGroup
	if err := db.Order("id").Find(&scenarioGroups).Error; err != nil {
		return nil, err
	}
	var scenarios []models.Scenario
	if err := db.Order("id").Find(&scenarios).Error; err != nil {
		return nil, err
	}
	var scenarioProducts []models.ScenarioProduct
	if err := db.Order("scenario_id, product_id").Find(&scenarioProducts).Error; err != nil {
		return nil, err
	}
	var scenarioServices []models.ScenarioService
	if err := db.Order("scenario_id, service_id").Find(&scenarioServices).Error; err != nil {
		return nil, err
	}
	return &models.BackupPayload{
		Version:            2,
		ExportedAt:         time.Now().UTC().Format(time.RFC3339),
		TopologyLayers:     layers,
		SolutionCategories: cats,
		SolutionVendors:    vendors,
		Users:              userRows,
		Products:           products,
		SalesScripts:       scripts,
		Cases:              cases,
		AuthSettings:       stPtr,
		Services:           services,
		ServiceCategories:  serviceCategories,
		ScenarioGroups:     scenarioGroups,
		Scenarios:          scenarios,
		ScenarioProducts:   scenarioProducts,
		ScenarioServices:   scenarioServices,
	}, nil
}

var backupTables = []string{
	"cases",
	"sales_scripts",
	"products",
	"solution_vendors",
	"solution_categories",
	"topology_layers",
	"users",
	"auth_settings",
}

var v2Tables = []string{
	"scenario_services",
	"scenario_products",
	"scenarios",
	"scenario_groups",
	"service_categories",
	"services",
}

func truncateTables(version int) []string {
	if version >= 2 {
		return append(append([]string{}, backupTables...), v2Tables...)
	}
	return backupTables
}

// Restore 全量覆盖恢复（事务 + TRUNCATE）
func (r *BackupRepository) Restore(ctx context.Context, payload *models.BackupPayload) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tables := truncateTables(payload.Version)
		truncateSQL := "TRUNCATE TABLE " + strings.Join(tables, ", ") + " RESTART IDENTITY CASCADE"
		if err := tx.Exec(truncateSQL).Error; err != nil {
			return fmt.Errorf("清空数据失败: %w", err)
		}
		sess := tx.Session(&gorm.Session{FullSaveAssociations: false})
		for i := range payload.TopologyLayers {
			if err := sess.Create(&payload.TopologyLayers[i]).Error; err != nil {
				return fmt.Errorf("写入分层失败: %w", err)
			}
		}
		for i := range payload.SolutionCategories {
			payload.SolutionCategories[i].Vendors = nil
			if err := sess.Create(&payload.SolutionCategories[i]).Error; err != nil {
				return fmt.Errorf("写入解决方案失败: %w", err)
			}
		}
		for i := range payload.SolutionVendors {
			if err := sess.Create(&payload.SolutionVendors[i]).Error; err != nil {
				return fmt.Errorf("写入厂商失败: %w", err)
			}
		}
		for _, row := range payload.Users {
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
		for i := range payload.Products {
			payload.Products[i].SolutionCategory = nil
			if err := sess.Create(&payload.Products[i]).Error; err != nil {
				return fmt.Errorf("写入产品失败: %w", err)
			}
		}
		for i := range payload.SalesScripts {
			if err := sess.Create(&payload.SalesScripts[i]).Error; err != nil {
				return fmt.Errorf("写入话术失败: %w", err)
			}
		}
		for i := range payload.Cases {
			if err := sess.Create(&payload.Cases[i]).Error; err != nil {
				return fmt.Errorf("写入案例失败: %w", err)
			}
		}
		if payload.AuthSettings != nil {
			st := *payload.AuthSettings
			if st.ID == 0 {
				st.ID = 1
			}
			if err := sess.Create(&st).Error; err != nil {
				return fmt.Errorf("写入认证配置失败: %w", err)
			}
		} else {
			if err := sess.Create(&models.AuthSettings{ID: 1}).Error; err != nil {
				return fmt.Errorf("初始化认证配置失败: %w", err)
			}
		}
		if payload.Version >= 2 {
			for i := range payload.Services {
				if err := sess.Create(&payload.Services[i]).Error; err != nil {
					return fmt.Errorf("写入服务失败: %w", err)
				}
			}
			for i := range payload.ServiceCategories {
				if err := sess.Create(&payload.ServiceCategories[i]).Error; err != nil {
					return fmt.Errorf("写入服务分类关联失败: %w", err)
				}
			}
			for i := range payload.ScenarioGroups {
				if err := sess.Create(&payload.ScenarioGroups[i]).Error; err != nil {
					return fmt.Errorf("写入场景组失败: %w", err)
				}
			}
			for i := range payload.Scenarios {
				if err := sess.Create(&payload.Scenarios[i]).Error; err != nil {
					return fmt.Errorf("写入场景失败: %w", err)
				}
			}
			for i := range payload.ScenarioProducts {
				if err := sess.Create(&payload.ScenarioProducts[i]).Error; err != nil {
					return fmt.Errorf("写入场景产品关联失败: %w", err)
				}
			}
			for i := range payload.ScenarioServices {
				if err := sess.Create(&payload.ScenarioServices[i]).Error; err != nil {
					return fmt.Errorf("写入场景服务关联失败: %w", err)
				}
			}
		}
		for _, tbl := range backupTables {
			q := fmt.Sprintf(
				`SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE((SELECT MAX(id) FROM %s), 1))`,
				tbl, tbl,
			)
			if err := tx.Exec(q).Error; err != nil {
				return fmt.Errorf("重置序列 %s 失败: %w", tbl, err)
			}
		}
		return nil
	})
}
