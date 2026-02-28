package common

// PageParams holds list request pagination.
type PageParams struct {
	Page     int
	PageSize int
}

// DefaultPageSize is used when PageSize <= 0.
const DefaultPageSize = 20

// MaxPageSize caps list size.
const MaxPageSize = 100

// Normalize ensures Page >= 1 and PageSize in [1, MaxPageSize].
func (p *PageParams) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
}

// Offset returns (page - 1) * pageSize for DB skip.
func (p PageParams) Offset() int {
	p.Normalize()
	return (p.Page - 1) * p.PageSize
}

// PageResult holds paginated list response metadata.
type PageResult struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalCount int64 `json:"total_count"`
}
