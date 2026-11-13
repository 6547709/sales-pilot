package apperr

import (
	"fmt"
	"testing"
)

func TestFromUnwraps(t *testing.T) {
	base := NotFound("product_not_found", "产品不存在")
	wrapped := fmt.Errorf("load: %w", base)
	got, ok := From(wrapped)
	if !ok || got.Code != "product_not_found" {
		t.Fatalf("got %+v ok=%v", got, ok)
	}
}

func TestCodeOfNonDomainError(t *testing.T) {
	if code := CodeOf(fmt.Errorf("boom")); code != "" {
		t.Fatalf("code=%q", code)
	}
}

func TestInternalKeepsCause(t *testing.T) {
	base := fmt.Errorf("pq: relation does not exist")
	e := Internal(base)
	if e.Kind != KindInternal || e.Code != "internal_error" {
		t.Fatalf("%+v", e)
	}
	if e.Unwrap() != base {
		t.Fatalf("cause lost")
	}
}
