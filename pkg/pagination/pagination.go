package pagination

import (
	"fmt"
	"strconv"
	"strings"
)

type SortOrder string

const (
	ASC  SortOrder = "asc"
	DESC SortOrder = "desc"
)

type Params struct {
	Page      int
	PerPage   int
	SortBy    string
	SortOrder SortOrder
	All       bool
}

type Meta struct {
	HasNext     bool  `json:"hasNext"`
	HasPrevious bool  `json:"hasPrevious"`
	Page        int   `json:"page"`
	PerPage     int   `json:"perPage"`
	Total       int64 `json:"total"`
}

func NewParams(q map[string]string) Params {
	toInt := func(s string, def int) int {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			return v
		}
		return def
	}

	page := toInt(q["page"], 1)
	per := toInt(q["perPage"], 10)
	sortBy := q["sortBy"]
	if sortBy == "" {
		sortBy = "created_at"
	}
	order := SortOrder(strings.ToLower(q["sortOrder"]))
	if order != ASC {
		order = DESC
	}
	all := strings.ToLower(q["all"]) == "true"

	return Params{
		Page:      page,
		PerPage:   per,
		SortBy:    sortBy,
		SortOrder: order,
		All:       all,
	}
}

func (p Params) Offset() int {
	if p.Page <= 1 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}

func (p Params) OrderClause(defaultCol string) string {
	col := p.SortBy
	if col == "" {
		col = defaultCol
	}
	return fmt.Sprintf("%s %s", col, p.SortOrder)
}

func BuildMeta(total int64, p Params) Meta {
	hasNext := total > int64(p.PerPage*p.Page)
	prevItems := (p.Page - 1) * p.PerPage
	hasPrev := p.Page > 1 && prevItems < int(total)

	return Meta{
		HasNext:     hasNext,
		HasPrevious: hasPrev,
		Page:        p.Page,
		PerPage:     p.PerPage,
		Total:       total,
	}
}
