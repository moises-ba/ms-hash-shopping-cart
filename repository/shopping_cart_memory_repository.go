package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/moises-ba/ms-hash-shopping-cart/model"
	"github.com/moises-ba/ms-hash-shopping-cart/utils"
)

var (
	products []*model.Product              //lista de produtos
	cart     map[string]*productListHolder //map que contem o carrinho de usuarios
)

//responsavel por encapsular uma lista de produtos e tratar concorrencia na insercao
type productListHolder struct {
	mu       *sync.Mutex
	products []*model.ItemProduct
}

func (ph *productListHolder) Add(pItemProduct *model.ItemProduct) {
	if ph.products == nil {
		ph.products = make([]*model.ItemProduct, 1)
	}
	ph.products = append(ph.products, pItemProduct)
}

func init() { //inicia as fontes de dados
	products = utils.ReadJSONProducts()
	cart = make(map[string]*productListHolder)
}

type shoppingCartMemoryRepository struct{}

func (r *shoppingCartMemoryRepository) FindAllProducts() []*model.Product {
	lProducts := make([]*model.Product, 0)

	lProducts = append(lProducts, products...)

	return lProducts
}

func (r *shoppingCartMemoryRepository) FindProducById(id int32) *model.Product {
	var lProduct *model.Product = nil

	for _, p := range products {
		if p.Id == id {
			lProduct = p
			break
		}
	}

	return lProduct
}

func (r *shoppingCartMemoryRepository) FindGifts() []*model.Product {
	var lProducts []*model.Product = make([]*model.Product, 0)

	for _, p := range products {
		if p.IsGift {
			lProducts = append(lProducts, p)
		}
	}

	return lProducts
}

func (r *shoppingCartMemoryRepository) AddToCart(user *model.User, itemProduct *model.ItemProduct) error {

	productFromDataSource := r.FindProducById(itemProduct.Id)
	if productFromDataSource == nil {
		return errors.New("Produto não encontrado.")
	}

	listProductHolder := cart[user.Id]

	if listProductHolder == nil {
		listProductHolder = &productListHolder{products: make([]*model.ItemProduct, 0)}
		cart[user.Id] = listProductHolder
	}
	listProductHolder.Add(itemProduct)

	return nil
}

func (r *shoppingCartMemoryRepository) AddGiftToCart(user *model.User, itemProduct *model.ItemProduct, canAddGift func(itensProducts []*model.ItemProduct) bool) error {
	if itemProduct.IsGift {

		var productsInCart = make([]*model.ItemProduct, 0)
		prodPlaceholder := cart[user.Id]
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

//como estamos usando um repositorio em memoria, retornamos de maneira fixa a data, deveria ser oriundo de uma base de feriados
func (r *shoppingCartMemoryRepository) FindBlackFridayDay() (time.Time, error) {
	layout := "02-01-2006"
	return time.Parse(layout, "26-11-20")
}

func (r *shoppingCartMemoryRepository) ResumeCart(user *model.User) *model.CartResume {

	resumeCart := &model.CartResume{TotalAmount: 0, TotalDiscount: 0, TotalAmountWithDiscount: 0, Products: make([]*model.ItemProduct, 10)}

	cartPlaceHolder := cart[user.Id]
	if cartPlaceHolder != nil {
		for _, p := range cartPlaceHolder.products {
			resumeCart.TotalAmount += p.UnitAmount * p.Quantity
			resumeCart.TotalDiscount += int32(p.Discount * float32(p.Quantity)) //ainda nao olhei a API para saber se o desconto chega como float
		}
	}

	return resumeCart
}

func NewShoppingCartMemoryRepository() ShoppingCartRepositoryIf {
	return &shoppingCartMemoryRepository{}
}
