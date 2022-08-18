package entity

import "github.com/google/uuid"

type Cart struct {
	ID       uuid.UUID `json:"id" db:"id"`
	FullName string    `json:"full_name" db:"full_name"`
}

func (e *Cart) GenerateUUID() {
	e.ID = uuid.New()
}
