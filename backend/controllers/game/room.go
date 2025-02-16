package game

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/rejxcy/logger"
)

// NewRoom 創建新房間並初始化內部資料結構
func NewRoom(id string) *Room {
	return &Room{
		ID:        id,
		Players:   make(map[string]*Player),
		Status:    RoomStatusWaiting,
		StartTime: time.Now(),
		mu:        sync.Mutex{},
	}
}

// AddPlayer 將新的玩家加入房間中
func (r *Room) AddPlayer(player *Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Players) >= MaxPlayers {
		return errors.New(ErrorMessages[ErrCodeRoomFull])
	}
	r.Players[player.ID] = player
	return nil
}

// RemovePlayer 從房間中移除玩家
func (r *Room) RemovePlayer(playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Players, playerID)
}

// BroadcastPlayerList 廣播更新後的玩家列表（包含進度、錯誤數、分數與排名）給所有玩家
func (r *Room) BroadcastPlayerList() {
	// 複製所有非房主玩家資訊到 slice 中，方便進行排序
	r.mu.Lock()
	playersSlice := make([]*Player, 0, len(r.Players))
	for _, p := range r.Players {
		if p.IsHost {
			continue
		}
		playersSlice = append(playersSlice, p)
	}
	r.mu.Unlock()

	// 根據分數進行排序
	sort.Slice(playersSlice, func(i, j int) bool {
		if playersSlice[i].Score != playersSlice[j].Score {
			return playersSlice[i].Score > playersSlice[j].Score // 降序排列
		}
		return playersSlice[i].Game.WrongCount < playersSlice[j].Game.WrongCount // 錯誤次數少者排前
	})

	// 為每位玩家分配排名，並準備要廣播給前端的資料
	rankingList := make([]map[string]interface{}, 0, len(playersSlice))
	for idx, p := range playersSlice {
		rankingList = append(rankingList, map[string]interface{}{
			"id":         p.ID,
			"name":       p.Name,
			"isReady":    p.IsReady,
			"progress":   p.Game.Progress,
			"wrongCount": p.Game.WrongCount,
			"score":      p.Score,
			"rank":       idx + 1,
		})
	}

	// 將整個玩家列表（含排名資訊）發送給所有連線的玩家
	msg := Message{
		Type:    MsgTypePlayerList,
		Payload: rankingList,
	}
	for _, p := range r.Players {
		if err := p.Send(msg); err != nil {
			logger.Output.Error("廣播玩家列表給 %s 失敗: %v", p.Name, err)
		}
	}
}

// IsGameStarted 判斷遊戲是否已開始
func (r *Room) IsGameStarted() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Status == RoomStatusPlaying
}

// StartGame 為所有玩家初始化獨立遊戲進度，並廣播初始狀態、遊戲開始訊息
func (r *Room) StartGame() error {
	r.mu.Lock()

	// 檢查所有非房主玩家是否皆準備好
	readyCount := 0
	for _, p := range r.Players {
		if !p.IsHost && p.IsReady {
			readyCount++
		}
	}
	if len(r.Players) < 2 || readyCount < len(r.Players)-1 {
		r.mu.Unlock()
		return errors.New("玩家不足或部分玩家尚未準備")
	}

	// 對每位玩家建立獨立的遊戲進度
	for _, p := range r.Players {
		p.Game = NewGame()
	}

	r.Status = RoomStatusPlaying
	r.mu.Unlock() // 釋放鎖後再發送訊息

	// 廣播每位玩家的初始遊戲狀態
	for _, p := range r.Players {
		state, err := p.Game.GetStatus()
		if err != nil {
			logger.Output.Error("取得 %s 遊戲狀態失敗: %v", p.Name, err)
			continue
		}
		gameStatePayload := map[string]interface{}{
			"quiz":         state.Quiz,
			"displayColor": state.DisplayColor,
			"progress":     state.Progress,
			"wrongCount":   state.WrongCount,
			"totalQuiz":    state.TotalQuiz,
			"isFinished":   state.IsFinished,
		}
		gameStateMsg := Message{
			Type:    "game_state",
			Payload: gameStatePayload,
		}
		if err := p.Send(gameStateMsg); err != nil {
			logger.Output.Error("廣播遊戲狀態給 %s 失敗: %v", p.Name, err)
		}
	}

	// 廣播遊戲開始訊息給所有玩家
	gameStartMsg := Message{
		Type:    MsgTypeGameStart,
		Payload: "遊戲開始",
	}
	for _, p := range r.Players {
		if err := p.Send(gameStartMsg); err != nil {
			logger.Output.Error("通知 %s 遊戲開始失敗: %v", p.Name, err)
		}
	}

	logger.Output.Info("房間 %s 遊戲開始, 總玩家數: %d", r.ID, len(r.Players)-1)
	return nil
}

