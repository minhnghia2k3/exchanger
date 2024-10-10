package store

import (
	"math"
	"strings"
)

type Filter struct {
	Page         int      `json:"page,omitempty" validate:"omitempty,min=1"`
	PageSize     int      `json:"pageSize,omitempty" validate:"omitempty,min=1"`
	Sort         string   `json:"sort,omitempty" validate:"omitempty"`
	SortSafeList []string `json:"sort_safe_list,omitempty"`
	Search       string   `json:"search,omitempty" validate:"omitempty"`
}

type Metadata struct {
	CurrentPage int `json:"current_page,omitempty"`
	PageSize    int `json:"page_size,omitempty"`
	FirstPage   int `json:"first_page,omitempty"`
	LastPage    int `json:"last_page,omitempty"`
	TotalRecord int `json:"total_record,omitempty"`
}

// Check if client provided sort column is in sortSafeList
func (f *Filter) sortColumn() string {
	val := strings.TrimPrefix(f.Sort, "-")
	for _, safeValue := range f.SortSafeList {
		if val == safeValue {
			return val
		}
	}

	// If user provided unsafe value => return id as default
	return "id"
}

// Return sort direction ("ASC", "DESC") depend on provided prefix
func (f *Filter) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

// Calculate offset
func (f *Filter) calculateOffset() int {
	return (f.Page - 1) * f.PageSize
}

// Calculate limit
func (f *Filter) limit() int {
	return f.PageSize
}

// Calculate metadata
func (f *Filter) calculateMetadata(totalRecords int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage: f.Page,
		PageSize:    f.PageSize,
		FirstPage:   1,
		LastPage:    int(math.Ceil(float64(totalRecords) / float64(f.PageSize))),
		TotalRecord: totalRecords,
	}
}
