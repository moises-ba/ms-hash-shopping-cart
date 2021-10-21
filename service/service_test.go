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

func (h *holidayServiceNotBlackFriday) IsTodayBlackFriday() bool { return true }

//servico de desconto no ar
type discountServiceOK struct{}

func (s *discountServiceOK) FindDiscount(p *model.Product) (float32, error) { return 0.5, nil }

//servico de desconto Fora
type discountServiceNOK struct{}

func (s *discountServiceNOK) FindDiscount(p *model.Product) (float32, error) {
	return -1, errors.New("Servico de desconto fora")
}

//arquvio products.json de teste colocado no caminho relativo
var repo repository.ShoppingCartRepositoryIf = repository.NewShoppingCartMemoryRepository()

func TestFindAllProducts(t *testing.T) {

	products := NewShoppinCartService(&holidayServiceBlackFriday{}, &discountServiceOK{}, repo).FindAllProducts()

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
