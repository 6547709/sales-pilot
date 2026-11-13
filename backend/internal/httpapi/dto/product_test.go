package dto

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/6547709/sales-pilot/backend/internal/models"
	"github.com/6547709/sales-pilot/backend/internal/view"
	"gorm.io/datatypes"
)

func TestCustomerProductHasNoInternalFields(t *testing.T) {
	p := models.Product{
		ID: 1, Name: "A",
		SalesContactName: "张三", SalesContactPhone: "138",
		DiscoveryQuestions: datatypes.JSON(`["q"]`),
		TargetPersonas:     datatypes.JSON(`{"CIO":"成本"}`),
		TriggerEvents:      "信创",
		CompetitorAnalysis: "竞品",
		ROIMetrics:         "35%",
	}
	b, err := json.Marshal(ProjectProduct(p, view.Customer))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, key := range []string{
		"sales_contact_name", "sales_contact_phone", "presales_contact_name",
		"discovery_questions", "target_personas", "trigger_events",
		"competitor_analysis", "roi_metrics",
	} {
		if strings.Contains(s, key) {
			t.Fatalf("customer view leaked %s: %s", key, s)
		}
	}
}

func TestSalesProductKeepsInternalFields(t *testing.T) {
	p := models.Product{ID: 1, Name: "A", SalesContactName: "张三", ROIMetrics: "35%"}
	b, _ := json.Marshal(ProjectProduct(p, view.Sales))
	s := string(b)
	for _, key := range []string{"sales_contact_name", "roi_metrics"} {
		if !strings.Contains(s, key) {
			t.Fatalf("sales view lost %s: %s", key, s)
		}
	}
}

func TestCustomerProductIncludesCategoryValue(t *testing.T) {
	p := models.Product{
		ID: 1, Name: "A",
		SolutionCategory: &models.SolutionCategory{
			CustomerValue:    "业务不中断替换",
			CustomerFeatures: datatypes.JSON(`["兼容 x86"]`),
		},
	}
	b, _ := json.Marshal(ProjectProduct(p, view.Customer))
	s := string(b)
	if !strings.Contains(s, "customer_value") || !strings.Contains(s, "customer_features") {
		t.Fatalf("missing customer fields: %s", s)
	}
}

func TestCustomerCaseDropsPainPoints(t *testing.T) {
	c := models.Case{ID: 1, ClientName: "某银行", PainPoints: "窗口不足", Solution: "S", ValueDelivered: "V"}
	b, _ := json.Marshal(ProjectCase(c, view.Customer))
	if strings.Contains(string(b), "pain_points") {
		t.Fatalf("leaked pain points: %s", b)
	}
}
