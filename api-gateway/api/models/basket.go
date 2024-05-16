package models



type Basket struct{
	Id string
	UserId string
	ProductId string
}
type BasketCeateReq struct {
	ProductId string
}
type ListBasket struct{
	Baskets []Basket
	Total int64
}