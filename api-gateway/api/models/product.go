package models

type (
	ProductCreateResponse struct {
		ProductID string `json:"product_id"`
	}

	Product struct {
		ID          string  `json:"product_id"`
		Name        string  `json:"product_name"`
		Category    string  `json:"category_id"`
		Description string  `json:"description"`
		MadeIn      string  `json:"made_in"`
		Color       string  `json:"color"`
		Size        int64   `json:"size"`
		Count       int64   `json:"count"`
		Cost        float64 `json:"cost"`
		Discount    float64 `json:"discount"`
		AgeMin      int64   `json:"age_min"`
		AgeMax      int64   `json:"age_max"`
		ForGender   string  `json:"for_gender"`
	}

	ProductReq struct {
		Name        string  `json:"product_name"`
		Category    string  `json:"category_id"`
		Description string  `json:"description"`
		MadeIn      string  `json:"made_in"`
		Color       string  `json:"color"`
		Size        int64   `json:"size"`
		Count       int64   `json:"count"`
		Cost        float64 `json:"cost"`
		Discount    float64 `json:"discount"`
		AgeMin      int64   `json:"age_min"`
		AgeMax      int64   `json:"age_max"`
		ForGender   string  `json:"for_gender"`
	}

	ListProduct struct {
		Products []Product `json:"products"`
		Total    uint64    `json:"total_count"`
	}

	Order struct {
		ID      string  `json:"order_id"`
		Product Product `json:"product"`
		User    User    `json:"user"`
		Status  string  `json:"status"`
	}

	OrderReq struct {
		ProductID string `json:"product_id"`
		UserID    string `json:"user_id"`
	}

	ListOrder struct {
		Orders []Order `json:"orders"`
		Total  uint64  `json:"total_count"`
	}

	LikeReq struct {
		ProductID string `json:"product_id"`
		UserID    string `json:"user_id"`
	}

	Like struct {
		ID      string  `json:"like_id"`
		Product Product `json:"product"`
		User    User    `json:"user"`
	}

	ListLike struct {
		Likes []Like `json:"likes"`
		Total uint64 `json:"totol_count"`
	}

	SaveReq struct {
		ProductID string `json:"product_id"`
		UserID    string `json:"user_id"`
	}

	Save struct {
		ID      string  `json:"like_id"`
		Product Product `json:"product"`
		User    User    `json:"user"`
	}

	ListSaves struct {
		Saves []Save `json:"saves"`
		Total uint64 `json:"total_count"`
	}

	StarReq struct {
		ProductID string `json:"product_id"`
		UserID    string `json:"user_id"`
		Star      int16  `json:"star_count"`
	}

	Star struct {
		ID      string  `json:"like_id"`
		Product Product `json:"product"`
		User    User    `json:"user"`
		Star    int16   `json:"star_count"`
	}

	ListStar struct {
		Stars []Star `json:"stars"`
		Totol uint64 `json:"total_count"`
	}

	CommentReq struct {
		ProductID string `json:"product_id"`
		UserID    string `json:"user_id"`
		Comment   string `json:"comment_message"`
	}

	Comment struct {
		ID      string  `json:"comment_id"`
		Product Product `json:"product"`
		User    User    `json:"user"`
		Comment string  `json:"comment_message"`
	}

	ListComment struct {
		Comments []Comment `json:"comments"`
		Total    uint64    `json:"total_count"`
	}

	CategoryReq struct {
		Name string `json:"category_name"`
	}

	Category struct {
		ID   string `json:"category_id"`
		Name string `json:"category_name"`
	}

	ListCategory struct {
		Categories []Category `json:"categories"`
		Total      uint64     `json:"total_count"`
	}
)
