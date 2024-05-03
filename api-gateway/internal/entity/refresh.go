package entity

import "time"

// RefreshToken holds refresh token information.
type RefreshToken struct {
	GUID        string    `json:"guid"`
	RefreshToken string    `json:"refresh_token"`
	ExpiryDate   time.Time `json:"expiry_date"`
	CreatedAt    time.Time `json:"created_at"`
}
