package dtos

type ApiOrderCreate struct {
	FatherOrderName string               `json:"father_name"`
	Orders          []ApiOrderItemCreate `json:"orders"`
}

type ApiOrderItemCreate struct {
	OrderName  string          `json:"name"`
	OrderLines []ApiOrderLines `json:"order_lines"`
}

type ApiOrderLines struct {
	Asin     string `json:"asin"`
	Quantity int    `json:"quantity"`
}
