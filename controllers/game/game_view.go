package game

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GameView(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "game_view.html", nil)
}