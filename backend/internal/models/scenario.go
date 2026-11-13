package models

import "time"

// ScenarioGroup 场景大类（4-6 个）
type ScenarioGroup struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Slug        string    `gorm:"uniqueIndex;size:128;not null" json:"slug"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"size:512" json:"description"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ScenarioGroup) TableName() string { return "scenario_groups" }

// Scenario 客户真实场景（核心数据库备份、VMware 替代……）
type Scenario struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	GroupID   uint           `gorm:"index;not null" json:"group_id"`
	Group     *ScenarioGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Slug      string         `gorm:"uniqueIndex;size:128;not null" json:"slug"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Summary   string         `gorm:"size:512" json:"summary"`
	PainPoint string         `gorm:"type:text" json:"pain_point"`
	ValueProp string         `gorm:"type:text" json:"value_prop"`
	IconKey   string         `gorm:"size:64" json:"icon_key"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (Scenario) TableName() string { return "scenarios" }

type ScenarioProduct struct {
	ScenarioID uint `gorm:"primaryKey" json:"scenario_id"`
	ProductID  uint `gorm:"primaryKey" json:"product_id"`
}

func (ScenarioProduct) TableName() string { return "scenario_products" }

type ScenarioService struct {
	ScenarioID uint `gorm:"primaryKey" json:"scenario_id"`
	ServiceID  uint `gorm:"primaryKey" json:"service_id"`
}

func (ScenarioService) TableName() string { return "scenario_services" }
