package model

type ProductItemRequest struct {
	Id       int32 `json:"id"`
	Quantity int32 `json:"quantity"`
}

type ProductsRequest struct {
	Products []*ProductItemRequest `json:"products"`
}
