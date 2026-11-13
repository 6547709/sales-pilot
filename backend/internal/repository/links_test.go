package repository

import (
	"reflect"
	"testing"
)

func TestDedupeIDs(t *testing.T) {
	got := dedupeIDs([]uint{3, 1, 3, 0, 2, 1})
	want := []uint{3, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}
