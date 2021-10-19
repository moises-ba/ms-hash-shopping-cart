package config

import "github.com/moises-ba/ms-hash-shopping-cart/utils"

const (
	MS_DISCOUNT_URL = "MS_DISCOUNT_URL"
)

func GetDiscountMSEndpoint() string {
	return utils.Getenv(MS_DISCOUNT_URL, "localhost:50051")
}
