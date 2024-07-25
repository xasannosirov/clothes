package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	ProductCreateResponse struct {
		ProductID string `json:"product_id"`
	}

	Product struct {
		ID          string   `json:"product_id"`
		Name        string   `json:"product_name"`
		Category    string   `json:"category_id"`
		Description string   `json:"description"`
		MadeIn      string   `json:"made_in"`
		Color       []string `json:"color"`
		Size        []string `json:"size"`
		Count       int64    `json:"count"`
		Cost        float64  `json:"cost"`
		Discount    float64  `json:"discount"`
		AgeMin      int64    `json:"age_min"`
		AgeMax      int64    `json:"age_max"`
		ForGender   string   `json:"for_gender"`
		Liked       bool     `json:"liked"`
		Basket      bool     `json:"basket"`
		ImageURL    []string `json:"image_url"`
	}

	ProductReq struct {
		Name        string   `json:"product_name"`
		Category    string   `json:"category_id"`
		Description string   `json:"description"`
		MadeIn      string   `json:"made_in"`
		Color       []string `json:"color"`
		Size        []string `json:"size"`
		Count       int64    `json:"count"`
		Cost        float64  `json:"cost"`
		Discount    float64  `json:"discount"`
		AgeMin      int64    `json:"age_min"`
		AgeMax      int64    `json:"age_max"`
		ForGender   string   `json:"for_gender"`
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

func isRoman(num string) bool {
	romanRegex := `^M{0,3}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})$`
	re := regexp.MustCompile(romanRegex)

	return re.MatchString(num)
}

func isColor(color string) bool {
	colorRegex := `(?i)^(?:red|green|blue|yellow|purple|orange|pink|black|white|gray|grey|brown|cyan|magenta|lime|indigo|violet|gold|silver|beige|maroon|navy|teal|olive|coral|salmon|khaki|orchid|lavender|turquoise|peach|tan|mint|plum|apricot|amber|ivory|saffron|crimson|rose|cherry|chocolate|jade|emerald|sapphire|ruby)$`
	re := regexp.MustCompile(colorRegex)

	return re.MatchString(color)
}

func (p ProductReq) Validate() error {
	for _, size := range p.Size {
		size = strings.ToUpper(size)
		if !isRoman(size) {
			return errors.New(fmt.Sprintf("invalid size for raman number: %s", size))
		}
	}

	for _, color := range p.Color {
		color = strings.ToLower(color)
		if !isColor(color) {
			return errors.New(fmt.Sprintf("invalid color: %s", color))
		}
	}

	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Category, validation.Required),
		validation.Field(&p.Description, validation.Required),

		validation.Field(&p.Count, validation.Required, validation.Min(0)),
		validation.Field(&p.Cost, validation.Required, validation.Min(0.0)),
		validation.Field(&p.Discount, validation.Required, validation.Min(0.0), validation.Max(100.0)),
		validation.Field(&p.AgeMin, validation.Required, validation.Min(0)),
		validation.Field(&p.AgeMax, validation.Required, validation.Min(p.AgeMin)),

		validation.Field(&p.ForGender, validation.Required, validation.In("Male", "Female")),

		validation.Field(&p.MadeIn, validation.Length(0, 100)),
	)
}
