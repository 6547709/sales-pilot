package dto

import (
	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/view"
)

// Case 案例投影：客户视图不含痛点。
type Case struct {
	ID             uint   `json:"id"`
	ClientName     string `json:"client_name"`
	Industry       string `json:"industry,omitempty"`
	Solution       string `json:"solution"`
	ValueDelivered string `json:"value_delivered"`
	PainPoints     string `json:"pain_points,omitempty"`
}

func ProjectCase(c models.Case, mode view.Mode) Case {
	out := Case{
		ID: c.ID, ClientName: c.ClientName, Industry: c.Industry,
		Solution: c.Solution, ValueDelivered: c.ValueDelivered,
	}
	if !mode.IsCustomer() {
		out.PainPoints = c.PainPoints
	}
	return out
}

func ProjectCases(list []models.Case, mode view.Mode) []Case {
	out := make([]Case, 0, len(list))
	for _, c := range list {
		out = append(out, ProjectCase(c, mode))
	}
	return out
}
