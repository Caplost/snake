package game

import (
	"fmt"
	"math/rand"
	"time"

	"snake/internal/constants"
	"snake/internal/food"
	"snake/internal/snake"

	"github.com/nsf/termbox-go"
)

type Game struct {
	score    int
	gameOver bool
	snake    *snake.Snake
	food     *food.Food
}

func New() *Game {
	rand.Seed(time.Now().UnixNano())
	s := snake.New(constants.GridSize/2, constants.GridSize/2)
	f := food.New(0, 0)
	f.Generate(constants.GridSize, s)
	return &Game{
		snake: s,
		food:  f,
	}
}

func (g *Game) Score() int {
	return g.score
}

func (g *Game) GameOver() bool {
	return g.gameOver
}

func (g *Game) HandleInput() {
	ev := termbox.PollEvent()
	if ev.Type == termbox.EventKey {
		switch ev.Key {
		case termbox.KeyArrowUp:
			g.trySetDirection(snake.Point{X: 0, Y: -1})
		case termbox.KeyArrowDown:
			g.trySetDirection(snake.Point{X: 0, Y: 1})
		case termbox.KeyArrowLeft:
			g.trySetDirection(snake.Point{X: -1, Y: 0})
		case termbox.KeyArrowRight:
			g.trySetDirection(snake.Point{X: 1, Y: 0})
		case termbox.KeyCtrlC, termbox.KeyEsc:
			g.gameOver = true
		}
		if g.gameOver && (ev.Ch == 'r' || ev.Ch == 'R') {
			g.reset()
		}
		switch ev.Ch {
		case 'w', 'W':
			g.trySetDirection(snake.Point{X: 0, Y: -1})
		case 's', 'S':
			g.trySetDirection(snake.Point{X: 0, Y: 1})
		case 'a', 'A':
			g.trySetDirection(snake.Point{X: -1, Y: 0})
		case 'd', 'D':
			g.trySetDirection(snake.Point{X: 1, Y: 0})
		case 'r', 'R':
			if g.gameOver {
				g.reset()
			}
		}
	}
}

func (g *Game) trySetDirection(d snake.Point) {
	if g.snake.CheckReverseDirection(d) {
		return
	}
	g.snake.SetDirection(d)
}

func (g *Game) Update() {
	if g.gameOver {
		return
	}

	if g.snake.CheckWallCollision(constants.GridSize) || g.snake.CheckSelfCollision() {
		g.gameOver = true
		return
	}

	head := g.snake.Head()
	foodPos := g.food.Position()
	if head.X == foodPos.X && head.Y == foodPos.Y {
		g.snake.Grow()
		g.score += 10
		g.food.Generate(constants.GridSize, g.snake)
	} else {
		g.snake.Move()
	}
}

func (g *Game) Render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw corners
	termbox.SetCell(0, 0, constants.CornerSymbol, termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(constants.GridSize, 0, constants.CornerSymbol, termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(0, constants.GridSize, constants.CornerSymbol, termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(constants.GridSize, constants.GridSize, constants.CornerSymbol, termbox.ColorWhite, termbox.ColorDefault)

	// Draw horizontal borders (top and bottom)
	for x := 1; x < constants.GridSize; x++ {
		termbox.SetCell(x, 0, constants.HorizontalBorder, termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, constants.GridSize, constants.HorizontalBorder, termbox.ColorWhite, termbox.ColorDefault)
	}
	// Draw vertical borders (left and right)
	for y := 1; y < constants.GridSize; y++ {
		termbox.SetCell(0, y, constants.VerticalBorder, termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(constants.GridSize, y, constants.VerticalBorder, termbox.ColorWhite, termbox.ColorDefault)
	}

	for _, p := range g.snake.Body() {
		termbox.SetCell(p.X+1, p.Y+1, constants.SnakeSymbol, termbox.ColorGreen, termbox.ColorDefault)
	}

	foodPos := g.food.Position()
	termbox.SetCell(foodPos.X+1, foodPos.Y+1, constants.FoodSymbol, termbox.ColorRed, termbox.ColorDefault)

	if g.gameOver {
		msg := fmt.Sprintf("Game Over! Score: %d. Press R to restart", g.score)
		for i, c := range msg {
			termbox.SetCell(constants.GridSize/2-10+i, constants.GridSize/2, c, termbox.ColorYellow, termbox.ColorDefault)
		}
	}

	scoreMsg := fmt.Sprintf("Score: %d", g.score)
	for i, c := range scoreMsg {
		termbox.SetCell(2+i, 0, c, termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.Flush()
}

func (g *Game) reset() {
	g.snake = snake.New(constants.GridSize/2, constants.GridSize/2)
	g.food.Generate(constants.GridSize, g.snake)
	g.score = 0
	g.gameOver = false
}

func (g *Game) Run() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)

	g.reset()

	ticker := time.NewTicker(time.Duration(constants.TickMs) * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for !g.gameOver {
			g.HandleInput()
			time.Sleep(16 * time.Millisecond)
		}
	}()

	for {
		g.Render()
		if g.gameOver {
			for {
				ev := termbox.PollEvent()
				if ev.Type == termbox.EventKey && (ev.Ch == 'r' || ev.Ch == 'R' || ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyEsc) {
					if ev.Ch == 'r' || ev.Ch == 'R' {
						g.reset()
						break
					}
					return
				}
			}
		}
		<-ticker.C
		g.Update()
	}
}
