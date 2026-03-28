# Test Specification - Snake CLI Game v3

## Overview

This test specification covers the v3 implementation of the Snake CLI game, which fixes critical concurrency bugs identified in v2. The primary changes are adding mutex protection for shared state and WaitGroup for proper goroutine cleanup.

## Project Structure

```
snake/
├── cmd/snake/main.go
├── internal/
│   ├── constants/constants.go
│   ├── snake/snake.go + snake_test.go
│   ├── food/food.go + food_test.go
│   └── game/game.go + game_test.go
├── go.mod
└── go.sum
```

## Environment & Setup

- **Language**: Go 1.21+
- **Dependencies**: github.com/nsf/termbox-go v1.1.1
- **Build**: `go build -o snake ./cmd/snake`
- **Test**: `go test ./... -v`
- **Lint**: `go vet ./...`
- **Race Detection**: `go test -race ./... -v`
- **Run**: `./snake` (requires terminal with ANSI support)

## Unit Test Coverage

### snake/snake_test.go

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestSnakeMove` | Snake moves in set direction | Head position updates correctly |
| `TestSnakeGrow` | Snake grows when eating food | Body length increases by 1 |
| `TestSnakeReverseDirection` | 180-degree turn is rejected | CheckReverseDirection returns true for opposite direction |
| `TestSnakeWallCollision` | Snake detects wall collision | CheckWallCollision returns true when head exits grid |
| `TestSnakeOccupies` | Snake occupancy detection | Occupies returns true for snake body positions |
| `TestSnakeSelfCollision` | Snake detects self collision | CheckSelfCollision returns true when head overlaps body |

### food/food_test.go

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestFoodGenerateNotOnSnake` | Food never spawns on snake | 100 iterations all produce valid positions |
| `TestFoodPosition` | Food position getter | Returns correct x,y coordinates |

### game/game_test.go

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestScoreIncrease` | Initial score is zero | Score starts at 0 |
| `TestScoreIncreaseOnEatingFood` | Score increases when eating food | Score increases by 10 |
| `TestGameOverOnWallCollision` | Game ends on wall collision | Game over flag set correctly |
| `TestNewGameInitialization` | New game initializes correctly | Score=0, gameOver=false, snake/food initialized |

## Integration Test Scenarios

### Scenario 1: Game Loop Lifecycle
1. Start game → Snake appears at center, food at random location
2. Move snake → Direction updates based on input
3. Eat food → Score +10, snake grows, new food spawns
4. Hit wall → Game over state triggered
5. Press R → Game resets to initial state

### Scenario 2: Input Handling
- Arrow keys (Up/Down/Left/Right) change snake direction
- WASD keys change snake direction
- 180-degree reverse is rejected
- Ctrl+C or Esc exits game

### Scenario 3: Collision Detection
- Wall collision: Head x < 0, x >= GridSize, y < 0, y >= GridSize
- Self collision: Head overlaps any body segment

### Scenario 4: Concurrency Safety (v3 Focus)
1. Start game with `go test -race` running
2. Rapidly press direction keys (WASD/Arrows)
3. Game should not panic or deadlock
4. Race detector should report no races

## Expected Results

### UI Layout (unchanged from v2)
```
+--------------------+
| SNAKE GAME    Score: 0 |
+--------------------+
|                    |
|    ████ (snake)    |
|    * (food)        |
|                    |
+--------------------+
| WASD/Arrows to move |
+--------------------+
```

### Game Rules (unchanged from v2)
- Grid: 20x20 cells
- Snake symbol: `█` (Unicode block)
- Food symbol: `*`
- Frame rate: ~10 FPS (100ms tick)
- Score per food: +10
- Border: corners `+`, horizontal `-`, vertical `|`

## Edge Cases

| Edge Case | Handling |
|-----------|----------|
| Food on snake | Regenerate until valid position found |
| 180-degree turn | Reject direction change |
| Window resize | Fixed 20x20 grid, no resize handling |
| Game over input | Wait for R key to restart |
| Rapid key presses | Input handled in separate goroutine with mutex protection |
| Concurrent state access | Mutex-protected getters/setters for gameOver |

## Test Commands

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run with race detector (CRITICAL for v3)
go test -race ./... -v

# Run linter
go vet ./...

# Build binary
go build -o snake ./cmd/snake

# Run game (requires terminal)
./snake
```

## Race Detection Verification (v3)

The most important verification for v3 is running with `-race` flag:

```bash
$ go test -race ./...
```

Expected: **No race warnings**. The following accesses are now protected:
- `g.gameOver` reads via `GameOver()` getter (mutex locked)
- `g.gameOver` writes via `setGameOver()` setter (mutex locked)
- `g.snake.direction` writes via `trySetDirection()` (mutex locked)
- `g.snake.direction` reads via `Move()` - protected by caller's mutex in Update()

## Coverage Report

| Package | Coverage |
|---------|----------|
| snake/internal/food | 100.0% |
| snake/internal/snake | 87.0% |
| snake/internal/game | 14.6% |
| snake/cmd/snake | 0.0% (no test files) |

**Note:** The game package has lower coverage because the Run() function (which contains the concurrent code paths) requires a terminal and cannot be tested in unit tests. The race detector is the primary verification for concurrent correctness.