package food

import (
	"math/rand"

	"snake/internal/snake"
)

type Food struct {
	position snake.Point
}

func New(x, y int) *Food {
	return &Food{
		position: snake.Point{X: x, Y: y},
	}
}

func (f *Food) Position() snake.Point {
	return f.position
}

func (f *Food) Generate(gridSize int, s *snake.Snake) {
	for {
		x := rand.Intn(gridSize)
		y := rand.Intn(gridSize)
		if !s.Occupies(x, y) {
			f.position = snake.Point{X: x, Y: y}
			return
		}
	}
}
