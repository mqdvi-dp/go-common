package types

import "math"

// CreateMetaData function to create metadata response
func CreateMetaData(page, limit, totalData int64) *Meta {
	var prevPage, nextPage int64
	totalPage := int64(math.Ceil(float64(totalData) / float64(limit)))

	if totalPage > 1 && totalPage > page {
		nextPage = page + 1
	}

	if page > 1 {
		prevPage = page - 1
	}

	return &Meta{
		TotalData:   totalData,
		TotalPage:   totalPage,
		PerPage:     limit,
		NextPage:    nextPage,
		CurrentPage: page,
		PrevPage:    prevPage,
	}
}

type Meta struct {
	TotalData   int64 `json:"total_data"`
	TotalPage   int64 `json:"total_page"`
	PerPage     int64 `json:"per_page"`
	NextPage    int64 `json:"next_page,omitempty"`
	CurrentPage int64 `json:"current_page,omitempty"`
	PrevPage    int64 `json:"prev_page,omitempty"`
}
