package service

import (
	"errors"
	"testing"

	"github.com/moises-ba/ms-hash-shopping-cart/model"
	"github.com/moises-ba/ms-hash-shopping-cart/repository"
)

//holiday retornando true para blackFriday
type holidayServiceBlackFriday struct{}

func (h *holidayServiceBlackFriday) IsTodayBlackFriday() bool { return true }

type holidayServiceNotBlackFriday struct{}

func (h *holidayServiceNotBlackFriday) IsTodayBlackFriday() bool { return false }

//servico de desconto no ar
type discountServiceOK struct{}

func (s *discountServiceOK) FindDiscount(p *model.Product) (float32, error) { return 0.05, nil }

//servico de desconto Fora
type discountServiceNOK struct{}

func (s *discountServiceNOK) FindDiscount(p *model.Product) (float32, error) {
	return -1, errors.New("Servico de desconto fora")
}

//arquvio products.json de teste colocado no caminho relativo
var repo repository.ShoppingCartRepositoryIf = repository.NewShoppingCartMemoryRepository()

func TestFindAllProducts(t *testing.T) {

	products := NewShoppinCartService(&holidayServiceNotBlackFriday{}, &discountServiceOK{}, repo).FindAllProducts()

	expectedTotalReg := 6

	if len(products) != expectedTotalReg {
		t.Fatalf(`FindAllProducts esperava %v registros`, expectedTotalReg)
	}

}

func TestFindById(t *testing.T) {

	product := repo.FindProducById(1)

	if product == nil {
		t.Fatalf(`FindProducById esperava %v registros`, 1)
	}

}

func TestFindGifts(t *testing.T) {

	gifts := repo.FindGifts()

	if gifts == nil {
		t.Fatalf(`FindGifts esperava %v registros`, 1)
	}

	if gifts[0].IsGift == false {
		t.Fatalf(`FindGifts retornou um registro que não é um gift`)
	}

}

func TestAddToCartNotBlackFriday(t *testing.T) {

	user := &model.User{Id: "fulano"}
	repo.EmptyCart(user)
	item := &model.ItemProduct{BaseProduct: model.BaseProduct{Id: 1}, Quantity: 1}
	srv := NewShoppinCartService(&holidayServiceNotBlackFriday{}, &discountServiceNOK{}, repo)

	err := srv.AddToCart(user, item)
	if err != nil {
		t.Fatalf(`TestAddToCart falha ao adicionar ao carrinho -> %v `, err.Error())
	}

	cartResume := srv.ResumeCart(user)

	if cartResume == nil || len(cartResume.Products) != 1 {
		t.Fatal("Carrinho deveria conter 1 item")
	}

	if cartResume.Products[0].Id != 1 {
		t.Fatal("Id do produto  no carrinho deveria ser 1 ")
	}

}

func TestAddToCartNotBlackFridayWithDiscountActive(t *testing.T) {

	user := &model.User{Id: "fulano"}
	repo.EmptyCart(user)
	item := &model.ItemProduct{BaseProduct: model.BaseProduct{Id: 1}, Quantity: 1}
	srv := NewShoppinCartService(&holidayServiceNotBlackFriday{}, &discountServiceOK{}, repo)

	err := srv.AddToCart(user, item)
	if err != nil {
		t.Fatalf(`TestAddToCartNotBlackFridayWithDiscountActive falha ao adicionar ao carrinho -> %v `, err.Error())
	}

	cartResume := srv.ResumeCart(user)

	if cartResume == nil || len(cartResume.Products) != 1 {
		t.Fatal("TestAddToCartNotBlackFridayWithDiscountActive Carrinho deveria conter 1 item")
	}

	unitAmountExpected := int32(15157)
	totalAmountExpected := int32(15149)
	totalDiscountExpected := int32(8)

	if cartResume.Products[0].UnitAmount != unitAmountExpected ||
		cartResume.Products[0].TotalAmount != totalAmountExpected ||
		cartResume.TotalDiscount != totalDiscountExpected {

		t.Fatalf(`TestAddToCartNotBlackFridayWithDiscountActive Valores esperados UnitAmount %v, TotalAmount %v, TotalDiscount %v `,
			unitAmountExpected, totalAmountExpected, totalDiscountExpected)
	}

}

func TestAddToCartBlackFriday(t *testing.T) {

	user := &model.User{Id: "fulano"}
	repo.EmptyCart(user)
	item := &model.ItemProduct{BaseProduct: model.BaseProduct{Id: 1}, Quantity: 1}
	srv := NewShoppinCartService(&holidayServiceBlackFriday{}, &discountServiceNOK{}, repo)

	err := srv.AddToCart(user, item)
	if err != nil {
		t.Fatalf(`TestAddToCartBlackFriday falha ao adicionar ao carrinho -> %v `, err.Error())
	}

	cartResume := srv.ResumeCart(user)

	if cartResume == nil || len(cartResume.Products) != 2 {
		t.Fatal("TestAddToCartBlackFriday Carrinho deveria conter 2 itens")
	}

	if cartResume.Products[0].IsGift != false || cartResume.Products[1].IsGift != true {
		t.Fatal("TestAddToCartBlackFriday Carrinho deveria ter dois produtos, um comprado e um presente")
	}

}
