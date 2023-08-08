package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("client/*")
	fmt.Println("Server Start!")
	router.GET("/lobby", Lobby)
	router.GET("/game_entry", GameEntry)
	router.GET("/gamer_name", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "gamer_name.html", nil)
	})
	router.POST("/player_view", PlayerView)

	router.Run()
}
