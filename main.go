package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moises-ba/ms-hash-shopping-cart/controller"
	"github.com/moises-ba/ms-hash-shopping-cart/repository"
	"github.com/moises-ba/ms-hash-shopping-cart/service/discountservice"
	"github.com/moises-ba/ms-hash-shopping-cart/service/holidayservice"
)

func main() {

	holidayservice := holidayservice.NewHolidayService()
	repo := repository.NewShoppingCartMemoryRepository()
	service := discountservice.NewDiscountService(holidayservice, repo)

	controller := controller.NewShoppingCartController(service)

	router := gin.Default()
	router.GET("/products", controller.ListProducts())
	router.POST("/checkout", controller.Checkout())

	router.Run("localhost:8080")

}
