package snake

type Point struct {
	X, Y int
}

type Snake struct {
	body      []Point
	direction Point
}

func New(x, y int) *Snake {
	return &Snake{
		body:      []Point{{X: x, Y: y}},
		direction: Point{X: 1, Y: 0},
	}
}

func (s *Snake) Body() []Point {
	return s.body
}

func (s *Snake) Direction() Point {
	return s.direction
}

func (s *Snake) SetDirection(d Point) {
	s.direction = d
}

func (s *Snake) Head() Point {
	return s.body[0]
}

func (s *Snake) Move() {
	head := s.Head()
	newHead := Point{X: head.X + s.direction.X, Y: head.Y + s.direction.Y}
	s.body = append([]Point{newHead}, s.body[:len(s.body)-1]...)
}

func (s *Snake) Grow() {
	head := s.Head()
	newHead := Point{X: head.X + s.direction.X, Y: head.Y + s.direction.Y}
	s.body = append([]Point{newHead}, s.body...)
}

func (s *Snake) CheckWallCollision(gridSize int) bool {
	head := s.Head()
	return head.X < 0 || head.X >= gridSize || head.Y < 0 || head.Y >= gridSize
}

func (s *Snake) CheckSelfCollision() bool {
	head := s.Head()
	for i := 1; i < len(s.body); i++ {
		if s.body[i].X == head.X && s.body[i].Y == head.Y {
			return true
		}
	}
	return false
}

func (s *Snake) CheckReverseDirection(newDir Point) bool {
	return s.direction.X == -newDir.X && s.direction.Y == -newDir.Y
}

func (s *Snake) Occupies(x, y int) bool {
	for _, p := range s.body {
		if p.X == x && p.Y == y {
			return true
		}
	}
	return false
}
