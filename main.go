package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moises-ba/ms-hash-shopping-cart/controller"
	"github.com/moises-ba/ms-hash-shopping-cart/repository"
	"github.com/moises-ba/ms-hash-shopping-cart/service/discountservice"
)

func main() {

	repo := repository.NewShoppingCartMemoryRepository()
	service := discountservice.NewDiscountService(repo)

	controller := controller.NewShoppingCartController(service)

	router := gin.Default()
	router.GET("/products", controller.ListProducts())
	router.POST("/checkout", controller.Checkout())

	router.Run("localhost:8080")

}
