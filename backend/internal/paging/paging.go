package paging

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type Request struct {
	Page     int
	PageSize int
	Sort     string
	Order    string
}

func New(page, pageSize int) Request {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Request{Page: page, PageSize: pageSize}
}

func (r Request) Offset() int {
	return (r.Page - 1) * r.PageSize
}

func (r Request) Limit() int {
	return r.PageSize
}

type Page[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}
