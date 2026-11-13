package repository

import (
	"testing"

	"github.com/6547709/sales-pilot/backend/internal/paging"
)

func TestOrderClauseWhitelist(t *testing.T) {
	cases := []struct {
		req  paging.Request
		want string
	}{
		{paging.Request{}, "updated_at desc"},
		{paging.Request{Sort: "id", Order: "asc"}, "id asc"},
		{paging.Request{Sort: "name", Order: "DESC"}, "name desc"},
		{paging.Request{Sort: "created_at", Order: "asc"}, "created_at asc"},
		{paging.Request{Sort: "description", Order: "asc"}, "updated_at desc"},
		{paging.Request{Sort: "id; drop table users", Order: "asc"}, "updated_at desc"},
	}
	for _, c := range cases {
		if got := orderClause(c.req); got != c.want {
			t.Fatalf("req=%+v got=%q want=%q", c.req, got, c.want)
		}
	}
}
