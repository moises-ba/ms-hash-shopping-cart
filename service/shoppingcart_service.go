package service

import (
	"errors"
	"log"
	"math/rand"
	"strconv"

	"github.com/moises-ba/ms-hash-shopping-cart/model"
	"github.com/moises-ba/ms-hash-shopping-cart/repository"
)

func NewShoppinCartService(pHolidayservice HolidayServiceIf, pDiscountService DiscountServiceIf, pRepo repository.ShoppingCartRepositoryIf) ShoppingCartServiceIf {
	return &shoppingCartService{repo: pRepo, holidayservice: pHolidayservice}
}

type shoppingCartService struct {
	holidayservice   HolidayServiceIf
	pDiscountService DiscountServiceIf
	repo             repository.ShoppingCartRepositoryIf
}

func (s *shoppingCartService) FindAllProducts() []*model.Product {
	return s.repo.FindAllProducts()
}

func (s *shoppingCartService) AddToCart(user *model.User, itemProduct *model.ItemProduct) error {

	productFromDataSource := s.repo.FindProducById(itemProduct.Id) //obtendo o produto na base de dados
	if productFromDataSource == nil {
		return errors.New("produto não encontrado: " + strconv.Itoa(int(itemProduct.Id)))
	}
	itemProduct.BaseProduct = productFromDataSource.BaseProduct
	itemProduct.UnitAmount = productFromDataSource.Amount

	discountProduct, err := s.pDiscountService.FindDiscount(productFromDataSource)
	if err != nil {
		log.Println("Falha na chamada da api de descontos.", err)
		discountProduct = 0
	}

	log.Println("Desconto concedido: ", discountProduct)

	itemProduct.Discount = discountProduct

	totalAmountItem := itemProduct.UnitAmount * itemProduct.Quantity //total dos produtos somados
	//calculando o desconto no valor total das somas dos itens
	itemProduct.TotalAmount = int32(float32(totalAmountItem) - (float32(totalAmountItem) * (discountProduct / 100.0)))

	//se o cliente escolheu comprar um presente ele deixa de ser presente
	itemProduct.IsGift = false

	err = s.repo.AddToCart(user, itemProduct)
	if err != nil {
		return err
	}

	if s.holidayservice.IsTodayBlackFriday() {
		gifts := s.repo.FindGifts()
		if len(gifts) > 0 {
			ramdomGift := gifts[rand.Intn(len(gifts))]
			s.repo.AddGiftToCart(user, &model.ItemProduct{ //adiciona o presente
				BaseProduct: ramdomGift.BaseProduct,
				TotalAmount: 0,
				Discount:    0,
				Quantity:    1,
			}, func(itensInCart []*model.ItemProduct) bool { //lambda q o repositorio utilizara para verificar se pode inserir ou nao um presente
				if len(itensInCart) == 0 { // o carrinho nao pode estar vazio
					return false
				}
				for _, v := range itensInCart { //nao pode haver presentes no carrinho
					if v.IsGift {
						return false
					}
				}
				return true
			})
		}
	}

	return nil
}

func (s *shoppingCartService) ResumeCart(user *model.User) *model.CartResume {

	cartResume := s.repo.ResumeCart(user)

	return cartResume
}

func (s *shoppingCartService) EmptyCart(user *model.User) {
	s.repo.EmptyCart(user)
}
