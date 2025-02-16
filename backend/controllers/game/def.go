package game

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RoomStatus 定義房間狀態類型
type RoomStatus string

// 房間狀態常數
const (
	RoomStatusWaiting  RoomStatus = "waiting"  // 等待玩家加入
	RoomStatusPlaying  RoomStatus = "playing"  // 遊戲進行中
	RoomStatusFinished RoomStatus = "finished" // 遊戲結束
)

// WebSocket 消息類型常數
const (
	MsgTypeAnswer      = "answer"
	MsgTypeGameState   = "game_state"
	MsgTypeGameEnd     = "game_end"
	MsgTypeError       = "error"
	MsgTypeJoinRoom    = "join_room"
	MsgTypeLeaveRoom   = "leave_room"
	MsgTypePlayerList  = "player_list"
	MsgTypeGameStart   = "game_start" // 遊戲開始（由玩家或房主發送）
	MsgTypeGameRank    = "game_rank"
	MsgTypeProgress    = "progress"     // 進度更新
	MsgTypeReady       = "ready"        // 準備狀態
	MsgTypeGameReset   = "game_reset"   // 新增重新開始的消息類型
)

// 遊戲相關常數
const (
	MaxPlayers    = 10
	MinPlayers    = 1
	TimeLimit     = 60
	AnswerTimeout = 5
)

// 錯誤碼與錯誤訊息
const (
	ErrCodeRoomFull         = "room_full"
	ErrCodeNotHost          = "not_host"
	ErrCodeGameNotStarted   = "game_not_started"
	ErrCodeGameInProgress   = "game_in_progress"
	ErrCodeNotReady         = "not_ready"
	ErrCodeNotEnoughPlayers = "not_enough_players"
	ErrCodeInvalidAnswer    = "invalid_answer"
	ErrCodePlayerNotFound   = "player_not_found"
	ErrCodeInvalidMessage   = "invalid_message"
)

var ErrorMessages = map[string]string{
	ErrCodeRoomFull:         "房間已滿",
	ErrCodeNotHost:          "不是房主",
	ErrCodeGameNotStarted:   "遊戲尚未開始",
	ErrCodeGameInProgress:   "遊戲進行中",
	ErrCodeNotReady:         "玩家未準備",
	ErrCodeNotEnoughPlayers: "玩家人數不足",
	ErrCodeInvalidAnswer:    "無效的答案",
	ErrCodePlayerNotFound:   "找不到玩家",
	ErrCodeInvalidMessage:   "無效的消息",
}

// Game 相關結構（遊戲內所有狀態與數據）
type Game struct {
	QuizList     []string  `json:"quiz_list"`
	ColorList    []string  `json:"color_list"`
	DisplayColor string    `json:"display_color"`
	Progress     int       `json:"progress"`
	TotalQuiz    int       `json:"total_quiz"`
	WrongCount   int       `json:"wrong_count"`
	IsFinished   bool      `json:"is_finished"`
	StartTime    time.Time `json:"start_time"`
	PlayerID     string    `json:"player_id"`
}

// GameStatus 為前端提供的遊戲狀態資訊
type GameStatus struct {
	Quiz         string  `json:"quiz"`
	DisplayColor string  `json:"displayColor"`
	Progress     int     `json:"progress"`
	Percentage   float64 `json:"percentage"`
	WrongCount   int     `json:"wrongCount"`
	IsFinished   bool    `json:"isFinished"`
	TotalQuiz    int     `json:"totalQuiz"`
}

// Player 相關結構
type Player struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	IsHost  bool            `json:"is_host"`
	IsReady bool            `json:"is_ready"`
	Score   int             `json:"score"`
	Conn    *websocket.Conn `json:"-"`
	Game    *Game           `json:"game"`
}

// Room 相關結構，包含房間 ID、房主、玩家列表與狀態等
type Room struct {
	ID        string
	Players   map[string]*Player
	Status    RoomStatus
	StartTime time.Time
	mu        sync.Mutex
}

// WSMessage 定義 WebSocket 消息結構
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// WSError 定義 WebSocket 錯誤結構
type WSError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// PlayerRank 排名資訊結構
type PlayerRank struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Score      int           `json:"score"`
	WrongCount int           `json:"wrong_count"`
	Duration   time.Duration `json:"duration"`
	IsFinished bool          `json:"is_finished"`
}

// Message 定義 WebSocket 的消息格式
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// GameState 定義遊戲狀態消息結構
type GameState struct {
	Quiz         string `json:"quiz"`
	DisplayColor string `json:"display_color"`
	Progress     int    `json:"progress"`
	TotalQuiz    int    `json:"total_quiz"`
	TimeLeft     int    `json:"time_left"`
	IsFinished   bool   `json:"is_finished"`
}
