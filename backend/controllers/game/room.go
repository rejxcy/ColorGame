package game

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/rejxcy/logger"
)

// NewRoom 創建新房間
func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Players: make(map[string]*Player),
		Status:  RoomStatusWaiting,
		mu:      sync.Mutex{},
	}
}

// AddPlayer 將玩家加入房間
func (r *Room) AddPlayer(player *Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	logger.Output.Info("Start: Attempting to add player %s to room %s", player.Name, r.ID)

	// 檢查玩家是否已經在房間中
	for _, existingPlayer := range r.Players {
		if existingPlayer.Name == player.Name {
			logger.Output.Info("Player %s already exists, updating connection", player.Name)
			// 關閉舊的連接
			existingPlayer.Conn.Close()
			// 更新玩家的連接
			existingPlayer.Conn = player.Conn
			// 保持原有的狀態（如準備狀態等）
			return nil
		}
	}

	// 檢查玩家數量限制
	if len(r.Players) >= MaxPlayers {
		logger.Output.Error("Failed: Room %s is full", r.ID)
		return errors.New(ErrorMessages[ErrCodeRoomFull])
	}

	// 加入新玩家
	r.Players[player.ID] = player
	logger.Output.Info("Player map updated: %v", r.Players)

	// 如果是第一個玩家，設置為房主
	if len(r.Players) == 1 {
		r.Host = player.ID
		player.IsHost = true
		logger.Output.Info("Set player %s as host of room %s", player.Name, r.ID)
	}

	logger.Output.Info("[Success: Player %s (ID: %s) added to room %s. Total players: %d",
		player.Name, player.ID, r.ID, len(r.Players))

	return nil
}

// RemovePlayer 從房間移除玩家
func (r *Room) RemovePlayer(playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if player, exists := r.Players[playerID]; exists {
		logger.Output.Info("Removing player %s from room %s", player.Name, r.ID)
		delete(r.Players, playerID)

		// 如果移除的是房主且還有其他玩家，選擇新房主
		if playerID == r.Host && len(r.Players) > 0 {
			// 選擇第一個玩家作為新房主
			for newHostID, newHost := range r.Players {
				r.Host = newHostID
				newHost.IsHost = true
				logger.Output.Info("New host selected: %s", newHost.Name)
				break
			}
		}

		logger.Output.Info("Player removed. Remaining players: %d", len(r.Players))
	}
}

// 廣播消息給所有玩家
func (r *Room) Broadcast(msg Message) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, player := range r.Players {
		player.Send(msg)
	}
}

// 向房間內所有玩家廣播最新的玩家列表
func (r *Room) BroadcastPlayerList() {
	playerList := r.GetPlayerList()
	logger.Output.Info("Broadcasting updated player list")

	r.Broadcast(Message{
		Type:    "player_list",
		Payload: playerList,
	})
}

// 返回當前房間的玩家列表
func (r *Room) GetPlayerList() []map[string]interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()

	playerList := make([]map[string]interface{}, 0)
	for _, p := range r.Players {
		if !p.IsHost { // 只返回非房主玩家
			playerInfo := map[string]interface{}{
				"id":      p.ID,
				"name":    p.Name,
				"isHost":  p.IsHost,
				"isReady": p.IsReady,
			}
			playerList = append(playerList, playerInfo)
		}
	}

	return playerList
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

	case MsgTypeStartGame:
		if !player.IsHost {
			logger.Output.Error("Non-host player tried to start game")
			return errors.New(ErrorMessages[ErrCodeNotHost])
		}
		return r.StartGame()

	default:
		return fmt.Errorf("未知的消息類型: %s", msg.Type)
	}
}

// 開始遊戲
func (r *Room) StartGame() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	logger.Output.Info("[StartGame] Attempting to start game in room %s", r.ID)

	// 檢查遊戲狀態
	if r.Status == RoomStatusPlaying {
		logger.Output.Error("[StartGame] Game already in progress")
		return errors.New(ErrorMessages[ErrCodeGameInProgress])
	}

	// 檢查玩家數量（不包括房主）
	readyPlayers := 0
	totalPlayers := 0
	for _, p := range r.Players {
		if !p.IsHost {
			totalPlayers++
			if p.IsReady {
				readyPlayers++
			}
		}
	}

	logger.Output.Info("[StartGame] Total players: %d, Ready players: %d", totalPlayers, readyPlayers)

	// 檢查玩家數量
	if totalPlayers < MinPlayers {
		logger.Output.Error("[StartGame] Not enough players")
		return errors.New(ErrorMessages[ErrCodeNotEnoughPlayers])
	}

	// 檢查是否所有玩家都準備好
	if readyPlayers < totalPlayers {
		logger.Output.Error("[StartGame] Not all players are ready")
		return errors.New(ErrorMessages[ErrCodeNotReady])
	}

	// 開始遊戲
	r.Status = RoomStatusPlaying
	r.StartTime = time.Now()

	// 廣播遊戲開始消息
	startGameMsg := Message{
		Type: MsgTypeStartGame,
		Payload: map[string]interface{}{
			"status":    string(r.Status),
			"startTime": r.StartTime,
		},
	}

	logger.Output.Info("[StartGame] Broadcasting game start message")
	r.Broadcast(startGameMsg)

	logger.Output.Info("[StartGame] Game started successfully in room %s", r.ID)
	return nil
}

// RestartGame 重新開始遊戲
func (r *Room) RestartGame() {
	r.mu.Lock()
	r.Status = RoomStatusWaiting
	for _, p := range r.Players {
		p.ResetGame()
	}
	r.mu.Unlock()
	r.BroadcastPlayerList()
}

// 檢查遊戲是否已開始
func (r *Room) IsGameStarted() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Status == RoomStatusPlaying
}

// 處理玩家答案
func (r *Room) HandleAnswer(playerID string, answer string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	player, exists := r.Players[playerID]
	if !exists {
		return errors.New(ErrorMessages[ErrCodePlayerNotFound])
	}

	if r.Status != RoomStatusPlaying {
		return errors.New(ErrorMessages[ErrCodeGameNotStarted])
	}

	correct := CheckAnswer(player.Game.CurrentQuiz, answer)
	player.UpdateScore(correct)

	if r.CheckGameFinish() {
		r.Status = RoomStatusFinished
		r.BroadcastFinalRanking()
	} else if player.Game.IsFinished {
		r.BroadcastRanking()
	} else {
		r.StartNewRound()
	}

	return nil
}

// 開始新回合
func (r *Room) StartNewRound() {
	quiz, displayColor := GenerateQuiz()
	r.Broadcast(Message{
		Type: MsgTypeGameState,
		Payload: map[string]interface{}{
			"quiz":         quiz,
			"displayColor": displayColor,
		},
	})
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

// 生成新題目
func GenerateQuiz() (string, string) {
	// TODO: 實現題目生成邏輯
	return "紅色", "藍色"
}

// 檢查答案
func CheckAnswer(quiz string, answer string) bool {
	// TODO: 實現答案檢查邏輯
	return true
}
