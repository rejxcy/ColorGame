package game

import (
	"sync"

	"github.com/gorilla/websocket"
)

// RoomStatus 定義房間狀態類型
type RoomStatus string

// 房間狀態常數
const (
	RoomStatusWaiting  RoomStatus = "waiting"
	RoomStatusPlaying  RoomStatus = "playing"
	RoomStatusFinished RoomStatus = "finished"
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
	MsgTypeGameStart   = "game_start"
	MsgTypeProgress    = "progress"
	MsgTypeReady       = "ready"
	MsgTypeGameReset   = "game_reset"
)

// 遊戲相關常數
const (
	MaxPlayers    = 10
	MinPlayers    = 1
	QuizCount = 10
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

type RoomManager struct {
	rooms map[string]*Room
	mu    sync.Mutex
}

type Room struct {
	ID        string
	Players   map[string]*Player
	Status    RoomStatus
	mu        sync.Mutex
}

// WebSocket 的消息格式
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// 遊戲內所有狀態與數據
type Game struct {
	QuizList     []string  `json:"quiz_list"`
	ColorList    []string  `json:"color_list"`
	DisplayColor string    `json:"display_color"`
	Progress     int       `json:"progress"`
	TotalQuiz    int       `json:"total_quiz"`
	WrongCount   int       `json:"wrong_count"`
	IsFinished   bool      `json:"is_finished"`
	PlayerID     string    `json:"player_id"`
}

// 為前端提供的遊戲狀態資訊
type GameStatus struct {
	Quiz         string  `json:"quiz"`
	DisplayColor string  `json:"displayColor"`
	Progress     int     `json:"progress"`
	WrongCount   int     `json:"wrongCount"`
	TotalQuiz    int     `json:"totalQuiz"`
	IsFinished   bool    `json:"isFinished"`
}

type Player struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	IsHost  bool            `json:"is_host"`
	IsReady bool            `json:"is_ready"`
	Score   int             `json:"score"`
	Conn    *websocket.Conn `json:"-"`
	Game    *Game           `json:"game"`
}
