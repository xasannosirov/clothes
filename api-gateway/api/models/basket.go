package models

type Basket struct {
	UserId    string
	ProductId []string
	Count     int64
}
type BasketCeateReq struct {
	ProductId string
}
type ListBasket struct {
	Baskets []Basket
	Total   int64
}
