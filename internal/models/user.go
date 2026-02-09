package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Bio           string    `json:"bio,omitempty"`
	LinkedinURL   string    `json:"linkedin_url,omitempty"`
	Skills        []string  `json:"skills,omitempty"`
	WalletAddress string    `json:"wallet_address,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}
