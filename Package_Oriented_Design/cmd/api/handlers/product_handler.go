package handlers

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/domain"
	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/products"
	"github.com/dahenao/goWeb/Package_Oriented_Design/pkg/web"
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

func tokenValidate(token string) bool {
	return os.Getenv("TOKEN") == token
}

type ProductHandler struct {
	Service products.Service
}

// Post godoc
// @Summary      Create a new product
// @Description  Create a new product in repository
// @Tags         products
// @Produce      json
// @Param        token header string true "token"
// @Param        product body domain.Product true "Product"
// @Router       /products [post]
func (handler *ProductHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if !tokenValidate(ctx.GetHeader("token")) {
			ctx.JSON(http.StatusUnauthorized, ErrInvalidToken.Error())
			return
		}
		var request RequestProduct //creamos variable de tipo request product

		if err := ctx.ShouldBindJSON(&request); err != nil { //recibimos el body y lo codificamos a json en la variable creada
			web.ErrorResp(ctx, err, http.StatusBadRequest)
			//ctx.JSON(400, web.ErrorRespose{Status: "Error", Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		productToCreate := request.ToDomain() //llamamos el metodo de requestProduct que convierte un request a domain product

		//enviamos el producto a la capa de servicio
		if err := handler.Service.Create(&productToCreate); err != nil {
			if errors.Is(err, products.ErrProductAlreadyExists) { //ErrProductAlreadyExists es una variable global en products
				//ctx.JSON(400, web.ErrorRespose{Status: "Error", Code: http.StatusBadRequest, Message: err.Error()}) //retorna el error que ya fue discriminado err.Error()
				web.ErrorResp(ctx, err, http.StatusBadRequest)
			} else {
				//ctx.JSON(http.StatusInternalServerError, "an internal error has occurred")
				web.ErrorResp(ctx, err, http.StatusInternalServerError)
			}
			return

		}
		web.OkResp(ctx, http.StatusCreated, productToCreate)
		//ctx.JSON(http.StatusCreated, productToCreate)
	}
}

func (handler *ProductHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !tokenValidate(ctx.GetHeader("token")) {
			ctx.JSON(http.StatusUnauthorized, ErrInvalidToken.Error())
			return
		}
		// Obtener los productos.
		products, err := handler.Service.GetAll()
		if err != nil {
			web.ErrorResp(ctx, err, http.StatusInternalServerError)
			//ctx.JSON(http.StatusInternalServerError, web.ErrorRespose{Status: "Error", Code: http.StatusInternalServerError, Message: "an internal error has occurred"})
			return
		}

		// Devolver la respuesta.
		web.OkResp(ctx, http.StatusOK, products)
		//ctx.JSON(http.StatusOK, products)
	}
}

func (handler *ProductHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !tokenValidate(ctx.GetHeader("token")) {
			ctx.JSON(http.StatusUnauthorized, ErrInvalidToken.Error())

			return
		}
		var request RequestProduct                  //creamos variable de tipo request product
		index, err := strconv.Atoi(ctx.Param("id")) //convertimos el parametro a entero
		if err != nil {
			//ctx.JSON(400, ErrBadIndex)

			web.ErrorResp(ctx, ErrBadIndex, http.StatusBadRequest)
		}
		if err := ctx.ShouldBindJSON(&request); err != nil { //recibimos el body y lo codificamos a json en la variable creada
			ctx.JSON(400, err)
			return
		}
		productToCreate := request.ToDomain() //llamamos el metodo de requestProduct que convierte un request a domain product

		//enviamos el producto a la capa de servicio
		if err := handler.Service.Update(index, &productToCreate); err != nil {
			//errors.Is(err, products.ErrProductAlreadyExists) //ErrProductAlreadyExists es una variable global en products
			//ctx.JSON(400, err.Error())                       //retorna el error que ya fue discriminado
			web.ErrorResp(ctx, products.ErrProductAlreadyExists, http.StatusBadRequest)
			return

		}

		//ctx.JSON(http.StatusOK, productToCreate)
		web.OkResp(ctx, http.StatusOK, productToCreate)
	}
}

func (handler *ProductHandler) UpdatePartial() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !tokenValidate(ctx.GetHeader("token")) {
			ctx.JSON(http.StatusUnauthorized, ErrInvalidToken.Error())
			return
		}
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
		if !tokenValidate(ctx.GetHeader("token")) {
			ctx.JSON(http.StatusUnauthorized, ErrInvalidToken.Error())
			return
		}
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
		if !tokenValidate(ctx.GetHeader("token")) {
			ctx.JSON(http.StatusUnauthorized, ErrInvalidToken.Error())
			return
		}
		// Obtener los productos.
		index, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, ErrBadIndex)
		}
		err = handler.Service.Delete(index)
		if err != nil {
			web.ErrorResp(ctx, errors.New("an internal error has occurred"), http.StatusInternalServerError)
			//ctx.JSON(http.StatusInternalServerError, "an internal error has occurred")
			return
		}

		// Devolver la respuesta.
		//ctx.JSON(http.StatusOK, "deleted product")
		web.OkResp(ctx, http.StatusOK, "deleted product")
	}
}
