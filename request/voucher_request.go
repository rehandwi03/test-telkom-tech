package request

import "interview-telkom-6/util"

type VoucherAddRequest struct {
	Name            string  `json:"name"`
	MinOrder        float64 `json:"min_order"`
	MaxUsagePerUser int     `json:"max_usage_per_user"`
	Value           float64 `json:"value"`
	StartDate       string  `json:"start_date"`
	EndDate         string  `json:"end_date"`
}

type VoucherCriteria struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	FullName  string `json:"full_name"`
	util.Pagination
}
