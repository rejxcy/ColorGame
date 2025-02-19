package game

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rejxcy/colorgame/backend/controllers"
	"github.com/rejxcy/logger"
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

// 封裝錯誤回應（發送錯誤消息並安全關閉連線）
func sendErrorAndClose(conn *websocket.Conn, err error) {
	if conn != nil {
		if writeErr := conn.WriteJSON(Message{
			Type:    MsgTypeError,
			Payload: err.Error(),
		}); writeErr != nil {
			logger.Output.Error("Error writing error message: %v", writeErr)
		}
		conn.Close()
	}
}

// 處理 WebSocket 連線的建立與參數驗證
func (c *controller) HandleWebSocket(ctx *gin.Context) {
	// 允許跨域，開發環境下允許所有來源
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// 取得必要參數
	roomID := ctx.Query("room_id")
	playerName := ctx.Query("player_name")
	isHost := ctx.Query("is_host") == "true"

	// 若缺少參數則返回錯誤
	if roomID == "" || playerName == "" {
		ctx.String(http.StatusBadRequest, "缺少必要參數")
		return
	}

	// 升級為 WebSocket 連線
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Output.Error("WebSocket upgrade failed: %v", err)
		return
	}

	logger.Output.Info("New WebSocket connection: room=%s, player=%s, host=%v", roomID, playerName, isHost)

	// 根據玩家身分決定創建或獲取房間
	var room *Room
	if isHost {
		room = c.createRoom(roomID)
		logger.Output.Info("Room %s created", room.ID)
	} else {
		room = c.getRoom(roomID)
		if room == nil {
			logger.Output.Error("Room %s not found", roomID)
			sendErrorAndClose(conn, errors.New("Room not found"))
			return
		}
		logger.Output.Info("Room %s found", room.ID)
	}

	// 建立玩家並加入房間
	player := NewPlayer(conn, playerName, isHost)
	logger.Output.Info("Created new player: %s (host: %v)", player.Name, player.IsHost)

	if err := room.AddPlayer(player); err != nil {
		logger.Output.Error("Failed to add player %s to room %s: %v", player.Name, room.ID, err)
		sendErrorAndClose(conn, err)
		return
	}
	logger.Output.Info("Player %s joined room %s", player.Name, room.ID)
	room.BroadcastPlayerList()

	// 進入持續接收並分發玩家消息的循環
	c.handlePlayerMessages(room, player)
}

// 持續接收玩家消息並委由玩家處理
func (c *controller) handlePlayerMessages(room *Room, player *Player) {
	defer func() {
		if r := recover(); r != nil {
			logger.Output.Error("Panic in handlePlayerMessages: %v", r)
		}
		room.RemovePlayer(player.ID)
		logger.Output.Info("Player %s left room %s", player.Name, room.ID)
		room.BroadcastPlayerList()
		player.Close()

		if len(room.Players) == 0 {
			c.removeRoom(room.ID)
			logger.Output.Info("Room %s deleted", room.ID)
		}
	}()

	for {
		_, messageData, err := player.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Output.Error("WebSocket unexpected close: %v", err)
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(messageData, &msg); err != nil {
			logger.Output.Error("Failed to unmarshal message: %v", err)
			player.SendError(ErrCodeInvalidMessage, "無效的消息格式")
			continue
		}

		logger.Output.Info("Received message from player %s: type=%s, payload=%v", player.Name, msg.Type, msg.Payload)
		if err := player.HandleMessage(msg, room); err != nil {
			logger.Output.Error("Error processing message for player %s: %v", player.Name, err)
			player.SendError(ErrCodeInvalidMessage, err.Error())
		}
	}
}

// 在 Context 中創建房間
func (c *controller) createRoom(id string) *Room {
	room := NewRoom(id)
	c.Base.GameRooms.Store(id, room)
	return room
}

// 從 Context 中獲取房間
func (c *controller) getRoom(id string) *Room {
	if room, ok := c.Base.GameRooms.Load(id); ok {
		return room.(*Room)
	}
	return nil
}

// 從 Context 中移除房間
func (c *controller) removeRoom(id string) {
	c.Base.GameRooms.Delete(id)
}
