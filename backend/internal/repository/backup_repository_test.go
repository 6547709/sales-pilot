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

func TestSequenceTablesVersioning(t *testing.T) {
	v1 := strings.Join(sequenceTables(1), ",")
	if strings.Contains(v1, "services") {
		t.Fatalf("v1 must not reset v2 sequences: %s", v1)
	}
	v2 := strings.Join(sequenceTables(2), ",")
	for _, tbl := range []string{"services", "scenario_groups", "scenarios"} {
		if !strings.Contains(v2, tbl) {
			t.Fatalf("v2 missing sequence %s: %s", tbl, v2)
		}
	}
	for _, join := range []string{"service_categories", "scenario_products", "scenario_services"} {
		if strings.Contains(v2, join) {
			t.Fatalf("join table has no id sequence: %s", join)
		}
	}
}
