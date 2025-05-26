package model

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
	Category    string  `json:"category"`
	ImageUrl    string  `json:"image_url,omitempty"`
}
