package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PlayerView(ctx *gin.Context) {
	playerName := ctx.PostForm("playerName")
	
	ctx.HTML(http.StatusOK, "player_view.html", gin.H{
		"playerName": playerName,
	})
}