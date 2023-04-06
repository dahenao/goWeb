package main

import (
	"os"

	"github.com/dahenao/goWeb/Package_Oriented_Design/cmd/api/handlers"
	"github.com/dahenao/goWeb/Package_Oriented_Design/docs"
	"github.com/dahenao/goWeb/Package_Oriented_Design/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func main() {
	godotenv.Load("./cmd/api/.env")
	// Generar un nuevo router en Gin.
	storeDb := store.NewStore("../products.json")

	server := gin.New() //engine de gin

	// Configurar el router.
	router := handlers.Router{
		Engine:  server,
		Storage: storeDb,
	}
	router.Setup()
	docs.SwaggerInfo.Host = os.Getenv("HOST")
	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Iniciar el servidor.
	server.Run(":8080")

}
