package handlers

import (
	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/products"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func (router *Router) Setup() {
	router.Engine.Use(gin.Logger()) //setea milddlewares por defecto
	router.Engine.Use(gin.Recovery())
	router.SetProductsRouter() //setear rutas

}

func (router *Router) SetProductsRouter() {
	repository := &products.Local_slice_DB{}

	service := products.ServiceDefault{
		BD: repository,
	}

	handler := ProductHandler{
		Service: service,
	}

	group := router.Engine.Group("/products")
	group.POST("", handler.Create())
	group.GET("", handler.GetAll())
	group.GET(":id", handler.getProductByID())
	group.PUT(":id", handler.Update())
	group.PATCH(":id", handler.UpdatePartial())
	group.DELETE(":id", handler.Delete())

}
