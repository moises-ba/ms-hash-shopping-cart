package repository

import (
	"github.com/moises-ba/ms-hash-shopping-cart/model"
)

type ShoppingCartRepositoryIf interface {
	FindAllProducts() []*model.Product
	FindGifts() []*model.Product
	FindProducById(id int32) *model.Product
	AddToCart(user *model.User, itemProduct *model.ItemProduct) error
	AddGiftToCart(user *model.User, itemProduct *model.ItemProduct, canAddGift func(itensProducts []*model.ItemProduct) bool) error
	ResumeCart(user *model.User) *model.CartResume
	EmptyCart(user *model.User)
}
