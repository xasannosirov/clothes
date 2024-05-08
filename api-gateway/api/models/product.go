package models

type (
	Product struct {
		ID          string  `json:"product_id"`
		Name        string  `json:"product_name"`
		Category    string  `json:"category"`
		Description string  `json:"descrition"`
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
		Category    string  `json:"category"`
		Description string  `json:"descrition"`
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

	LikeReq struct {
		ProductID string `json:"product_id"`
		UserID    string `json:"user_id"`
	}

	Like struct {
		ID      string  `json:"like_id"`
		Product Product `json:"product"`
		User    User    `json:"user"`
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

	CommentReq struct {
		ProductID string `json:"produt_id"`
		UserID    string `json:"user_id"`
		Comment   string `json:"comment_message"`
	}

	Comment struct {
		ID      string  `json:"comment_id"`
		Product Product `json:"product"`
		User    User    `json:"user"`
		Comment string  `json:"comment_message"`
	}
)
