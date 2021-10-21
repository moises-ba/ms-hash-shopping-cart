package service

import (
	"context"
	"time"

	"github.com/moises-ba/ms-hash-shopping-cart/config"
	pb "github.com/moises-ba/ms-hash-shopping-cart/grpc/discount"
	"github.com/moises-ba/ms-hash-shopping-cart/model"
	"github.com/moises-ba/ms-hash-shopping-cart/utils"
)

type discountService struct {
}

func NewDiscountService() DiscountServiceIf {
	return &discountService{}
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
