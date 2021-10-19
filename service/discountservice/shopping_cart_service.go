package discountservice

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/moises-ba/ms-hash-shopping-cart/config"
	pb "github.com/moises-ba/ms-hash-shopping-cart/grpc/discount"
	"github.com/moises-ba/ms-hash-shopping-cart/model"
	"github.com/moises-ba/ms-hash-shopping-cart/repository"
	"github.com/moises-ba/ms-hash-shopping-cart/utils"
)

func NewDiscountService(pRepo repository.ShoppingCartRepositoryIf) DiscountServiceIf {
	return &discountService{repo: pRepo}
}

type discountService struct {
	repo repository.ShoppingCartRepositoryIf
}

func (s *discountService) FindDiscount(p *model.Product) (float32, error) {

	conn, err := utils.ConnectGRPCEndPoint(config.GetDiscountMSEndpoint())
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	client := pb.NewDiscountClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reponse, err := client.GetDiscount(ctx, &pb.GetDiscountRequest{ProductID: p.Id})
	if err != nil {
		return -1, err
	}

	return reponse.Percentage, nil

}

func (s *discountService) FindAllProducts() []*model.Product {
	return s.repo.FindAllProducts()
}

func (s *discountService) AddToCart(user *model.User, itemProduct *model.ItemProduct) error {

	discountProduct, err := s.FindDiscount(&itemProduct.Product)
	if err != nil {
		log.Println("Falha na chamada da api de descontos.", err)
		discountProduct = 0
	}

	itemProduct.Discount = discountProduct

	totalAmountItem := itemProduct.UnitAmount * itemProduct.Quantity //total dos produtos somados
	//calculando o desconto no valor total das somas dos itens
	itemProduct.TotalAmount = totalAmountItem - (totalAmountItem * (int32(discountProduct) / 100))

	//se o cliente escolheu comprar um presente ele deixa de ser presente
	itemProduct.IsGift = false

	err = s.repo.AddToCart(user, itemProduct)
	if err != nil {
		return err
	}

	dtBlackFriday, err := s.repo.FindBlackFridayDay()
	if err != nil {
		return err
	}
	dtNow := time.Now()

	if dtBlackFriday.Day() == dtNow.Day() && dtBlackFriday.Month() == dtNow.Month() { //eh blackfriday?
		gifts := s.repo.FindGifts()
		if len(gifts) > 0 {
			ramdomGift := gifts[rand.Intn(len(gifts)+1)]
			s.repo.AddGiftToCart(user, &model.ItemProduct{ //adiciona o presente
				Product:     *ramdomGift,
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
				return false
			})
		}
	}

	return nil
}

func (s *discountService) ResumeCart(user *model.User) *model.CartResume {
	return s.repo.ResumeCart(user)
}
