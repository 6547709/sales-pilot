package models

import (
	"time"

	"gorm.io/datatypes"
)

// Service 我公司提供的服务（V2V、备份、容灾编排等）
type Service struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Slug         string         `gorm:"uniqueIndex;size:128;not null" json:"slug"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Kind         string         `gorm:"size:32;not null;index" json:"kind"` // v2v|backup|dr|migration|upgrade|ops
	Summary      string         `gorm:"type:text" json:"summary"`
	Highlights   datatypes.JSON `gorm:"type:jsonb" json:"highlights"`
	Deliverables datatypes.JSON `gorm:"type:jsonb" json:"deliverables"`
	ValueProp    string         `gorm:"type:text" json:"value_prop"`
	Features     datatypes.JSON `gorm:"type:jsonb" json:"features"`
	IconKey      string         `gorm:"size:64" json:"icon_key"`
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (Service) TableName() string { return "services" }

// ServiceCategory 服务出现在哪些解决方案方框里
type ServiceCategory struct {
	ServiceID  uint `gorm:"primaryKey" json:"service_id"`
	CategoryID uint `gorm:"primaryKey" json:"category_id"`
}

func (ServiceCategory) TableName() string { return "service_categories" }
