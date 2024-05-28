package entity

type Category struct {
	ID   string
	Name string
}

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
}

type Basket struct {
	ProductIDs []string
	UserID     string
	TotalCount int64
}
type GetBAsketReq struct {
	UserId string
	Page   int64
	Limit  int64
}

type BasketCreateReq struct {
	ProductID string
	UserID    string
}

type Order struct {
	Id        string
	ProductID string
	UserID    string
	Count     uint64
	Status    string
}

type Like struct {
	Id        string
	ProductID string
	UserID    string
}

type Params struct {
	Filter map[string]string
}

type SearchRequest struct {
	Page   uint64
	Limit  uint64
	Params map[string]string
}

type ListRequest struct {
	Page  int64
	Limit int64
}

type MoveResponse struct {
	Status bool
}

type GetWithID struct {
	ID string
}

type StatsResponse struct {
	Stats []struct {
		Product Product
		Rating  uint64
	}
}

type ListBasket struct {
	Baskets    []*Basket
	TotalCount uint64
}

type ListProduct struct {
	Products   []*Product
	TotalCount uint64
}

type ListOrders struct {
	Orders     []*Order
	TotalCount uint64
}

type ListLikes struct {
	Likes      []*Like
	TotalCount uint64
}

type LiestCategory struct {
	Categories []*Category
	TotalCount uint64
}
