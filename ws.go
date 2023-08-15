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

var clients = make(map[string]*websocket.Conn)
var player *Player
var game *Game
var isGameStart = false

func WS(ctx *gin.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	client := ctx.Param("client")

	clients[client] = ws
	defer delete(clients, client)

	if client != "waiting" {
		gamerReady(client)
		sendQuestion(ws)
	}

	for {
		_, m, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Error with read message:", err)
			return
		}

		if string(m) == "Restart" {
			game.restart()
			sendQuestionToAll()
			
		} else if game.isAnswer(string(m)) {
			if !game.isGameEnd {
				sendQuestionToAll()
			}
		}

		if game.isGameEnd {
			sendGameEndding(clients["Game"])
		}
	}
}

func sendQuestionToAll() {
	for _, ws := range clients {
		sendQuestion(ws)
	}
}

func sendQuestion(ws *websocket.Conn) {
	quiz, color := game.getQuiz()
	if quiz != "" && color != "" {
		message := fmt.Sprintf("%s,%s", quiz, color)
		if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			fmt.Println("Error sending question:", err)
		}
	}
}

func sendGameEndding(ws *websocket.Conn) {
	message := fmt.Sprintf("GameEnd,%d", game.wrongCount)
	if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		fmt.Println("Error sending question:", err)
	}
}

func gamerReady(playerName string) {
	if !isGameStart {
		player = NewPlayer(playerName)
		game = NewGame(player)
		if err := clients["waiting"].WriteMessage(websocket.TextMessage, []byte(playerName)); err != nil {
			fmt.Println("Error sending question:", err)
		}
	}
	isGameStart = true
}
