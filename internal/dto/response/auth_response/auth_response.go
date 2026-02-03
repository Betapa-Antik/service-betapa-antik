package authresponse

import "github.com/google/uuid"

type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresIn int64     `json:"expires_in"` // seconds
	UserID    uuid.UUID `json:"user_id"`
}
