package dto

import (
	"time"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/view"
	"gorm.io/datatypes"
)

// Service 服务投影：服务内容默认全部对客户可见。
type Service struct {
	ID           uint           `json:"id"`
	Slug         string         `json:"slug"`
	Name         string         `json:"name"`
	Kind         string         `json:"kind"`
	Summary      string         `json:"summary"`
	Highlights   datatypes.JSON `json:"highlights"`
	Deliverables datatypes.JSON `json:"deliverables"`
	ValueProp    string         `json:"value_prop,omitempty"`
	Features     datatypes.JSON `json:"features,omitempty"`
	IconKey      string         `json:"icon_key,omitempty"`
	SortOrder    int            `json:"sort_order"`
	IsActive     bool           `json:"is_active"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func ProjectService(s models.Service, _ view.Mode) Service {
	return Service{
		ID: s.ID, Slug: s.Slug, Name: s.Name, Kind: s.Kind, Summary: s.Summary,
		Highlights: s.Highlights, Deliverables: s.Deliverables, ValueProp: s.ValueProp,
		Features: s.Features, IconKey: s.IconKey, SortOrder: s.SortOrder,
		IsActive: s.IsActive, UpdatedAt: s.UpdatedAt,
	}
}

func ProjectServices(list []models.Service, mode view.Mode) []Service {
	out := make([]Service, 0, len(list))
	for _, s := range list {
		out = append(out, ProjectService(s, mode))
	}
	return out
}
