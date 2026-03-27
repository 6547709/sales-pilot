package business

import "time"

// TopologyLayer 拓扑层级
type TopologyLayer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TopologyLayer) TableName() string {
	return "topology_layers"
}

// TopologyCategory 拓扑分类
type TopologyCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	LayerID   uint      `gorm:"index;not null" json:"layer_id"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
