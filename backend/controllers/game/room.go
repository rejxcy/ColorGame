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

// BroadcastPlayerList 廣播更新後的玩家列表（包含進度與錯誤數）給所有玩家
func (r *Room) BroadcastPlayerList() {
	// 複製一份玩家狀態，避免長時間持有鎖
	r.mu.Lock()
	playerList := make([]map[string]interface{}, 0, len(r.Players))
	for _, p := range r.Players {
		if p.IsHost { // 房主不列入
			continue
		}
		progress := 0
		wrongCount := 0
		if p.Game != nil {
			progress = p.Game.Progress
			wrongCount = p.Game.WrongCount
		}
		info := map[string]interface{}{
			"id":         p.ID,
			"name":       p.Name,
			"isHost":     p.IsHost,
			"isReady":    p.IsReady,
			"progress":   progress,
			"wrongCount": wrongCount,
		}
		playerList = append(playerList, info)
	}
	r.mu.Unlock()

	msg := Message{
		Type:    MsgTypePlayerList,
		Payload: playerList,
	}

	// 對每位玩家發送更新訊息（不在鎖區段中執行 Send）
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

	// 對每位玩家建立獨立的遊戲進度，但題目相同
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

	logger.Output.Info("房間 %s 遊戲開始, 總玩家數: %d", r.ID, len(r.Players))
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
	correct, finished, err := player.Game.Answer(answer)
	r.mu.Unlock() // 先釋放鎖
	logger.Output.Info("處理玩家 %s 提交答案結果: correct=%v, finished=%v, err=%v", player.Name, correct, finished, err)

	if err != nil {
		logger.Output.Error("處理玩家 %s 提交答案失敗: %v", player.Name, err)
		return err
	}

	// 傳送答案結果給該玩家
	answerResultMsg := Message{
		Type: "answer_result",
		Payload: map[string]interface{}{
			"correct":    correct,
			"progress":   player.Game.Progress,
			"wrongCount": player.Game.WrongCount,
			"isFinished": finished,
		},
	}
	if err := player.Send(answerResultMsg); err != nil {
		logger.Output.Error("傳送答案結果給 %s 失敗: %v", player.Name, err)
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

	case MsgTypeRestartGame:
		if !player.IsHost {
			logger.Output.Error("Non-host player tried to restart game")
			return errors.New(ErrorMessages[ErrCodeNotHost])
		}
		return r.RestartGame()

	default:
		return fmt.Errorf("未知的消息類型: %s", msg.Type)
	}
}

// RestartGame 重置每位玩家的遊戲狀態，並發送重新開始的通知
func (r *Room) RestartGame() error {
	r.mu.Lock()
	r.Status = RoomStatusWaiting
	for _, p := range r.Players {
		p.ResetGame() // 每位玩家自行重置遊戲狀態
	}
	r.mu.Unlock()

	restartMsg := Message{
		Type:    "game_restart",
		Payload: "遊戲重新開始",
	}
	for _, p := range r.Players {
		if err := p.Send(restartMsg); err != nil {
			logger.Output.Error("傳送重新開始訊息給 %s 失敗: %v", p.Name, err)
		}
	}
	r.BroadcastPlayerList()

	return nil
}

// 廣播排名
func (r *Room) BroadcastRanking() {
	rankings := r.GetRankings()
	r.Broadcast(Message{
		Type:    MsgTypeGameRank,
		Payload: rankings,
	})
}

// 廣播最終排名
func (r *Room) BroadcastFinalRanking() {
	r.BroadcastRanking()
}

// 獲取排名列表
func (r *Room) GetRankings() []PlayerRank {
	rankings := make([]PlayerRank, 0, len(r.Players))
	for _, p := range r.Players {
		rankings = append(rankings, p.GetRank())
	}

	sort.Slice(rankings, func(i, j int) bool {
		if rankings[i].Score != rankings[j].Score {
			return rankings[i].Score > rankings[j].Score
		}
		return rankings[i].Duration < rankings[j].Duration
	})
	return rankings
}

// 檢查遊戲是否結束
func (r *Room) CheckGameFinish() bool {
	for _, p := range r.Players {
		if !p.Game.IsFinished {
			return false
		}
	}
	return true
}
