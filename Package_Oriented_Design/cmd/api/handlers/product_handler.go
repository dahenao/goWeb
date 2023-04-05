package handlers

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/domain"
	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/products"
	"github.com/gin-gonic/gin"
)

var (
	ErrBadIndex     = errors.New("Error: index is bad")
	ErrInvalidToken = errors.New("invalid token")
)

type RequestProduct struct {
	Name         string  `json:"name" validate:"required"` //required no permite que ese envien valores por defecto desde la peticion
	Quantity     int     `json:"quantity" validate:"required"`
	Code_value   string  `json:"code_value" validate:"required"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration" validate:"required"`
	Price        float64 `json:"price" validate:"required"`
}

func (request RequestProduct) ToDomain() domain.Product {

	return domain.Product{
		Id:           0,
		Name:         request.Name,
		Quantity:     request.Quantity,
		Code_value:   request.Code_value,
		Is_published: request.Is_published,
		Expiration:   request.Expiration,
		Price:        request.Price}
}

type ProductHandler struct {
	Service products.Service
}

func (handler *ProductHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if ctx.GetHeader("token") != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, ErrInvalidToken.Error())
			return
		}
		var request RequestProduct //creamos variable de tipo request product

		if err := ctx.ShouldBindJSON(&request); err != nil { //recibimos el body y lo codificamos a json en la variable creada
			ctx.JSON(400, err)
			return
		}
		productToCreate := request.ToDomain() //llamamos el metodo de requestProduct que convierte un request a domain product

		//enviamos el producto a la capa de servicio
		if err := handler.Service.Create(&productToCreate); err != nil {
			if errors.Is(err, products.ErrProductAlreadyExists) { //ErrProductAlreadyExists es una variable global en products
				ctx.JSON(400, err.Error()) //retorna el error que ya fue discriminado
			} else {
				ctx.JSON(http.StatusInternalServerError, "an internal error has occurred")
			}
			return

		}

		ctx.JSON(http.StatusCreated, productToCreate)
	}
}

func (handler *ProductHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Obtener los productos.
		products, err := handler.Service.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "an internal error has occurred")
			return
		}

		// Devolver la respuesta.
		ctx.JSON(http.StatusOK, products)
	}
}

func (handler *ProductHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var request RequestProduct                  //creamos variable de tipo request product
		index, err := strconv.Atoi(ctx.Param("id")) //convertimos el parametro a entero
		if err != nil {
			ctx.JSON(400, ErrBadIndex)
		}
		if err := ctx.ShouldBindJSON(&request); err != nil { //recibimos el body y lo codificamos a json en la variable creada
			ctx.JSON(400, err)
			return
		}
		productToCreate := request.ToDomain() //llamamos el metodo de requestProduct que convierte un request a domain product

		//enviamos el producto a la capa de servicio
		if err := handler.Service.Update(index, &productToCreate); err != nil {
			errors.Is(err, products.ErrProductAlreadyExists) //ErrProductAlreadyExists es una variable global en products
			ctx.JSON(400, err.Error())                       //retorna el error que ya fue discriminado
			return

		}

		ctx.JSON(http.StatusOK, productToCreate)
	}
}

func (handler *ProductHandler) UpdatePartial() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//creamos variable de tipo request product
		index, err := strconv.Atoi(ctx.Param("id")) //convertimos el parametro a entero
		if err != nil {
			ctx.JSON(400, ErrBadIndex)
		}

		prod, err := handler.Service.GetProductByID(index)
		if err != nil {
			ctx.JSON(400, err)
		}

		if err := ctx.ShouldBindJSON(&prod); err != nil { //recibimos el body y lo codificamos a el producto consultado
			ctx.JSON(400, err)
			return
		}
		prod.Id = index
		//enviamos el producto a la capa de servicio
		if err := handler.Service.Update(index, &prod); err != nil {
			errors.Is(err, products.ErrProductAlreadyExists) //ErrProductAlreadyExists es una variable global en products
			ctx.JSON(400, err.Error())                       //retorna el error que ya fue discriminado
			return

		}

		ctx.JSON(http.StatusOK, prod)
	}
}

func (handler *ProductHandler) getProductByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Obtener los productos.
		index, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, err)
		}
		product, err := handler.Service.GetProductByID(index)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "an internal error has occurred")
			return
		}

		// Devolver la respuesta.
		ctx.JSON(http.StatusOK, product)
	}
}

func (handler *ProductHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Obtener los productos.
		index, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, ErrBadIndex)
		}
		err = handler.Service.Delete(index)
		if err != nil {

			ctx.JSON(http.StatusInternalServerError, "an internal error has occurred")
			return
		}

		// Devolver la respuesta.
		ctx.JSON(http.StatusOK, "deleted product")
	}
}
