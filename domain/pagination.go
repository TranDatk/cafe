package domain

type Pagination struct {
	TotalRecords int64 `json:"total_records"`
	TotalPages   int64 `json:"total_pages"`
	Page         int   `json:"page"`
	PageSize     int   `json:"page_size"`
}

type PaginationResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Logic string

const (
	And Logic = "AND"
	Or  Logic = "OR"
)

type SortOrder string

const (
	Asc  SortOrder = "ASC"
	Desc SortOrder = "DESC"
)

type SortOption struct {
	Field     string    `form:"field" binding:"omitempty"`
	SortOrder SortOrder `form:"sort_order" binding:"omitempty"`
}

type Criterion struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
	Logic    Logic       `json:"logic"`
}
