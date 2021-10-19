package model

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type BaseProduct struct { //criado um produto base pois o campo valor no produto data_source é 'Amount' e no ItemProduct é 'unit_amount'
	Id          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsGift      bool   `json:"is_gift"`
}

type Product struct {
	BaseProduct
	Amount int32 `json:"amount"`
}

type ItemProduct struct {
	BaseProduct
	Quantity    int32   `json:"quantity"`
	TotalAmount int32   `json:"total_amount"`
	Discount    float32 `json:"discount"`
	UnitAmount  int32   `json:"unit_amount"`
}

type CartResume struct {
	TotalAmount             int32          `json:"total_amount"`
	TotalAmountWithDiscount float32        `json:"total_amount_with_discount"`
	TotalDiscount           int32          `json:"total_discount"`
	Products                []*ItemProduct `json:"products"`
}
