package database

import (
	"github.com/6547709/sales-pilot/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect 连接 PostgreSQL 并迁移
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&models.User{},
		&models.AdminAPIKey{},
		&models.TopologyLayer{},
		&models.SolutionCategory{},
		&models.SolutionVendor{},
		&models.Product{},
		&models.SalesScript{},
		&models.Case{},
		&models.AuthSettings{},
	); err != nil {
		return nil, err
	}
	return db, nil
}
