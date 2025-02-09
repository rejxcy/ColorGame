package game

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 枚舉類型
type RoomStatus string

// 常量定義
const (
	// Room 狀態
	RoomStatusWaiting  RoomStatus = "waiting"  // 等待玩家加入
	RoomStatusPlaying  RoomStatus = "playing"  // 遊戲進行中
	RoomStatusFinished RoomStatus = "finished" // 遊戲結束

	// WebSocket 消息類型
	MsgTypeAnswer     = "answer"
	MsgTypeRestart    = "restart"
	MsgTypeGameState  = "game_state"
	MsgTypeGameOver   = "game_over"
	MsgTypeError      = "error"
	MsgTypeJoinRoom   = "join_room"
	MsgTypeLeaveRoom  = "leave_room"
	MsgTypePlayerList = "player_list"
	MsgTypeGameStart  = "game_start"
	MsgTypeGameRank   = "game_rank"
	MsgTypeProgress   = "progress"   // 新增：進度更新
	MsgTypeReady      = "ready"      // 新增：準備狀態
	MsgTypeStartGame  = "start_game" // 開始遊戲

	// 遊戲相關常數
	MaxPlayers    = 10
	MinPlayers    = 1
	TimeLimit     = 60
	AnswerTimeout = 5
)

// 錯誤碼和訊息
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

// 錯誤訊息
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

// 消息類型常數
const (
	msgTypeAnswer     = "answer"      // 答案
	msgTypeRestart    = "restart"     // 重新開始
	msgTypeGameState  = "game_state"  // 遊戲狀態
	msgTypeGameOver   = "game_over"   // 遊戲結束
	msgTypeError      = "error"       // 錯誤
	msgTypeJoinRoom   = "join_room"   // 加入房間
	msgTypeLeaveRoom  = "leave_room"  // 離開房間
	msgTypePlayerList = "player_list" // 玩家列表
	msgTypeGameStart  = "game_start"  // 遊戲開始
	msgTypeGameRank   = "game_rank"   // 遊戲排名
	msgTypeProgress   = "progress"    // 進度更新
	msgTypeReady      = "ready"       // 準備狀態
)

// 錯誤訊息
const (
	errRoomFull        = "房間已滿"
	errRoomNotFound    = "房間不存在"
	errPlayerNotFound  = "玩家不存在"
	errGameNotStarted  = "遊戲尚未開始"
	errGameInProgress  = "遊戲進行中"
	errInvalidAnswer   = "無效的答案"
	errNotHost         = "只有房主可以執行此操作"
	errNotEnoughPlayer = "玩家人數不足"
)

// Game 相關結構
type Game struct {
	QuizList     []string  `json:"quiz_list"`
	ColorList    []string  `json:"color_list"`
	CurrentQuiz  string    `json:"current_quiz"`
	DisplayColor string    `json:"display_color"`
	Progress     int       `json:"progress"`
	TotalQuiz    int       `json:"total_quiz"`
	WrongCount   int       `json:"wrong_count"`
	IsFinished   bool      `json:"is_finished"`
	StartTime    time.Time `json:"start_time"`
	WhichQuiz    int       `json:"which_quiz"`
	PlayerID     string    `json:"player_id"`
}

// 遊戲狀態
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

// Room 相關結構
type Room struct {
	ID        string             `json:"id"`
	Host      string             `json:"host"`      // 房主ID
	Players   map[string]*Player `json:"players"`   // 玩家列表
	Status    RoomStatus         `json:"status"`    // 房間狀態
	StartTime time.Time          `json:"startTime"` // 遊戲開始時間
	mu        sync.Mutex
}

// WebSocket 相關結構
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type WSError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// PlayerRank 排名信息
type PlayerRank struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Score      int           `json:"score"`
	WrongCount int           `json:"wrong_count"`
	Duration   time.Duration `json:"duration"`
	IsFinished bool          `json:"is_finished"`
}

// Message 定義 WebSocket 消息結構
type Message struct {
	Type    string      `json:"type"`    // 消息類型
	Payload interface{} `json:"payload"` // 消息內容
}

// GameState 定義遊戲狀態結構
type GameState struct {
	Quiz         string `json:"quiz"`
	DisplayColor string `json:"display_color"`
	Progress     int    `json:"progress"`
	TotalQuiz    int    `json:"total_quiz"`
	TimeLeft     int    `json:"time_left"`
	IsFinished   bool   `json:"is_finished"`
}
