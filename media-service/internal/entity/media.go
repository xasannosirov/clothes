package entity

import "time"

type Media struct {
	Id        string
	ProductID string
	ImageUrl  string
	FileName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
