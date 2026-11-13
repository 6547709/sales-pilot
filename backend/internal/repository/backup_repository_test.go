package repository

import (
	"strings"
	"testing"
)

func TestTruncateTablesVersioning(t *testing.T) {
	v1 := strings.Join(truncateTables(1), ",")
	if strings.Contains(v1, "services") || strings.Contains(v1, "scenarios") {
		t.Fatalf("v1 must not touch new tables: %s", v1)
	}
	v2 := strings.Join(truncateTables(2), ",")
	for _, tbl := range []string{"services", "service_categories", "scenarios", "scenario_groups", "scenario_products", "scenario_services"} {
		if !strings.Contains(v2, tbl) {
			t.Fatalf("v2 missing %s: %s", tbl, v2)
		}
	}
}
