package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Lobby(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "lobby.html", nil)
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("client/*")
	fmt.Println("Server Start!")
	router.GET("/lobby", Lobby)

	router.Run()
}
