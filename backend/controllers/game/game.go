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

const (
	QuizCount = 10
	TimeUnit  = int64(time.Millisecond)
)

// 可用的顏色列表
var ValidColors = []string{"red", "green", "blue", "yellow", "orange", "purple"}

func NewGame() *Game {
    game := &Game{
        QuizList:   make([]string, QuizCount),
        ColorList:  make([]string, QuizCount),
        rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
    }
    game.generateColors()
    return game
}

func (g *Game) GetStatus() (GameStatus, error) {
	if g.WhichQuiz >= len(g.QuizList) {
		return GameStatus{}, ErrGameNotStarted
	}

	timeUsed := g.calculateTimeUsed()
	progress := float64(g.WhichQuiz) / float64(QuizCount) * 100

	return GameStatus{
		Quiz:         g.QuizList[g.WhichQuiz],
		DisplayColor: g.ColorList[g.WhichQuiz],
		Progress:     g.WhichQuiz,
		Percentage:   progress,
		WrongCount:   g.WrongCount,
		IsFinished:   g.IsFinished,
		TimeUsed:     timeUsed,
		TotalQuiz:    QuizCount,
	}, nil
}

func (g *Game) Answer(color string) (bool, bool, error) {
	if g.WhichQuiz >= len(g.QuizList) {
		return false, false, ErrGameNotStarted
	}

	if g.IsFinished {
		return false, true, ErrGameFinished
	}

	if !isValidColor(color) {
		return false, false, ErrInvalidColor
	}

	correct := color == g.QuizList[g.WhichQuiz]
	if correct {
		g.WhichQuiz++
		g.IsFinished = g.WhichQuiz >= len(g.QuizList)
		if g.IsFinished {
			g.EndTime = time.Now().UnixNano() / TimeUnit
		}
	} else {
		g.WrongCount++
	}

	return correct, g.IsFinished, nil
}

func (g *Game) Restart() {
	g.generateColors()
	g.WhichQuiz = 0
	g.WrongCount = 0
	g.StartTime = time.Now().UnixNano() / TimeUnit
	g.EndTime = 0
	g.IsFinished = false
}

// 將顏色生成邏輯抽取為單獨的方法
func (g *Game) generateColors() {
    for i := 0; i < QuizCount; i++ {
        g.QuizList[i] = ValidColors[g.rng.Intn(len(ValidColors))]
        g.ColorList[i] = ValidColors[g.rng.Intn(len(ValidColors))]
    }
}

func (g *Game) calculateTimeUsed() float64 {
	if g.EndTime > 0 {
		return float64(g.EndTime-g.StartTime) / 1000 // 轉換為秒
	}
	currentTime := time.Now().UnixNano() / TimeUnit
	return float64(currentTime-g.StartTime) / 1000
}

// 使用 map 優化顏色驗證
var validColorMap = make(map[string]bool)

func init() {
	for _, color := range ValidColors {
		validColorMap[color] = true
	}
}

func isValidColor(color string) bool {
	return validColorMap[color]
}
