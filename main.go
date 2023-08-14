package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("client/*")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})
	router.GET("/player_name", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "player_name.html", nil)
	})
	router.GET("/game_entry", GameEntry)
	router.GET("/game_view", GameView)
	router.GET("/ws/:client", WS)

	router.POST("/player_view", PlayerView)

	router.Run()
}

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}
	return "", fmt.Errorf("no valid local ip found")
}