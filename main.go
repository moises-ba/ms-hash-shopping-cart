package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moises-ba/ms-hash-shopping-cart/controller"
	"github.com/moises-ba/ms-hash-shopping-cart/repository"
	"github.com/moises-ba/ms-hash-shopping-cart/service"
)

func main() {

	holidayservice := service.NewHolidayService()
	discountService := service.NewDiscountService()

	repo := repository.NewShoppingCartMemoryRepository()

	service := service.NewShoppinCartService(holidayservice, discountService, repo)

	controller := controller.NewShoppingCartController(service)

	router := gin.Default()
	router.GET("/products", controller.ListProducts())
	router.POST("/checkout", controller.Checkout())

	router.Run(":8080")

}
