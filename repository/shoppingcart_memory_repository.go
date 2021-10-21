package repository

import (
	"sync"

	"github.com/moises-ba/ms-hash-shopping-cart/model"
	"github.com/moises-ba/ms-hash-shopping-cart/utils"
)

//responsavel por encapsular uma lista de produtos e tratar concorrencia na insercao
type productListHolder struct {
	mu       *sync.Mutex
	products []*model.ItemProduct
}

func (ph *productListHolder) Add(pItemProduct *model.ItemProduct) {
	if ph.products == nil {
		ph.products = make([]*model.ItemProduct, 0)
	}
	ph.products = append(ph.products, pItemProduct)
}

type shoppingCartMemoryRepository struct {
	products []*model.Product              //lista de produtos
	cart     map[string]*productListHolder //map que contem o carrinho de usuarios
}

func (r *shoppingCartMemoryRepository) FindAllProducts() []*model.Product {
	lProducts := make([]*model.Product, 0)

	lProducts = append(lProducts, r.products...)

	return lProducts
}

func (r *shoppingCartMemoryRepository) FindProducById(id int32) *model.Product {
	var lProduct *model.Product = nil

	for _, p := range r.products {
		if p.Id == id {
			lProduct = p
			break
		}
	}

	return lProduct
}

func (r *shoppingCartMemoryRepository) FindGifts() []*model.Product {
	var lProducts []*model.Product = make([]*model.Product, 0)

	for _, p := range r.products {
		if p.IsGift {
			lProducts = append(lProducts, p)
		}
	}

	return lProducts
}

func (r *shoppingCartMemoryRepository) AddToCart(user *model.User, itemProduct *model.ItemProduct) error {

	listProductHolder := r.cart[user.Id]

	if listProductHolder == nil {
		listProductHolder = &productListHolder{products: make([]*model.ItemProduct, 0)}
		r.cart[user.Id] = listProductHolder
	}
	listProductHolder.Add(itemProduct)

	return nil
}

func (r *shoppingCartMemoryRepository) AddGiftToCart(user *model.User, itemProduct *model.ItemProduct, canAddGift func(itensProducts []*model.ItemProduct) bool) error {
	if itemProduct.IsGift {

		var productsInCart = make([]*model.ItemProduct, 0)
		prodPlaceholder := r.cart[user.Id]
		if prodPlaceholder != nil {
			//locando o carrinho do usuario para checagem
			if prodPlaceholder.mu == nil {
				prodPlaceholder.mu = &sync.Mutex{}
			}
			prodPlaceholder.mu.Lock()
			defer prodPlaceholder.mu.Unlock()

			if prodPlaceholder != nil && prodPlaceholder.products != nil && len(prodPlaceholder.products) > 0 {
				productsInCart = append(productsInCart, prodPlaceholder.products...)
			}
		}

		if canAddGift(productsInCart) {
			return r.AddToCart(user, itemProduct)
		}
	}

	return nil
}

func (r *shoppingCartMemoryRepository) ResumeCart(user *model.User) *model.CartResume {

	resumeCart := &model.CartResume{TotalAmount: 0, TotalDiscount: 0, TotalAmountWithDiscount: 0, Products: make([]*model.ItemProduct, 0)}

	cartPlaceHolder := r.cart[user.Id]

	if cartPlaceHolder != nil {

		for _, p := range cartPlaceHolder.products {
			resumeCart.TotalAmount += p.TotalAmount
			resumeCart.TotalDiscount += ((p.UnitAmount * p.Quantity) - p.TotalAmount)
			resumeCart.Products = append(resumeCart.Products, &model.ItemProduct{
				BaseProduct: p.BaseProduct,
				UnitAmount:  p.UnitAmount,
				Quantity:    p.Quantity,
				Discount:    p.Discount,
				TotalAmount: p.TotalAmount,
			})
		}
	}

	return resumeCart
}

func (r *shoppingCartMemoryRepository) EmptyCart(user *model.User) {
	delete(r.cart, user.Id)
}

func NewShoppingCartMemoryRepository() ShoppingCartRepositoryIf {
	return &shoppingCartMemoryRepository{
		products: utils.ReadJSONProducts(),
		cart:     make(map[string]*productListHolder),
	}
}
