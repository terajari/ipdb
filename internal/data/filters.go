package data

import (
	"math"

	"github.com/terajari/ipdb/internal/validator"
)

type Filters struct {
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
	Sort         string `form:"sort"`
	SortSafelist []string
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.PermitedValues[string](f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

func DefaultsFilters(f Filters) *Filters {
	if f.PageSize == 0 {
		f.PageSize = 10
	}
	if f.Page == 0 {
		f.Page = 1
	}
	if f.Sort == "" {
		f.Sort = "title"
	}
	return &f
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}
