package handlers

import (
	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/products"
	"github.com/dahenao/goWeb/Package_Oriented_Design/pkg/store"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine  *gin.Engine
	Storage store.Store
}

func (router *Router) Setup() {
	router.Engine.Use(gin.Logger()) //setea milddlewares por defecto
	router.Engine.Use(gin.Recovery())
	router.SetProductsRouter() //setear rutas

}

func (router *Router) SetProductsRouter() {
	storage := router.Storage
	repository := &products.Local_slice_DB{Storage: storage}

	service := products.ServiceDefault{
		BD: repository,
	}

	handler := ProductHandler{
		Service: service,
	}

	group := router.Engine.Group("/products")
	{
		group.Use(TokenMiddlewareValidate())
		group.POST("", handler.Create())
		group.GET("", handler.GetAll())
		group.GET(":id", handler.getProductByID())
		group.PUT(":id", handler.Update())
		group.PATCH(":id", handler.UpdatePartial())
		group.DELETE(":id", handler.Delete())
	}

}
