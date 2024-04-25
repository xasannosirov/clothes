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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
