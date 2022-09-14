package utils

import (
	"math"
	"strconv"
)

type Meta struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type PaginationMeta struct {
	Status      bool        `json:"status"`
	Message     string      `json:"message"`
	CurrentPage int         `json:"current_page"`
	NextPage    interface{} `json:"next_page"`
	PrevPage    interface{} `json:"prev_page"`
	PerPage     int         `json:"per_page"`
	PageCount   int         `json:"page_count"`
	TotalCount  int         `json:"total_count"`
}
type ModelResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}
type ModelPaginationResponse struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}
type EmptyStruct struct{}
type EmptyResponse struct {
	Content interface{} `json:"content"`
}
type HrefPagination struct {
	Href string `json:"href"`
}

func Response(msg string, data interface{}) ModelResponse {
	res := ModelResponse{}
	meta := new(Meta)

	meta.Message = msg
	res.Meta = *meta
	res.Data = data

	return res
}

func PaginationResponse(total int, page, perPage string, data interface{}) *ModelPaginationResponse {
	res := new(ModelPaginationResponse)
	convPage, _ := strconv.Atoi(page)
	convPerPage, _ := strconv.Atoi(perPage)
	page_count := int(math.Ceil(float64(total) / float64(convPerPage)))

	hasNext := false
	if float64(convPage) < float64(page_count) {
		hasNext = true
	}

	meta := PaginationMeta{
		Message:     "success",
		Status:      true,
		CurrentPage: convPage,
		NextPage:    hasNext,
		PrevPage:    convPage > 1,
		PerPage:     convPerPage,
		PageCount:   page_count,
		TotalCount:  total,
	}
	res.Meta = meta
	res.Data = data

	return res
}
