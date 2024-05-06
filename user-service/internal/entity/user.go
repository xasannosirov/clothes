package entity

import "time"

type User struct {
	GUID        string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Password    string
	Gender      string
	Age         uint8
	Role        string
	Refresh     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type IsUnique struct {
	Email string
}

type UpdateRefresh struct {
	UserID       string
	Role         string
	RefreshToken string
}

type UpdatePassword struct {
	UserID      string
	Role        string
	NewPassword string
}

type Response struct {
	Status bool
}
