package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type Product struct {
	ID                uuid.UUID       `json:"id" db:"id"`
	Name              string          `json:"name" db:"name"`
	Price             float64         `json:"price" db:"price"`
	Description       string          `json:"description" db:"description"`
	IsDiscount        bool            `json:"is_discount" db:"is_discount"`
	DiscountValue     sql.NullFloat64 `json:"discount_value" db:"discount_value"`
	StartDateDiscount sql.NullTime    `json:"start_date_discount" db:"start_date_discount"`
	EndDateDiscount   sql.NullTime    `json:"end_date_discount" db:"end_date_discount"`
}

func (e *Product) GenerateUUID() {
	e.ID = uuid.New()
}
