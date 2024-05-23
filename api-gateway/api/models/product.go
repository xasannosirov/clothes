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
		ImageURL    string  `json:"image_url"`
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
