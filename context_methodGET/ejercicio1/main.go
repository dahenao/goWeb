package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id           int
	Name         string
	Quantity     int
	Code_value   string
	Is_published bool
	Expiration   string
	Price        float64
}

var Products []Product

func getProductByID(id int) (Product, error) {

	for _, p := range Products {
		if id == p.Id {
			return p, nil
		}
	}
	return Product{}, fmt.Errorf("no exite el prodructo con id: %d", id)
}

func getPricesGT(price float64) []Product {
	var prices []Product
	for _, p := range Products {
		if price < p.Price {
			prices = append(prices, p)
		}
	}
	return prices
}

func createDatabase(path string) {
	fileJSON, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fileJSON.Close()

	err = json.NewDecoder(fileJSON).Decode(&Products) //lee el archivo json y lo decodifica en el slice products
	if err != nil {
		panic(err)
	}
}

func main() {

	createDatabase("products.json")
	router := gin.Default() //crear router

	router.GET("/ping", func(c *gin.Context) { //endpint Handler
		c.String(200, "pong") //respondemos un texto plano
	})

	router.GET("/products", func(c *gin.Context) { //endpint Handler

		//json.NewDecoder(fileJSON).Decode(&Products)
		c.JSON(200, Products) //respondemos un json decodificando el slice e products
	})

	router.GET("/products/:id", func(c *gin.Context) { //endpint Handler
		id, err := strconv.Atoi(c.Param("id")) //convertimos el parametro a entero
		if err != nil {
			fmt.Println(err)
		}
		prd, err := getProductByID(id) //la funcion debuelve un producto
		if err != nil {
			fmt.Println(err)
		}
		//json.NewDecoder(fileJSON).Decode(&Products)
		c.JSON(200, prd) //respondemos un json decodificando con el producto
	})

	router.GET("/products/search", func(c *gin.Context) { //endpint Handler
		price, err := strconv.ParseFloat(c.Query("priceGt"), 64) //convertimos el parametro a entero
		if err != nil {
			fmt.Println(err)
		}
		prd := getPricesGT(price) //la funcion debuelve un producto
		c.JSON(200, prd)          //respondemos un json decodificando con el producto
	})

	router.Run()

	/*
	    Crear una ruta /products/:id que nos devuelva un producto por su id.
	Crear una ruta /products/search que nos permita buscar por parámetro los productos cuyo precio sean mayor a un valor priceGt.*/

}
