package entity

import (
	"time"
)

type Product struct {
	Id             string
	Name           string
	Description    string
	Category       string
	MadeIn         string
	Color          string
	Count          int64
	Cost           float32
	Discount       float32
	AgeMin         int64
	AgeMax         int64
	TemperatureMin int64
	TemperatureMax int64
	ForGender      string
	Size           int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Filter struct {
	Name           string
	Description    string
	Category       string
	MadeIn         string
	Color          string
	Count          int64
	Cost           int64
	Discount       int64
	AgeMin         int64
	AgeMax         int64
	TemperatureMin int64
	TemperatureMax int64
	ForGender      string
	Page           int64
	Limit          int64
}

type Order struct {
	Id        string
	ProductID string
	UserID    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Recom struct {
	Gender string
	Age    uint8
}

type MoveResponse struct {
	Status bool
}

type ListRequest struct {
	Page  int64
	Limit int64
}

type GetWithID struct {
	ID string
}

type DeleteResponse struct {
	Status bool
}

type LikeProduct struct {
	Id         string
	Product_id string
	User_id    string
	Created_at time.Time
	Updated_at time.Time
	Deleted_at time.Time
}
type SaveProduct struct {
	Id         string
	Product_id string
	User_id    string
	Created_at time.Time
	Updated_at time.Time
	Deleted_at time.Time
}

type CommentToProduct struct{
	Id string
	Product_Id string
	UserId string
	Comment string
	Created_at time.Time
	Updated_at time.Time
	Deleted_at time.Time
}


   