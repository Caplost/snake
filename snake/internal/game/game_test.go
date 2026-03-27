package game

import (
	"testing"

	"snake/internal/constants"
	"snake/internal/food"
	"snake/internal/snake"
)

func TestScoreIncrease(t *testing.T) {
	g := &Game{
		score:    0,
		gameOver: false,
		snake:    snake.New(constants.GridSize/2, constants.GridSize/2),
		food:     food.New(0, 0),
	}

	if g.score != 0 {
		t.Errorf("Initial score should be 0, got %d", g.score)
	}
}

func TestGameOverOnWallCollision(t *testing.T) {
	s := snake.New(0, 0)
	s.SetDirection(snake.Point{X: -1, Y: 0})
	s.Move()
	if !s.CheckWallCollision(constants.GridSize) {
		t.Error("Snake should be in wall collision state after moving left from x=0")
	}
}

func TestNewGameInitialization(t *testing.T) {
	g := New()
	if g.score != 0 {
		t.Errorf("New game score should be 0, got %d", g.score)
	}
	if g.gameOver {
		t.Error("New game should not be game over")
	}
	if g.snake == nil {
		t.Error("Snake should be initialized")
	}
	if g.food == nil {
		t.Error("Food should be initialized")
	}
}
