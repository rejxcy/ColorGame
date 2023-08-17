package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db *playerDB

func main() {
	db = initialPlayerDb()

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

func initialPlayerDb() *playerDB {
	db, err := NewPlayerDB()
	if err != nil {
		log.Fatalln("Cann't initial playerDB, err: ", err)
	}
	return db
}

func CheckPlayer(name string) *Player {
	player, err := db.selectPlayerByName(name)
	if err != nil {
		fmt.Println("player not found")
	}

	if player == nil {
		player := Player{
			name:       name,
			timeRecord: 0,
		}
		db.insert(player)
	}
	
	return player
}

func UpdatePlayerRecord(player *Player) {
	db.update(*player)
}
