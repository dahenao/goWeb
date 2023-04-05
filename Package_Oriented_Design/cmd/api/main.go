package main

import (
	"github.com/dahenao/goWeb/Package_Oriented_Design/cmd/api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("./cmd/api/.env")
	// Generar un nuevo router en Gin.

	server := gin.New()

	// Configurar el router.
	router := handlers.Router{
		Engine: server,
	}
	router.Setup()

	// Iniciar el servidor.
	server.Run(":8080")

}
