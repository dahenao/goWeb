package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	ErrEmptyField = errors.New("campo vacio")
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

type Request struct {
	Name         string  `json:"name" validate:"required"` //required no permite que ese envien valores por defecto desde la peticion
	Quantity     int     `json:"quantity" validate:"required"`
	Code_value   string  `json:"code_value" validate:"required"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration" validate:"required"`
	Price        float64 `json:"price" validate:"required"`
}

func (p *Product) validateCode() error {

	for _, v := range Products {
		if v.Code_value == p.Code_value {
			fmt.Errorf("Invalid product code ")
		}

	}
	return nil
}

func (p *Product) validate() error {
	if p.Id == 0 {
		return fmt.Errorf("no envio el campo Id : %w", ErrEmptyField)
	}
	if p.Name == "" {
		return fmt.Errorf("no envio el campo Name : %w", ErrEmptyField)
	}
	if p.Quantity == 0 {
		return fmt.Errorf("no envio el campo Code_value : %w", ErrEmptyField)
	}
	if p.Code_value == "" {
		return fmt.Errorf("no envio el campo Code_value: %w", ErrEmptyField)
	}
	return nil
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

func Guardar() gin.HandlerFunc {

	return func(c *gin.Context) { //endpint Handler

		var prdRequest Request

		if err := c.ShouldBindJSON(&prdRequest); err != nil {
			c.JSON(400, err)
		}

		prd := &Product{
			Id:           LastId + 1,
			Name:         prdRequest.Name,
			Quantity:     prdRequest.Quantity,
			Code_value:   prdRequest.Code_value,
			Is_published: prdRequest.Is_published,
			Expiration:   prdRequest.Expiration,
			Price:        prdRequest.Price}

		if err := prd.validate(); err != nil {
			if errors.Is(err, ErrEmptyField) {
				c.JSON(http.StatusBadRequest, gin.H{"message": "invalid product", "data": err.Error()})
				return
			}
		}

		if err := prd.validateCode(); err != nil {
			if errors.Is(err, ErrEmptyField) {
				c.JSON(http.StatusBadRequest, gin.H{"message": "invalid product", "data": err.Error()})
				return
			}
		}

		Products = append(Products, *prd)

		c.JSON(200, prd) //respondemos un json decodificando el slice e products
		LastId++
	}
}

var LastId int

func main() {

	router := gin.Default() //crear router

	router.POST("/products", Guardar()) //puede recibir un elipsis de handlers (middlewarez)

	router.GET("/products", func(c *gin.Context) { //endpint Handler
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
