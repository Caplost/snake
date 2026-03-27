package snake

import (
	"testing"
)

func TestSnakeMove(t *testing.T) {
	s := New(5, 5)
	initialHead := s.Head()
	s.SetDirection(Point{X: 1, Y: 0})
	s.Move()
	newHead := s.Head()
	if newHead.X != initialHead.X+1 || newHead.Y != initialHead.Y {
		t.Errorf("Expected head at (%d, %d), got (%d, %d)", initialHead.X+1, initialHead.Y, newHead.X, newHead.Y)
	}
}

func TestSnakeGrow(t *testing.T) {
	s := New(5, 5)
	initialLen := len(s.body)
	s.SetDirection(Point{X: 1, Y: 0})
	s.Grow()
	if len(s.body) != initialLen+1 {
		t.Errorf("Expected length %d, got %d", initialLen+1, len(s.body))
	}
}

func TestSnakeReverseDirection(t *testing.T) {
	s := New(5, 5)
	s.SetDirection(Point{X: 1, Y: 0})
	if !s.CheckReverseDirection(Point{X: -1, Y: 0}) {
		t.Error("Expected reverse direction to be detected")
	}
	if s.CheckReverseDirection(Point{X: 0, Y: -1}) {
		t.Error("90-degree turn should be allowed")
	}
}

func TestSnakeWallCollision(t *testing.T) {
	s := New(0, 0)
	s.SetDirection(Point{X: -1, Y: 0})
	s.Move()
	if !s.CheckWallCollision(20) {
		t.Error("Expected wall collision after moving left from x=0")
	}
}

func TestSnakeOccupies(t *testing.T) {
	s := New(5, 5)
	if !s.Occupies(5, 5) {
		t.Error("Expected snake to occupy its head position")
	}
	if s.Occupies(6, 6) {
		t.Error("Snake should not occupy position (6,6)")
	}
}

func TestSnakeSelfCollision(t *testing.T) {
	s := New(5, 5)
	s.SetDirection(Point{X: 1, Y: 0})
	// Grow snake to length 3
	s.Grow()
	s.Grow()
	// Now body should be: [head at (5,5), (4,5), (3,5)]
	// Force a self-collision scenario by manually creating overlapping segments
	s.body = []Point{{X: 5, Y: 5}, {X: 5, Y: 5}, {X: 4, Y: 5}}
	if !s.CheckSelfCollision() {
		t.Error("Expected self collision when head overlaps with body")
	}
}
