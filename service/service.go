package service

import "github.com/moises-ba/ms-hash-shopping-cart/model"

type DiscountServiceIf interface {
	FindDiscount(p *model.Product) (float32, error)
}

type ShoppingCartServiceIf interface {
	FindAllProducts() []*model.Product
	AddToCart(user *model.User, itemProduct *model.ItemProduct) error
	ResumeCart(user *model.User) *model.CartResume
	EmptyCart(user *model.User)
}
