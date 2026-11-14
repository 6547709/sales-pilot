package repository

import (
	"context"

	"gorm.io/gorm"
)

type LabelCount struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type StatsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) ProductsByCategory(ctx context.Context) ([]LabelCount, error) {
	var rows []LabelCount
	err := r.db.WithContext(ctx).Table("products").
		Select("COALESCE(NULLIF(TRIM(category), ''), '未分类') AS label, COUNT(*) AS count").
		Group("label").Order("count DESC").Scan(&rows).Error
	return rows, err
}

func (r *StatsRepository) ProductsByMarket(ctx context.Context) ([]LabelCount, error) {
	var rows []LabelCount
	err := r.db.WithContext(ctx).Table("products").
		Select(`CASE LOWER(COALESCE(NULLIF(TRIM(vendor_market), ''), 'all'))
			WHEN 'domestic' THEN '国内'
			WHEN 'foreign' THEN '国外'
			ELSE '不区分' END AS label, COUNT(*) AS count`).
		Group("label").Order("count DESC").Scan(&rows).Error
	return rows, err
}

func (r *StatsRepository) UsersSummary(ctx context.Context) (total, admins, active int64, err error) {
	if err = r.db.WithContext(ctx).Table("users").Count(&total).Error; err != nil {
		return
	}
	if err = r.db.WithContext(ctx).Table("users").Where("role = ?", "admin").Count(&admins).Error; err != nil {
		return
	}
	err = r.db.WithContext(ctx).Table("users").Where("is_active = ?", true).Count(&active).Error
	return
}

func (r *StatsRepository) ScenariosByGroup(ctx context.Context) ([]LabelCount, error) {
	var rows []LabelCount
	err := r.db.WithContext(ctx).Table("scenarios AS s").
		Select("g.name AS label, COUNT(*) AS count").
		Joins("JOIN scenario_groups AS g ON g.id = s.group_id").
		Group("g.name").Order("count DESC").Scan(&rows).Error
	return rows, err
}

func (r *StatsRepository) ServicesSummary(ctx context.Context) (total, active int64, err error) {
	if err = r.db.WithContext(ctx).Table("services").Count(&total).Error; err != nil {
		return
	}
	err = r.db.WithContext(ctx).Table("services").Where("is_active = ?", true).Count(&active).Error
	return
}
