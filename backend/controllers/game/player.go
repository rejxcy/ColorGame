package game

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rejxcy/logger"
)

// NewPlayer 創建新玩家並初始化資料
func NewPlayer(conn *websocket.Conn, name string, isHost bool) *Player {
	return &Player{
		ID:      uuid.New().String(),
		Name:    name,
		IsHost:  isHost,
		IsReady: false,
		Score:   0,
		Conn:    conn,
		Game:    NewGame(),
	}
}

// Send 透過 WriteJSON 發送消息給玩家前先檢查連線是否有效
func (p *Player) Send(msg Message) error {
	if p.Conn == nil {
		logger.Output.Error("連線為 nil")
		return errors.New("連線為 nil")
	}
	return p.Conn.WriteJSON(msg)
}

// SendError 發送錯誤消息給客戶端
func (p *Player) SendError(code string, message string) {
	p.Send(Message{
		Type: MsgTypeError,
		Payload: map[string]interface{}{
			"code":    code,
			"message": message,
		},
	})
}

// HandleMessage 根據消息類型分派處理（所有消息均由玩家處理）
func (p *Player) HandleMessage(msg Message, room *Room) error {
	switch msg.Type {
	case MsgTypeReady:
		// 房主不需要設置準備狀態
		if p.IsHost {
			return errors.New("房主無需設置準備狀態")
		}
		ready, ok := msg.Payload.(bool)
		if !ok {
			return errors.New("無效的準備狀態")
		}
		p.IsReady = ready
		logger.Output.Info("Player %s ready state changed to: %v", p.Name, ready)
		room.BroadcastPlayerList()
		return nil

	case MsgTypeGameStart:
		return p.handleGameStart(room)

	case MsgTypeAnswer:
		return p.handleAnswer(msg.Payload, room)

	case MsgTypeGameReset:
		return p.handleGameReset(room)

	default:
		return errors.New("未知的消息類型")
	}
}

// handleGameStart 處理開始遊戲的請求（僅允許房主觸發）
func (p *Player) handleGameStart(room *Room) error {
	if !p.IsHost {
		return errors.New(ErrorMessages[ErrCodeNotHost])
	}
	return room.StartGame()
}

// handleAnswer 處理玩家提交的答案
func (p *Player) handleAnswer(payload interface{}, room *Room) error {
	if !room.IsGameStarted() {
		return errors.New(ErrorMessages[ErrCodeGameNotStarted])
	}
	answer, ok := payload.(string)
	if !ok {
		return errors.New(ErrorMessages[ErrCodeInvalidAnswer])
	}
	return room.HandleAnswer(p.ID, answer)
}

// handleGameReset 處理重置遊戲的請求（僅允許房主觸發）
func (p *Player) handleGameReset(room *Room) error {
	if !p.IsHost {
		return errors.New(ErrorMessages[ErrCodeNotHost])
	}
	return room.GameReset()
}

// ResetGame 重置玩家遊戲狀態（例如重新開始時使用）
func (p *Player) ResetGame() {
	p.Game = NewGame()
	p.IsReady = false
	p.Score = 0
}

// UpdateScore 根據答案正確與否更新分數與進度
func (p *Player) UpdateScore(correct bool) {
	if correct {
		p.Score += 10
		p.Game.Progress++
	} else {
		p.Game.WrongCount ++
		p.Score -= 5
	}
	if p.Game.Progress >= p.Game.TotalQuiz {
		p.Game.IsFinished = true
	}
}

// Close 安全地關閉玩家連線
func (p *Player) Close() {
	if p.Conn != nil {
		p.Conn.Close()
	}
}
