package dtos

type OrderLineLable struct {
	Brand        string  `json:"brand"`
	BrandAddress string  `json:"brand_address"`
	BrandEmail   string  `json:"brand_email"`
	Ean          string  `json:"ean"`
	Asin         *string `json:"asin"`
}
