package models

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Skills         []string  `json:"skills,omitempty"`
	Salary         string    `json:"salary,omitempty"`
	Location       string    `json:"location,omitempty"`
	UserID         uuid.UUID `json:"user_id"`
	PaymentTxHash  string    `json:"payment_tx_hash,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}
