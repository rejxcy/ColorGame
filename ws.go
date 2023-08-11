package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var questions = []string{"紅色", "綠色", "藍色", "黃色", "橘色", "紫色"}
var colors = []string{"red", "green", "blue", "yellow", "orange", "purple"}

var clients = make(map[*websocket.Conn]bool)
var questionCount = 0

func WS(ctx *gin.Context) {

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	defer delete(clients, conn)

	sendQuestion(conn)

	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error with read message:", err)
			return
		}
		
		if string(m) == colors[questionCount] {
			NewQuestion()
			for c := range clients {
                sendQuestion(c)
            }
		} 
	}
}

func sendQuestion(conn *websocket.Conn) {
	message := fmt.Sprintf("%s,%s", questions[questionCount], colors[questionCount])
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		fmt.Println("Error sending question:", err)
	}
}

func NewQuestion() {
	//TODO: game logic
	questionCount ++
}
