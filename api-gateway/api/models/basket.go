package models

type Basket struct {
	UserId    string
	ProductId []Product
	TotalCount     int64
}
type BasketCeateReq struct {
	ProductId string
}
type ListBasket struct {
	Baskets []Basket
	Total   int64
}
