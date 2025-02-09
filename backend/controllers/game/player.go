package game

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rejxcy/logger"
)

// NewPlayer 創建新玩家
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

// Send 發送消息給玩家
func (p *Player) Send(msg Message) error {
	return p.Conn.WriteJSON(msg)
}

// SendError 發送錯誤消息給玩家
func (p *Player) SendError(code string, message string) {
	p.Send(Message{
		Type: MsgTypeError,
		Payload: map[string]interface{}{
			"code":    code,
			"message": message,
		},
	})
}

// HandleMessage 處理玩家消息
func (p *Player) HandleMessage(msg Message, room *Room) error {
	switch msg.Type {
	case "ready":
		// 房主不能設置準備狀態
		if p.IsHost {
			return errors.New("房主不需要準備")
		}
		ready, ok := msg.Payload.(bool)
		if !ok {
			return errors.New("無效的準備狀態")
		}
		p.IsReady = ready
		logger.Output.Info("Player %s ready state changed to: %v", p.Name, ready)

		// 廣播玩家列表更新
		room.BroadcastPlayerList()
		return nil
	case MsgTypeGameStart:
		return p.handleGameStart(room)
	case MsgTypeAnswer:
		return p.handleAnswer(msg.Payload, room)
	default:
		return errors.New("未知的消息類型")
	}
}

// handleGameStart 處理開始遊戲
func (p *Player) handleGameStart(room *Room) error {
	if !p.IsHost {
		return errors.New(ErrorMessages[ErrCodeNotHost])
	}
	return room.StartGame()
}

// handleAnswer 處理答案
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

// ResetGame 重置玩家遊戲狀態
func (p *Player) ResetGame() {
	p.Game = NewGame()
	p.IsReady = false
	p.Score = 0
}

// UpdateScore 更新玩家分數
func (p *Player) UpdateScore(correct bool) {
	if correct {
		p.Score += 10
	} else {
		p.Game.WrongCount++
	}
	p.Game.Progress++

	if p.Game.Progress >= p.Game.TotalQuiz {
		p.Game.IsFinished = true
	}
}

// GetRank 獲取玩家排名資訊
func (p *Player) GetRank() PlayerRank {
	duration := time.Since(p.Game.StartTime)
	return PlayerRank{
		ID:         p.ID,
		Name:       p.Name,
		Score:      p.Score,
		WrongCount: p.Game.WrongCount,
		Duration:   duration,
		IsFinished: p.Game.IsFinished,
	}
}

// Close 關閉玩家連接
func (p *Player) Close() {
	if p.Conn != nil {
		p.Conn.Close()
	}
}
