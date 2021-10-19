package model

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	Id          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UnitAmount  int32  `json:"unit_amount"`
	IsGift      bool   `json:"is_gift"`
}

type ItemProduct struct {
	Product
	Quantity    int32   `json:"quantity"`
	TotalAmount int32   `json:"total_amount"`
	Discount    float32 `json:"discount"`
}

type CartResume struct {
	TotalAmount             int32          `json:"total_amount"`
	TotalAmountWithDiscount float32        `json:"total_amount_with_discount"`
	TotalDiscount           int32          `json:"total_discount"`
	Products                []*ItemProduct `json:"products"`
}
