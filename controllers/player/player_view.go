package player

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PlayerView(ctx *gin.Context) {
	localIP, err := GetLocalIP()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error getting local IP")
		return
	}
	playerName := ctx.PostForm("playerName")
	
	ctx.HTML(http.StatusOK, "player_view.html", gin.H{
		"localIP": localIP,
		"playerName": playerName,
	})
}