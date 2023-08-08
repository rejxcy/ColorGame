package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GameEntry(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "game_entry.html", nil)
}