// HandleAnswer 處理玩家提交的答案，並更新該玩家獨立的遊戲進度
func (r *Room) HandleAnswer(playerID, answer string) error {
	r.mu.Lock()
	player, exists := r.Players[playerID]
	if !exists {
		r.mu.Unlock()
		return errors.New(ErrorMessages[ErrCodePlayerNotFound])
	}
	if player.Game == nil {
		r.mu.Unlock()
		return errors.New(ErrorMessages[ErrCodeGameNotStarted])
	}
	logger.Output.Info("玩家 %s 提交答案: %s", player.Name, answer)

	// 利用玩家自身的 Game 處理答案
	correct, err := player.Game.Answer(answer)
	// 更新玩家分數
	player.UpdateScore(correct)
	
	r.mu.Unlock()

	if err != nil {
		logger.Output.Error("處理玩家 %s 提交答案失敗: %v", player.Name, err)
		return err
	}


	// 新增：發送更新後的遊戲狀態（包含最新題目等資訊）
	state, err := player.Game.GetStatus()
	if err != nil {
		logger.Output.Error("取得 %s 遊戲狀態失敗: %v", player.Name, err)
	} else {
		updatedGameStateMsg := Message{
			Type: "game_state",
			Payload: map[string]interface{}{
				"quiz":         state.Quiz,
				"displayColor": state.DisplayColor,
				"progress":     state.Progress,
				"wrongCount":   state.WrongCount,
				"totalQuiz":    state.TotalQuiz,
				"isFinished":   state.IsFinished,
			},
		}
		if err := player.Send(updatedGameStateMsg); err != nil {
			logger.Output.Error("更新遊戲狀態給 %s 失敗: %v", player.Name, err)
		}
	}

	// 檢查遊戲是否結束
	if r.gameFinish() {
		logger.Output.Info("遊戲結束，廣播結束訊息")
		r.Broadcast(Message{
			Type: MsgTypeGameEnd,
			Payload: "遊戲結束",
		})
	}
	
	// 廣播更新後的玩家列表（狀態）
	r.BroadcastPlayerList()
	return nil
}

// 廣播消息給所有玩家
func (r *Room) Broadcast(msg Message) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, player := range r.Players {
		player.Send(msg)
	}
}

// 處理房間消息
func (r *Room) HandleMessage(playerID string, msg Message) error {
	r.mu.Lock()
	player, exists := r.Players[playerID]
	r.mu.Unlock()

	if !exists {
		return errors.New(ErrorMessages[ErrCodePlayerNotFound])
	}

	switch msg.Type {
	case MsgTypeReady:
		if player.IsHost {
			logger.Output.Error("Host cannot set ready state")
			return errors.New(ErrorMessages[ErrCodeNotHost])
		}

		ready, ok := msg.Payload.(bool)
		if !ok {
			logger.Output.Error("Invalid ready state payload")
			return errors.New(ErrorMessages[ErrCodeInvalidMessage])
		}

		player.IsReady = ready
		logger.Output.Info("Player %s ready state changed to: %v", player.Name, ready)
		r.BroadcastPlayerList()
		return nil

	case MsgTypeGameStart:
		if !player.IsHost {
			logger.Output.Error("Non-host player tried to start game")
			return errors.New(ErrorMessages[ErrCodeNotHost])
		}
		return r.StartGame()

	case MsgTypeGameReset:
		if !player.IsHost {
			logger.Output.Error("Non-host player tried to restart game")
			return errors.New(ErrorMessages[ErrCodeNotHost])
		}
		return r.GameReset()

	default:
		return fmt.Errorf("未知的消息類型: %s", msg.Type)
	}
}

// RestartGame 重置每位玩家的遊戲狀態，並發送重新開始的通知
func (r *Room) GameReset() error {
	r.mu.Lock()
	r.Status = RoomStatusWaiting
	for _, p := range r.Players {
		p.ResetGame() // 每位玩家自行重置遊戲狀態
	}
	r.mu.Unlock()

	gameResetMsg := Message{
		Type:    MsgTypeGameReset,
		Payload: "遊戲重新開始",
	}
	r.Broadcast(gameResetMsg)

	r.BroadcastPlayerList()

	return nil
}

// 檢查遊戲是否結束
func (r *Room) gameFinish() bool {
	for _, p := range r.Players {
		if !p.IsHost && !p.Game.IsFinished {
			return false
		}
	}
	return true
}
