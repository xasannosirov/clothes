package entity

import (
	"time"
)

type Product struct {
	Id          string
	Name        string
	Description string
	Category    string
	MadeIn      string
	Color       string
	Count       int64
	Cost        float32
	Discount    float32
	AgeMin      int64
	AgeMax      int64
	ForGender   string
	Size        int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ListProduct struct {
	Products   []*Product
	TotalCount uint64
}

type Filter struct {
	Name string
}

type Order struct {
	Id        string
	ProductID string
	UserID    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListOrders struct {
	Orders     []*Order
	TotalCount uint64
}

type Recom struct {
	Gender string
	Age    uint8
}

type MoveResponse struct {
	Status bool
}

type ListProductRequest struct {
	Page  uint64
	Limit uint64
	Name  string
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
	Id        string
	ProductID string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListLikes struct {
	Likes      []*LikeProduct
	TotalCount uint64
}

type SaveProduct struct {
	Id        string
	ProductID string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListSaves struct {
	Saves      []*SaveProduct
	TotalCount uint64
}

type CommentToProduct struct {
	Id        string
	ProductID string
	UserID    string
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListComments struct {
	Comments   []*CommentToProduct
	TotalCount uint64
}

type StarProduct struct {
	Id        string
	ProductID string
	UserID    string
	Stars     int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListStars struct {
	Stars      []*StarProduct
	TotalCount uint64
}

type Category struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LiestCategory struct {
	Categories []*Category
	TotalCount uint64
}
