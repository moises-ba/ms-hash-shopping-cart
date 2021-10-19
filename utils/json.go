package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/moises-ba/ms-hash-shopping-cart/model"
)

type Products struct {
	Products []*model.Product `json:"products"`
}

func ReadJSONProducts() []*model.Product {

	jsonFile, err := os.Open("products.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var products Products
	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		panic(err)
	}

	return products.Products

}
