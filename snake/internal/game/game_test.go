package game

import (
	"testing"

	"snake/internal/constants"
	"snake/internal/food"
	"snake/internal/snake"

	"github.com/nsf/termbox-go"
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

func TestScoreIncreaseOnEatingFood(t *testing.T) {
	g := New()
	initialScore := g.score
	head := g.snake.Head()
	g.food = food.New(head.X, head.Y)
	g.Update()
	if g.score != initialScore+10 {
		t.Errorf("Score should increase by 10 when eating food, got %d", g.score)
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

// --- reset() tests ---

func TestResetClearsScore(t *testing.T) {
	g := New()
	g.score = 100
	g.reset()
	if g.score != 0 {
		t.Errorf("reset() should clear score, got %d", g.score)
	}
}

func TestResetClearsGameOver(t *testing.T) {
	g := New()
	g.setGameOver(true)
	g.reset()
	if g.gameOver {
		t.Error("reset() should clear gameOver flag")
	}
}

func TestResetReinitializesSnake(t *testing.T) {
	g := New()
	// Move snake away from center
	g.snake.Move()
	g.reset()
	head := g.snake.Head()
	expectedX := constants.GridSize / 2
	if head.X != expectedX || head.Y != expectedX {
		t.Errorf("reset() should reinit snake at center (%d,%d), got (%d,%d)",
			expectedX, expectedX, head.X, head.Y)
	}
}

func TestResetGeneratesNewFood(t *testing.T) {
	g := New()
	g.reset()
	snakeBody := g.snake.Body()
	foodPos := g.food.Position()
	for _, p := range snakeBody {
		if p.X == foodPos.X && p.Y == foodPos.Y {
			t.Error("reset() should generate food not overlapping snake")
		}
	}
}

// --- handleKeyEvent() tests ---

func TestHandleKeyEvent_RestartsGame(t *testing.T) {
	g := New()
	g.setGameOver(true)
	g.score = 999
	ev := termbox.Event{Type: termbox.EventKey, Ch: 'r'}
	g.handleKeyEvent(ev)
	if g.gameOver {
		t.Error("R key should restart game (clear gameOver)")
	}
	if g.score != 0 {
		t.Errorf("R key should reset score, got %d", g.score)
	}
}

func TestHandleKeyEvent_ExitOnEsc(t *testing.T) {
	g := New()
	ev := termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEsc}
	g.handleKeyEvent(ev)
	if !g.GameOver() {
		t.Error("Esc should trigger game over")
	}
}

func TestHandleKeyEvent_ExitOnCtrlC(t *testing.T) {
	g := New()
	ev := termbox.Event{Type: termbox.EventKey, Key: termbox.KeyCtrlC}
	g.handleKeyEvent(ev)
	if !g.GameOver() {
		t.Error("Ctrl+C should trigger game over")
	}
}

func TestHandleKeyEvent_IgnoresNonKeyEvents(t *testing.T) {
	g := New()
	g.score = 123
	ev := termbox.Event{Type: termbox.EventResize}
	g.handleKeyEvent(ev)
	if g.GameOver() || g.score != 123 {
		t.Error("Non-key events should be ignored")
	}
}

func TestHandleKeyEvent_WASDAndArrows(t *testing.T) {
	// Initial direction is right {1,0}. Test only non-reverse directions.
	tests := []struct {
		name     string
		ev       termbox.Event
		expected snake.Point
	}{
		{"ArrowUp", termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowUp}, snake.Point{X: 0, Y: -1}},
		{"W", termbox.Event{Type: termbox.EventKey, Ch: 'w'}, snake.Point{X: 0, Y: -1}},
		{"ArrowDown", termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowDown}, snake.Point{X: 0, Y: 1}},
		{"S", termbox.Event{Type: termbox.EventKey, Ch: 's'}, snake.Point{X: 0, Y: 1}},
		{"ArrowRight", termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowRight}, snake.Point{X: 1, Y: 0}},
		{"D", termbox.Event{Type: termbox.EventKey, Ch: 'd'}, snake.Point{X: 1, Y: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			g.handleKeyEvent(tt.ev)
			dir := g.snake.Direction()
			if dir != tt.expected {
				t.Errorf("key event should set direction to %v, got %v", tt.expected, dir)
			}
		})
	}
}

func TestHandleKeyEvent_LeftBlockedFromRight(t *testing.T) {
	g := New()
	// Initial direction is right {1,0}; left is reverse and should be blocked
	g.handleKeyEvent(termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowLeft})
	dir := g.snake.Direction()
	if dir.X != 1 || dir.Y != 0 {
		t.Errorf("left should be blocked when moving right, got %v", dir)
	}
}

func TestHandleKeyEvent_ABlockedFromRight(t *testing.T) {
	g := New()
	g.handleKeyEvent(termbox.Event{Type: termbox.EventKey, Ch: 'a'})
	dir := g.snake.Direction()
	if dir.X != 1 || dir.Y != 0 {
		t.Errorf("A should be blocked when moving right, got %v", dir)
	}
}

// --- trySetDirection() tests ---

func TestTrySetDirection_BlocksReverse(t *testing.T) {
	g := New()
	// Current direction is (1,0) - moving right
	// Try to go left (reverse) - should be blocked
	g.trySetDirection(snake.Point{X: -1, Y: 0})
	dir := g.snake.Direction()
	if dir.X != 1 || dir.Y != 0 {
		t.Errorf("trySetDirection should block reverse direction, got %v", dir)
	}
}

func TestTrySetDirection_AllowsValidDirection(t *testing.T) {
	g := New()
	// Current direction is (1,0) - moving right
	// Try to go up - should be allowed
	g.trySetDirection(snake.Point{X: 0, Y: -1})
	dir := g.snake.Direction()
	if dir.X != 0 || dir.Y != -1 {
		t.Errorf("trySetDirection should allow valid direction, got %v", dir)
	}
}
