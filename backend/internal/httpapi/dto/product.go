package dto

import (
	"time"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/view"
	"gorm.io/datatypes"
)

// Product 目录投影：客户视图只含白名单字段；销售视图追加内部字段。
type Product struct {
	ID                 uint           `json:"id"`
	Name               string         `json:"name"`
	Category           string         `json:"category"`
	VendorMarket       string         `json:"vendor_market"`
	SolutionCategoryID *uint          `json:"solution_category_id,omitempty"`
	ManufacturerName   string         `json:"manufacturer_name,omitempty"`
	ManufacturerLogo   string         `json:"manufacturer_logo,omitempty"`
	Description        string         `json:"description"`
	Highlights         datatypes.JSON `json:"highlights"`
	CustomerValue      string         `json:"customer_value,omitempty"`
	CustomerFeatures   datatypes.JSON `json:"customer_features,omitempty"`
	IsDraft            bool           `json:"is_draft"`
	UpdatedAt          time.Time      `json:"updated_at"`

	// 以下仅销售视图
	CreatedAt            *time.Time     `json:"created_at,omitempty"`
	SalesContactName     string         `json:"sales_contact_name,omitempty"`
	SalesContactPhone    string         `json:"sales_contact_phone,omitempty"`
	SalesContactEmail    string         `json:"sales_contact_email,omitempty"`
	PresalesContactName  string         `json:"presales_contact_name,omitempty"`
	PresalesContactPhone string         `json:"presales_contact_phone,omitempty"`
	PresalesContactEmail string         `json:"presales_contact_email,omitempty"`
	DiscoveryQuestions   datatypes.JSON `json:"discovery_questions,omitempty"`
	TargetPersonas       datatypes.JSON `json:"target_personas,omitempty"`
	TriggerEvents        string         `json:"trigger_events,omitempty"`
	CompetitorAnalysis   string         `json:"competitor_analysis,omitempty"`
	ROIMetrics           string         `json:"roi_metrics,omitempty"`
}

func ProjectProduct(p models.Product, mode view.Mode) Product {
	out := Product{
		ID:                 p.ID,
		Name:               p.Name,
		Category:           p.Category,
		VendorMarket:       p.VendorMarket,
		SolutionCategoryID: p.SolutionCategoryID,
		ManufacturerName:   p.ManufacturerName,
		ManufacturerLogo:   p.ManufacturerLogo,
		Description:        p.Description,
		Highlights:         p.Highlights,
		IsDraft:            p.IsDraft,
		UpdatedAt:          p.UpdatedAt,
	}
	if p.SolutionCategory != nil {
		out.CustomerValue = p.SolutionCategory.CustomerValue
		out.CustomerFeatures = p.SolutionCategory.CustomerFeatures
	}
	if mode.IsCustomer() {
		return out
	}
	created := p.CreatedAt
	out.CreatedAt = &created
	out.SalesContactName = p.SalesContactName
	out.SalesContactPhone = p.SalesContactPhone
	out.SalesContactEmail = p.SalesContactEmail
	out.PresalesContactName = p.PresalesContactName
	out.PresalesContactPhone = p.PresalesContactPhone
	out.PresalesContactEmail = p.PresalesContactEmail
	out.DiscoveryQuestions = p.DiscoveryQuestions
	out.TargetPersonas = p.TargetPersonas
	out.TriggerEvents = p.TriggerEvents
	out.CompetitorAnalysis = p.CompetitorAnalysis
	out.ROIMetrics = p.ROIMetrics
	return out
}

func ProjectProducts(list []models.Product, mode view.Mode) []Product {
	out := make([]Product, 0, len(list))
	for _, p := range list {
		out = append(out, ProjectProduct(p, mode))
	}
	return out
}
