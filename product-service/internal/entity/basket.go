package entity

import "time"



type Basket struct{
	Id string
	UserId string
	ProductId string
	Created_at time.Time
	Updated_at time.Time
}
type ListBasketReq struct {
	Page  int64
	Limit int64
}
type ListBasketRes struct {
	TotalCount int
	Basket  []*Basket
}