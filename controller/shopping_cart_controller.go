package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moises-ba/ms-hash-shopping-cart/model"
	"github.com/moises-ba/ms-hash-shopping-cart/service"
)

type shoppingCartController struct {
	service service.ShoppingCartServiceIf
}

func NewShoppingCartController(pService service.ShoppingCartServiceIf) ShoppingCartControllerIf {
	return &shoppingCartController{
		service: pService,
	}
}

func (controller *shoppingCartController) ListProducts() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, controller.service.FindAllProducts())
	}
}

func (controller *shoppingCartController) Checkout() func(c *gin.Context) {
	return func(c *gin.Context) {

		productRequest := model.ProductsRequest{}
		if err := c.BindJSON(&productRequest); err != nil {
			log.Println(err)
			return
		}

		log.Println(len(productRequest.Products))

		usuarioSimulado := &model.User{Id: "fulano", Name: "Fulano"} //deveria vir de um token na vida real ou de uma base de dados.

		if len(productRequest.Products) > 0 {

			for _, v := range productRequest.Products {
				product := &model.Product{
					BaseProduct: model.BaseProduct{Id: v.Id},
				}

				err := controller.service.AddToCart(usuarioSimulado, &model.ItemProduct{
					BaseProduct: product.BaseProduct,
					Quantity:    v.Quantity,
				})

				if err != nil {
					log.Println(err.Error())
				}
			}

			cartResume := controller.service.ResumeCart(usuarioSimulado)
			//VERIFICAR SE DEVEMOS APAGAR O CARRINHO DE COMPRAS
			controller.service.EmptyCart(usuarioSimulado)

			c.IndentedJSON(http.StatusOK, cartResume)

		}
	}
}

func (controller *shoppingCartController) ResumeCart() func(c *gin.Context) {
	return nil
}
