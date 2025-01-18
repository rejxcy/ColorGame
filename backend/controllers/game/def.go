package game

import (
	"math/rand"
	"sync"

	"github.com/gorilla/websocket"
)

type Game struct {
    QuizList   []string    
    ColorList  []string    
    WhichQuiz  int         
    WrongCount int         
    StartTime  int64       
    EndTime    int64       
    IsFinished bool
    PlayerID   string
    rng        *rand.Rand     
}

// 遊戲狀態
type GameStatus struct {
    Quiz         string  `json:"quiz"`
    DisplayColor string  `json:"displayColor"`
    Progress     int     `json:"progress"`
    Percentage   float64 `json:"percentage"`
    WrongCount   int     `json:"wrongCount"`
    IsFinished   bool    `json:"isFinished"`
    TimeUsed     float64 `json:"timeUsed,omitempty"`
    TotalQuiz    int     `json:"totalQuiz"`
}

type GameConnection struct {
    conn     *websocket.Conn
    game     *Game
    mu       sync.Mutex
    done     chan struct{}
    playerID string
}

// WebSocket 消息類型
const (
    MsgTypeAnswer     = "answer"
    MsgTypeRestart    = "restart"
    MsgTypeGameState  = "game_state"
    MsgTypeGameOver   = "game_over"
    MsgTypeError      = "error"
)

type WSMessage struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

type WSError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}