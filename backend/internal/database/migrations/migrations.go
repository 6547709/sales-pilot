package migrations

import (
	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/gorm"
)

// Migrate 按依赖顺序注册全部模型；新增表只改这里。
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 用户与凭证
		&models.User{},
		&models.AdminAPIKey{},
		// 拓扑：层 → 分类 → 厂商
		&models.TopologyLayer{},
		&models.SolutionCategory{},
		&models.SolutionVendor{},
		// 产品与关联内容
		&models.Product{},
		&models.SalesScript{},
		&models.Case{},
		// 系统设置
		&models.AuthSettings{},
		&models.SystemSettings{},
	)
}
