package dto

import (
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/view"
)

// Scenario 场景投影：场景内容对销售与客户一致。
type Scenario struct {
	ID        uint   `json:"id"`
	GroupID   uint   `json:"group_id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Summary   string `json:"summary"`
	PainPoint string `json:"pain_point,omitempty"`
	ValueProp string `json:"value_prop,omitempty"`
	IconKey   string `json:"icon_key,omitempty"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}

func ProjectScenario(s models.Scenario, _ view.Mode) Scenario {
	return Scenario{
		ID: s.ID, GroupID: s.GroupID, Slug: s.Slug, Name: s.Name, Summary: s.Summary,
		PainPoint: s.PainPoint, ValueProp: s.ValueProp, IconKey: s.IconKey,
		SortOrder: s.SortOrder, IsActive: s.IsActive,
	}
}
