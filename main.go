package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("client/*")
	fmt.Println("Server Start!")
	router.GET("/lobby", Lobby)

	router.Run()
}
