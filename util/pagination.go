package util

import (
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

type Pagination struct {
	Limit   int    `json:"limit"`
	Page    int    `json:"page"`
	SortBy  string `json:"sort_by"`
	OrderBy string `json:"order_by"`
}

type PaginationResponse struct {
	Data   interface{} `json:"data"`
	Paging Paging      `json:"paging"`
}

type Paging struct {
	TotalRecord int    `json:"total_record"`
	TotalPage   int    `json:"total_page"`
	Page        int    `json:"page"`
	OrderBy     string `json:"order_by"`
	SortBy      string `json:"sort_by"`
	Limit       int    `json:"limit"`
}

func GeneratePaginationFromRequest(c *gin.Context) Pagination {
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}

	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}

	sortBy := c.Query("sort_by")
	if sortBy == "" {
		sortBy = "id"
	}

	orderBy := c.Query("order_by")
	if orderBy == "" {
		orderBy = "asc"
	}

	return Pagination{
		Limit:   limit,
		Page:    page,
		SortBy:  sortBy,
		OrderBy: orderBy,
	}
}

func BuildPagination(pagination Pagination, data interface{}, totalRow int64) *PaginationResponse {
	var response PaginationResponse

	response.Data = data
	response.Paging.TotalRecord = int(totalRow)
	response.Paging.Page = pagination.Page
	response.Paging.Limit = pagination.Limit
	totalPage := int(math.Ceil(float64(totalRow)) / float64(pagination.Limit))
	if (((pagination.Limit * totalPage) - int(totalRow)) * -1) > 0 {
		totalPage++
	}
	if totalRow == 1 {
		totalPage = 1
	}

	response.Paging.TotalPage = totalPage
	response.Paging.OrderBy = pagination.OrderBy
	response.Paging.SortBy = pagination.SortBy

	return &response
}
