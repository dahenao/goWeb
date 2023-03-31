package main

import "github.com/gin-gonic/gin"


func main(){
	router := gin.Default()  //crear router

	router.GET("/ping", func(c *gin.Context){  //endpint Handler
		//c.String(200, "v1: %s %s", c.Request.Method, c.Request.URL.Path)
		c.String(200,"pong") //respondemos un texto plano
	})

	router.Run()

}