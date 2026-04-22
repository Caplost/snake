# snake-game - Implementation Report v1

## 1. Overview

**Project:** snake-game (贪吃蛇 CLI 游戏)
**Language:** Go
**Version:** v1
**Date:** 2026-04-22

## 2. Acceptance Criteria Status

| Criteria | Status | Details |
|----------|--------|---------|
| 20x20 game grid | ✓ PASS | `constants.GridSize = 20` |
| Snake displayed with `@` or `█` | ✓ PASS | `constants.SnakeSymbol = '█'` |
| Food displayed with `*` | ✓ PASS | `constants.FoodSymbol = '*'` |
| Arrow keys / WASD control | ✓ PASS | `handleKeyEvent()` supports both |
| Score +10 when eating food | ✓ PASS | `g.score += 10` in `Update()` |
| Game over on wall/self collision | ✓ PASS | `CheckWallCollision()` / `CheckSelfCollision()` |
| R key to restart | ✓ PASS | `case 'r', 'R': g.reset()` |
| ~10 FPS (100ms interval) | ✓ PASS | `constants.TickMs = 100` |

## 3. File Structure

```
snake/
├── cmd/snake/main.go              # Entry point
├── internal/
│   ├── constants/constants.go      # GridSize=20, TickMs=100, symbols
│   ├── game/game.go               # Game loop, state management
│   ├── food/food.go               # Food position, generation
│   └── snake/snake.go             # Snake body, direction, movement
├── screenshots/                    # E2E screenshot files
├── scripts/
│   ├── e2e-screenshot.sh           # asciinema/script recording
│   └── asciinema.sh               # Alternative recording
├── go.mod                          # Module declaration
└── go.sum                          # Dependencies
```

## 4. Components

### 4.1 constants.go
```go
const GridSize = 20
const TickMs = 100
const SnakeSymbol = '█'
const FoodSymbol = '*'
const CornerSymbol = '+'
const HorizontalBorder = '-'
const VerticalBorder = '|'
```

### 4.2 snake.go
- `New(x, y int)` — creates snake at position
- `Body() []Point` — returns snake body segments
- `Head() Point` — returns head position
- `Move()` — moves snake forward (append new head, remove tail)
- `Grow()` — moves snake forward (add head, do NOT remove tail)
- `CheckWallCollision(gridSize int)` — wall collision detection
- `CheckSelfCollision()` — self-collision detection
- `CheckReverseDirection(d Point)` — prevents 180° turns
- `Occupies(x, y int)` — checks if position is occupied by snake

### 4.3 food.go
- `New(x, y int)` — creates food at position
- `Position() snake.Point` — returns food position
- `Generate(gridSize int, s *snake.Snake)` — generates food avoiding snake

### 4.4 game.go
- `New()` — initializes game with snake at center and food
- `handleKeyEvent(ev termbox.Event)` — processes keyboard input
- `trySetDirection(d snake.Point)` — sets direction (blocks reverse)
- `Update()` — game tick: movement, collision, eating
- `Render()` — draws game to terminal using termbox
- `reset()` — restarts game
- `Run()` — main game loop

### 4.5 main.go
```go
func main() {
    g := game.New()
    g.Run()
}
```

## 5. Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/nsf/termbox-go` | Terminal handling, keyboard input, rendering |

## 6. Test Coverage

### 6.1 Unit Tests (12 tests)

| Package | Tests |
|---------|-------|
| `snake/internal/food` | 2 tests |
| `snake/internal/game` | 4 tests |
| `snake/internal/snake` | 6 tests |

### 6.2 Test Files

- `snake/internal/food/food_test.go`
- `snake/internal/game/game_test.go`
- `snake/internal/snake/snake_test.go`

## 7. Build & Verification

```bash
# Build
go build ./...           # PASS

# Vet
go vet ./...             # PASS

# Test with race detection
go test -race ./...      # PASS (3 packages, no races)
```

## 8. Known Limitations

- **termbox-go requirement**: Requires a real TTY terminal; headless environments cannot run the game directly
- **Platform dependency**: Uses ANSI escape codes (works on macOS/Linux/Windows Terminal)
- **No E2E automation**: E2E testing requires manual verification via asciinema recordings

## 9. Conclusion

All acceptance criteria met. Implementation is complete and verified.

---
**Status:** ✓ COMPLETE