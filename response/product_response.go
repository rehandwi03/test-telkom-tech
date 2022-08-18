package response

import (
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
	ID                uuid.UUID  `json:"id"`
	Name              string     `json:"name"`
	Price             float64    `json:"price"`
	Description       string     `json:"description"`
	IsDiscount        bool       `json:"is_discount"`
	StartDateDiscount *time.Time `json:"start_date_discount"`
	EndDateDiscount   *time.Time `json:"end_date_discount"`
	DiscountValue     float64    `json:"discount_value"`
}
