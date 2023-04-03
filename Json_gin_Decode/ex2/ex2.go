package main

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Client struct {
	Nombre   string
	Apellido string
}

func main() {
	router := gin.Default() //crear router

	router.POST("/saludo", func(c *gin.Context) { //endpint Handler

		var cl Client

		//if err := c.ShouldBindJSON(&cl); err != nil {
		//json.NewDecoder(c.Request.Body).Decode(&cl)
		//if  errm := json.Unmarshal([]byte (c.Request.body),&cl); errm != nil {
		if errD := json.NewDecoder(c.Request.Body).Decode(&cl); errD != nil { //decodifica el json en la estructura Cliente
			c.JSON(400, gin.H{"error": errD})
			return
		}

		c.JSON(200, gin.H{"answer": "Hola " + cl.Nombre + cl.Apellido})
	})

	router.Run()

}
