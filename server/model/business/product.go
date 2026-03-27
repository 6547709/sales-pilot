package business

import (
	"time"

	"gorm.io/datatypes"
)

// Product 销售赋能产品内容
type Product struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"size:255;not null" json:"name"`
	Category           string         `gorm:"size:128" json:"category"` // 与所属解决方案名称同步，便于列表展示
	SolutionCategoryID *uint          `gorm:"index" json:"solution_category_id,omitempty"`
	ManufacturerName   string         `gorm:"size:255" json:"manufacturer_name"`
	SalesContactName   string         `gorm:"size:128" json:"sales_contact_name"`
	SalesContactPhone  string         `gorm:"size:64" json:"sales_contact_phone"`
	SalesContactEmail  string         `gorm:"size:255" json:"sales_contact_email"`
	PresalesContactName  string       `gorm:"size:128" json:"presales_contact_name"`
	PresalesContactPhone string       `gorm:"size:64" json:"presales_contact_phone"`
	PresalesContactEmail string       `gorm:"size:255" json:"presales_contact_email"`
	Description        string         `gorm:"type:text" json:"description"`
	Highlights         datatypes.JSON `gorm:"type:jsonb" json:"highlights"`             // []string，3 大亮点
	TargetPersonas     datatypes.JSON `gorm:"type:jsonb" json:"target_personas"`         // []string
	TriggerEvents      string         `gorm:"type:text" json:"trigger_events"`           // 触发事件
	DiscoveryQuestions datatypes.JSON `gorm:"type:jsonb" json:"discovery_questions"`     // 黄金三问
	CompetitorAnalysis string         `gorm:"type:text" json:"competitor_analysis"`      // 竞品分析
	ROIMetrics         string         `gorm:"type:text" json:"roi_metrics"`             // ROI 指标
	// VendorMarket 与首页架构卡「国内/国外」筛选一致：all=两侧均展示；domestic/foreign=仅对应分区
	VendorMarket string    `gorm:"size:16;default:all;index" json:"vendor_market"`
	IsDraft      bool      `gorm:"default:false;index" json:"is_draft"` // true=仅后台可见，前台不展示
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
