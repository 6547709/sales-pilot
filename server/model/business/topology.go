package business

import "time"

// TopologyLayer 拓扑层级
type TopologyLayer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	Level     int       `gorm:"default:0" json:"level"`         // 层级深度（用于排序）
	Title     string    `gorm:"size:128" json:"title"`          // 显示标题
	Subtitle  string    `gorm:"size:256" json:"subtitle"`       // 副标题
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TopologyLayer) TableName() string {
	return "topology_layers"
}

// TopologyCategory 拓扑分类
type TopologyCategory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	LayerID    uint      `gorm:"index" json:"layer_id"`               // 关联层级（central 类型分类用）
	ColumnType string    `gorm:"size:16;not null" json:"column_type"`  // security/ops/central
	Name       string    `gorm:"size:128;not null" json:"name"`
	Slug       string    `gorm:"size:128" json:"slug"`                 // URL slug
	Label      string    `gorm:"size:128" json:"label"`               // 显示名称
	IconKey    string    `gorm:"size:64" json:"icon_key"`             // 图标 key
	Keywords   string    `gorm:"type:text" json:"keywords"`           // JSON array of keywords
	Hint       string    `gorm:"size:256" json:"hint"`               // 提示信息
	SortOrder  int       `gorm:"default:0" json:"sort_order"`         // 排序
	IsActive   bool      `gorm:"default:true" json:"is_active"`       // 是否激活
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (TopologyCategory) TableName() string {
	return "topology_categories"
}

// TopologyVendor 拓扑厂商
type TopologyVendor struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CategoryID   uint      `gorm:"index;not null" json:"category_id"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	Website      string    `gorm:"size:512" json:"website"`
	Description  string    `gorm:"type:text" json:"description"`
	VendorMarket string    `gorm:"size:16;default:all" json:"vendor_market"` // all/domestic/foreign
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (TopologyVendor) TableName() string {
	return "topology_vendors"
}

// TopologyLayerBlock 层级区块（包含层级信息和该层级的所有分类）
type TopologyLayerBlock struct {
	Layer      TopologyLayer      `json:"layer"`
	Categories []TopologyCategory `json:"categories"`
}

// TopologyFullResponse 完整拓扑响应（供首页全景图使用）
type TopologyFullResponse struct {
	Security      []TopologyCategory `json:"security"`       // 安全体系分类
	Ops           []TopologyCategory `json:"ops"`           // 运维体系分类
	CentralLayers []TopologyLayerBlock `json:"central_layers"` // 中心层级
}
