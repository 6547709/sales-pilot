package models

import (
	"time"

	"gorm.io/datatypes"
)

// TopologyLayer 中部「云平台 / 网络 / 硬件」等分层标题（仅 column_type=central 下的分类会挂到某一层）
type TopologyLayer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Level     int       `gorm:"not null;index" json:"level"` // 展示序号，越大越靠上（如 5=应用层）
	Title     string    `gorm:"size:255;not null" json:"title"`
	Subtitle  string    `gorm:"size:512" json:"subtitle"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SolutionCategory 首页拓扑中的「解决方案」节点（如虚拟化、零信任）
type SolutionCategory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Slug        string         `gorm:"uniqueIndex;size:128;not null" json:"slug"`
	Label       string         `gorm:"size:255;not null" json:"label"`
	IconKey     string         `gorm:"size:64" json:"icon_key"` // 与前端 lucide 名称对应
	ColumnType  string         `gorm:"size:32;not null;index" json:"column_type"` // security | ops | central
	LayerID     *uint          `gorm:"index" json:"layer_id,omitempty"`
	Layer       *TopologyLayer `gorm:"foreignKey:LayerID" json:"layer,omitempty"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Keywords    datatypes.JSON `gorm:"type:jsonb" json:"keywords"` // []string，兼容搜索
	Hint        string         `gorm:"size:512" json:"hint"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Vendors     []SolutionVendor `gorm:"foreignKey:CategoryID" json:"vendors,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// SolutionVendor 某解决方案下的厂商（国内 / 国外）
type SolutionVendor struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CategoryID uint      `gorm:"index;not null" json:"category_id"`
	Category   *SolutionCategory `gorm:"foreignKey:CategoryID" json:"-"`
	Market     string    `gorm:"size:16;not null" json:"market"` // domestic | foreign
	Name       string    `gorm:"size:255;not null" json:"name"`
	SortOrder  int       `gorm:"default:0" json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
