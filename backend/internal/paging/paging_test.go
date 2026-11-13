package paging

import "testing"

func TestNewClamps(t *testing.T) {
	r := New(0, 0)
	if r.Page != 1 || r.PageSize != 20 {
		t.Fatalf("%+v", r)
	}
	r = New(3, 500)
	if r.PageSize != 100 || r.Offset() != 200 {
		t.Fatalf("%+v", r)
	}
}

func TestOffsetAndLimit(t *testing.T) {
	r := New(2, 10)
	if r.Offset() != 10 || r.Limit() != 10 {
		t.Fatalf("offset=%d limit=%d", r.Offset(), r.Limit())
	}
}
