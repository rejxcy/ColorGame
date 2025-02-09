package game

import (
	"encoding/json"
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

// handleWebSocket 處理 WebSocket 連接
func (c *controller) HandleWebSocket(ctx *gin.Context) {
	// 允許跨域
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true // 開發環境下允許所有來源
	}

	// 獲取參數
	roomID := ctx.Query("room_id")
	playerName := ctx.Query("player_name")
	isHost := ctx.Query("is_host") == "true"

	// 驗證參數
	if roomID == "" || playerName == "" {
		ctx.String(http.StatusBadRequest, "缺少必要參數")
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Output.Error("WebSocket upgrade failed: %v", err)
		return
	}

	// 記錄連接信息
	logger.Output.Info("New WebSocket connection: room=%s, player=%s, host=%v",
		roomID, playerName, isHost)

	// 獲取或創建房間
	var room *Room
	if isHost {
		room = c.createRoom(roomID)
		logger.Output.Info("Room %s created", room.ID)
	} else {
		room = c.getRoom(roomID)
		if room == nil {
			logger.Output.Error("Room %s not found", roomID)
			conn.WriteJSON(Message{
				Type:    msgTypeError,
				Payload: errRoomNotFound,
			})
			conn.Close()
			return
		}
		logger.Output.Info("Room %s found", room.ID)
	}

	// 創建玩家
	player := NewPlayer(conn, playerName, isHost)
	logger.Output.Info("Created new player: %s (host: %v)", player.Name, player.IsHost)

	// 將玩家加入房間
	if err := room.AddPlayer(player); err != nil {
		logger.Output.Error("Failed to add player %s to room %s: %v", player.Name, room.ID, err)
		conn.WriteJSON(Message{
			Type:    msgTypeError,
			Payload: err.Error(),
		})
		conn.Close()
		return
	}
	logger.Output.Info("Player %s joined room %s", player.Name, room.ID)
	logger.Output.Info("Room players: %v", room.Players)

	// 廣播更新後的玩家列表
	room.BroadcastPlayerList()

	// 開始處理玩家消息
	c.handlePlayerMessages(room, player)
}

// handlePlayerMessages 處理玩家消息
func (c *controller) handlePlayerMessages(room *Room, player *Player) {
	defer func() {
		if r := recover(); r != nil {
			logger.Output.Error("Panic in handlePlayerMessages: %v", r)
		}
		room.RemovePlayer(player.ID)
		logger.Output.Info("Player %s left room %s", player.Name, room.ID)
		room.BroadcastPlayerList()
		player.Conn.Close()

		if len(room.Players) == 0 {
			c.removeRoom(room.ID)
			logger.Output.Info("Room %s deleted", room.ID)
		}
	}()

	for {
		_, p, err := player.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				logger.Output.Error("WebSocket unexpected close: %v", err)
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(p, &msg); err != nil {
			logger.Output.Error("Failed to unmarshal message: %v", err)
			player.SendError(ErrCodeInvalidMessage, "無效的消息格式")
			continue
		}

		logger.Output.Info("Received message from player %s: type=%s, payload=%v",
			player.Name, msg.Type, msg.Payload)

		// 根據消息類型處理
		switch msg.Type {
		case MsgTypeStartGame:
			if err := room.HandleMessage(player.ID, msg); err != nil {
				logger.Output.Error("Failed to handle start game: %v", err)
				player.SendError(ErrCodeInvalidMessage, err.Error())
			}

		case MsgTypeReady:
			if err := room.HandleMessage(player.ID, msg); err != nil {
				logger.Output.Error("Failed to handle ready state: %v", err)
				player.SendError(ErrCodeInvalidMessage, err.Error())
			}
			logger.Output.Info("Player %s ready state changed to: %v", player.Name, msg.Payload)
			room.BroadcastPlayerList()

		default:
			logger.Output.Error("Unknown message type: %s", msg.Type)
			player.SendError(ErrCodeInvalidMessage, "未知的消息類型")
		}
	}
}

// createRoom 在 Context 中創建房間
func (c *controller) createRoom(id string) *Room {
	room := NewRoom(id)
	c.Base.GameRooms.Store(id, room)
	return room
}

// getRoom 從 Context 中獲取房間
func (c *controller) getRoom(id string) *Room {
	if room, ok := c.Base.GameRooms.Load(id); ok {
		return room.(*Room)
	}
	return nil
}

// removeRoom 從 Context 中移除房間
func (c *controller) removeRoom(id string) {
	c.Base.GameRooms.Delete(id)
}
