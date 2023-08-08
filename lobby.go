package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Lobby(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "lobby.html", nil)
}