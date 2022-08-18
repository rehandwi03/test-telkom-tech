package request

import (
	"interview-telkom-6/util"
)

type ProductAddRequest struct {
	Name              string  `json:"name" binding:"required"`
	Price             float64 `json:"price" binding:"required"`
	Description       string  `json:"description" binding:"required"`
	IsDiscount        bool    `json:"is_discount"`
	StartDateDiscount string  `json:"start_date_discount"`
	EndDateDiscount   string  `json:"end_date_discount"`
	DiscountValue     float64 `json:"discount_value"`
}

type ProductCriteria struct {
	Search  string `json:"search"`
	Popular string `json:"popular"`
	util.Pagination
}
