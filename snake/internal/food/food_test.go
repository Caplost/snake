package food

import (
	"testing"

	"snake/internal/snake"
)

func TestFoodGenerateNotOnSnake(t *testing.T) {
	s := snake.New(5, 5)
	f := &Food{position: snake.Point{X: 0, Y: 0}}

	for i := 0; i < 100; i++ {
		f.Generate(20, s)
		if s.Occupies(f.position.X, f.position.Y) {
			t.Errorf("Food generated on snake at (%d, %d)", f.position.X, f.position.Y)
		}
	}
}

func TestFoodPosition(t *testing.T) {
	f := New(3, 4)
	pos := f.Position()
	if pos.X != 3 || pos.Y != 4 {
		t.Errorf("Expected position (3,4), got (%d,%d)", pos.X, pos.Y)
	}
}
