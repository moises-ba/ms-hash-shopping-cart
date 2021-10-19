package controller

import "github.com/gin-gonic/gin"

type ShoppingCartControllerIf interface {
	ListProducts() func(c *gin.Context)
	Checkout() func(c *gin.Context)
	ResumeCart() func(c *gin.Context)
}
