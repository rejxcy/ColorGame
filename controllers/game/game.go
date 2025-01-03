package game

import (
	"math/rand"
	"strconv"
)

type Game struct {
	player	   *Player
	isGameEnd  bool
	quizList   []string
	colorList  []string
	whichQuiz  int
	wrongCount int
}

func NewGame(player *Player) *Game {

	return &Game{
		player: 	player,
		isGameEnd:  false,
		quizList:   randomColors(10),
		colorList:  randomColors(10),
		whichQuiz:  0,
		wrongCount: 0,
	}
}

func (g *Game) restart(gameTime string) {
	g.isGameEnd = false
	g.quizList = randomColors(10)
	g.colorList = randomColors(10)
	g.whichQuiz = 0
	time, err := strconv.ParseFloat(gameTime, 64)
	if err != nil {
		time = 0
	}
	g.player.UpdateRecord(time)
}

func (g *Game) getQuiz() (string, string) {
	if g.whichQuiz >= len(g.quizList) {
		g.isGameEnd = true
		return "", ""
	}

	quiz := g.quizList[g.whichQuiz]
	color := g.colorList[g.whichQuiz]

	return quiz, color
}

func (g *Game) isAnswer(color string) bool {
	if color == g.quizList[g.whichQuiz] {
		g.whichQuiz++
		return true
	}
	
	g.wrongCount ++
	return false
}

func randomColors(count int) []string {
	colors := []string{"red", "green", "blue", "yellow", "orange", "purple"}

	result := make([]string, count)
	for i := 0; i < count; i++ {
		randomIndex := rand.Intn(len(colors))
		result[i] = colors[randomIndex]
	}

	return result
}
