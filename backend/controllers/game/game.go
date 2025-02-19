package game

import (
	"errors"
	"math/rand"
	"time"
)

var (
	ErrGameFinished   = errors.New("遊戲已結束")
	ErrInvalidColor   = errors.New("無效的顏色選擇")
	ErrGameNotStarted = errors.New("遊戲尚未開始")
)

// 可用的顏色列表
var ValidColors = []string{"red", "green", "blue", "yellow", "orange", "purple"}

// 使用 map 優化顏色驗證
var validColorMap = make(map[string]bool)

func init() {
	for _, color := range ValidColors {
		validColorMap[color] = true
	}
}

// NewGame 創建並初始化一個新遊戲
func NewGame() *Game {
	game := &Game{
		QuizList:   make([]string, QuizCount),
		ColorList:  make([]string, QuizCount),
		TotalQuiz:  QuizCount,
		WrongCount: 0,
		IsFinished: false,
	}
	game.generateColors()
	return game
}

// 返回目前遊戲狀態
func (g *Game) GetStatus() (GameStatus, error) {
	if g.IsFinished {
		return GameStatus{
			Progress:   QuizCount,
			WrongCount: g.WrongCount,
			IsFinished: true,
			TotalQuiz:  QuizCount,
		}, nil
	}

	return GameStatus{
		Quiz:         g.QuizList[g.Progress],
		DisplayColor: g.ColorList[g.Progress],
		Progress:     g.Progress,
		WrongCount:   g.WrongCount,
		IsFinished:   g.IsFinished,
		TotalQuiz:    QuizCount,
	}, nil
}

// 判斷答案是否正確，並更新遊戲狀態
func (g *Game) Answer(color string) (bool, error) {
	if g.IsFinished {
		return false, ErrGameFinished
	}

	if !isValidColor(color) {
		return false, ErrInvalidColor
	}

	correct := color == g.QuizList[g.Progress]
	if correct {
		if g.Progress >= QuizCount {
			g.IsFinished = true
		}
	}
	return correct, nil
}

// 重置遊戲狀態
func (g *Game) Restart() {
	g.generateColors()
	g.Progress = 0
	g.WrongCount = 0
	g.IsFinished = false
}

// 產生新的題目與顏色列表
func (g *Game) generateColors() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < QuizCount; i++ {
		g.QuizList[i] = ValidColors[rng.Intn(len(ValidColors))]
		g.ColorList[i] = ValidColors[rng.Intn(len(ValidColors))]
	}
}

func isValidColor(color string) bool {
	return validColorMap[color]
}
