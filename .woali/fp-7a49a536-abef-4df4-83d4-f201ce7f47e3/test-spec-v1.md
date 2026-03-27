# Test Specification - Snake CLI Game (v1)

## Unit Coverage

### `snake/internal/snake/snake_test.go`

| Test | Description | Status |
|------|-------------|--------|
| TestSnakeMove | Verifies snake moves in set direction | PASS |
| TestSnakeGrow | Verifies snake grows by 1 when Grow() called | PASS |
| TestSnakeReverseDirection | Verifies 180° turn is detected, 90° allowed | PASS |
| TestSnakeWallCollision | Verifies collision detected after moving off grid | PASS |
| TestSnakeOccupies | Verifies Occupies() returns correct positions | PASS |

### `snake/internal/food/food_test.go`

| Test | Description | Status |
|------|-------------|--------|
| TestFoodGenerateNotOnSnake | Verifies food never generates on snake body | PASS |
| TestFoodPosition | Verifies food position getter works | PASS |

### `snake/internal/game/game_test.go`

| Test | Description | Status |
|------|-------------|--------|
| TestScoreIncrease | Verifies initial score is 0 | PASS |
| TestGameOverOnWallCollision | Verifies game detects wall collision | PASS |
| TestNewGameInitialization | Verifies all game fields initialize correctly | PASS |

## Integration/E2E Scenarios

### Manual Testing Required
Since termbox-go requires a real terminal, the following scenarios should be tested manually:

1. **Game Start**: Snake appears at center, food at random location
2. **Movement**: Arrow keys and WASD control direction
3. **Eating Food**: Snake grows +1, score +10
4. **Wall Collision**: Game ends when snake hits boundary
5. **Self Collision**: Game ends when snake hits itself
6. **Reverse Direction**: 180° turn is prevented
7. **Restart**: R key restarts after game over
8. **Exit**: Ctrl+C or ESC exits game

## Environment/Setup

### Build
```bash
cd snake
go mod tidy
go build -o snake ./cmd/snake
```

### Run
```bash
./snake
```

### Test
```bash
go test ./... -v
```

## Coverage Summary
- Lines: ~200 LOC across 6 source files
- Test coverage: 11 passing tests
- Packages: 4 (cmd, constants, snake, food, game)
