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
		QuizList:  make([]string, QuizCount),
		ColorList: make([]string, QuizCount),
	}
	game.generateColors()
	return game
}

func (g *Game) GetStatus() (GameStatus, error) {
	if g.IsFinished {
		return GameStatus{
			Progress:   QuizCount,
			Percentage: 100,
			WrongCount: g.WrongCount,
			IsFinished: true,
			TotalQuiz:  QuizCount,
		}, nil
	}

	return GameStatus{
		Quiz:         g.QuizList[g.WhichQuiz],
		DisplayColor: g.ColorList[g.WhichQuiz],
		Progress:     g.WhichQuiz,
		Percentage:   float64(g.WhichQuiz) / float64(QuizCount) * 100,
		WrongCount:   g.WrongCount,
		IsFinished:   g.IsFinished,
		TotalQuiz:    QuizCount,
	}, nil
}

func (g *Game) Answer(color string) (bool, bool, error) {
	if g.IsFinished {
		return false, true, ErrGameFinished
	}

	if !isValidColor(color) {
		return false, false, ErrInvalidColor
	}

	correct := color == g.QuizList[g.WhichQuiz]
	if correct {
		g.WhichQuiz++
		if g.WhichQuiz >= len(g.QuizList) {
			g.IsFinished = true
			return true, true, nil
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
	g.IsFinished = false
}

func (g *Game) generateColors() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < QuizCount; i++ {
		g.QuizList[i] = ValidColors[rng.Intn(len(ValidColors))]
		g.ColorList[i] = ValidColors[rng.Intn(len(ValidColors))]
	}
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
