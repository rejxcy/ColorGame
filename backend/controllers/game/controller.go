package game

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rejxcy/colorgame/controllers"
	"github.com/rejxcy/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func New(base *controllers.Context) *controller {
	return &controller{Base: base}
}

type controller struct {
	Base *controllers.Context
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *controller) HandleWebSocket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Output.Error("WebSocket upgrade failed: %v", err)
		return
	}

	gameConn := &GameConnection{
		conn:     conn,
		game:     NewGame(),
		done:     make(chan struct{}),
		playerID: uuid.New().String(),
	}

	defer func() {
		close(gameConn.done)
		gameConn.conn.Close()
	}()

	// 啟動心跳檢測
	go gameConn.pingHandler()

	// 發送初始遊戲狀態
	gameConn.sendGameState()

	logger.Output.Info("player gameConn, Id: %v", gameConn.playerID)

	// 處理接收到的消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Output.Error("WebSocket error: %v", err)
			}
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			gameConn.sendError(1001, "無效的消息格式")
			logger.Output.Info("worng websocket msg type: %v", err)
			continue
		}

		gameConn.handleMessage(msg)
	}
}

func (gc *GameConnection) handleMessage(msg WSMessage) {
	switch msg.Type {
	case MsgTypeAnswer:
		color, ok := msg.Payload.(string)
		if !ok {
			gc.sendError(1002, "無效的答案格式")
			return
		}

		correct, isFinished, err := gc.game.Answer(color)
		if err != nil {
			gc.sendError(1003, err.Error())
			return
		}

		// 發送答案結果
		gc.sendMessage(MsgTypeAnswer, map[string]interface{}{
			"correct":  correct,
			"finished": isFinished,
		})

		// 更新遊戲狀態
		if isFinished {
			status, err := gc.game.GetStatus()
			if err != nil {
				gc.sendError(1004, "獲取遊戲狀態失敗")
				return
			}
			gc.sendMessage(MsgTypeGameOver, status)
		} else {
			status, err := gc.game.GetStatus()
			if err != nil {
				gc.sendError(1004, "獲取遊戲狀態失敗")
				return
			}
			gc.sendMessage(MsgTypeGameState, status)
		}

	case MsgTypeRestart:
		gc.game.Restart()
		status, err := gc.game.GetStatus()
		if err != nil {
			gc.sendError(1004, "獲取遊戲狀態失敗")
			return
		}
		gc.sendMessage(MsgTypeGameState, status)

	default:
		gc.sendError(1001, "未知的消息類型")
	}
}

func (gc *GameConnection) sendGameState() {
	status, err := gc.game.GetStatus()
	if err != nil {
		gc.sendError(1004, "獲取遊戲狀態失敗")
		logger.Output.Error("獲取遊戲狀態失敗: %v, %s", err, gc.game)
		return
	}
	gc.sendMessage(MsgTypeGameState, status)
}

func (gc *GameConnection) sendError(code int, message string) {
	gc.sendMessage(MsgTypeError, WSError{
		Code:    code,
		Message: message,
	})
}

func (gc *GameConnection) sendMessage(msgType string, payload interface{}) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	msg := WSMessage{
		Type:    msgType,
		Payload: payload,
	}

	if err := gc.conn.WriteJSON(msg); err != nil {
		logger.Output.Error("Write error: %v", err)
	}
}

func (gc *GameConnection) pingHandler() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			gc.mu.Lock()
			err := gc.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
			gc.mu.Unlock()
			if err != nil {
				logger.Output.Error("Ping error: %v", err)
				return
			}
		case <-gc.done:
			return
		}
	}
}
