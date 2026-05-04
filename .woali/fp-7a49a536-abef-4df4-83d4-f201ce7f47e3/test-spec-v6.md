# Test Specification v6

## Sub-feature
贪吃蛇命令游戏开发 (v6)

## Scope
v6 is a test-coverage-only iteration. No functional changes. The test suite verifies:
1. `reset()` behavior (score, gameOver, snake position, food placement)
2. `handleKeyEvent()` behavior (restart, exit, ignore, direction keys)
3. `trySetDirection()` behavior (blocks reverse, allows valid)
4. `Render()` no longer has unprotected concurrent reads (verified by race detector)

---

## Unit Tests

### Package: `snake/internal/game`

#### Existing Tests (retained from v5)
- `TestScoreIncrease` — verify initial score is 0
- `TestScoreIncreaseOnEatingFood` — score +10 when eating food
- `TestGameOverOnWallCollision` — wall collision detection
- `TestNewGameInitialization` — New() initializes all fields

#### New Tests (v6)

**reset() coverage**
- `TestResetClearsScore`
  - Setup: `g := New(); g.score = 100; g.reset()`
  - Expected: `g.score == 0`

- `TestResetClearsGameOver`
  - Setup: `g := New(); g.setGameOver(true); g.reset()`
  - Expected: `g.gameOver == false`

- `TestResetReinitializesSnake`
  - Setup: `g := New(); g.snake.Move(); g.reset()`
  - Expected: snake head at `(GridSize/2, GridSize/2)`

- `TestResetGeneratesNewFood`
  - Setup: `g := New(); g.reset()`
  - Expected: `g.food.Position()` not in `g.snake.Body()`

**handleKeyEvent() coverage**
- `TestHandleKeyEvent_RestartsGame`
  - Setup: `g := New(); g.setGameOver(true); g.score = 999`
  - Input: `termbox.Event{Type: EventKey, Ch: 'r'}`
  - Expected: `g.gameOver == false && g.score == 0`

- `TestHandleKeyEvent_ExitOnEsc`
  - Setup: `g := New()`
  - Input: `termbox.Event{Type: EventKey, Key: KeyEsc}`
  - Expected: `g.GameOver() == true`

- `TestHandleKeyEvent_ExitOnCtrlC`
  - Setup: `g := New()`
  - Input: `termbox.Event{Type: EventKey, Key: KeyCtrlC}`
  - Expected: `g.GameOver() == true`

- `TestHandleKeyEvent_IgnoresNonKeyEvents`
  - Setup: `g := New(); g.score = 123`
  - Input: `termbox.Event{Type: EventResize}`
  - Expected: `g.GameOver() == false && g.score == 123`

- `TestHandleKeyEvent_WASDAndArrows`
  - Subtests (fresh Game each): ArrowUp, W → direction {0,-1}; ArrowDown, S → direction {0,1}; ArrowRight, D → direction {1,0}
  - Expected: direction matches input

- `TestHandleKeyEvent_LeftBlockedFromRight`
  - Setup: `g := New()` (initial direction right)
  - Input: `KeyArrowLeft`
  - Expected: direction stays {1,0}

- `TestHandleKeyEvent_ABlockedFromRight`
  - Setup: `g := New()` (initial direction right)
  - Input: `Ch: 'a'`
  - Expected: direction stays {1,0}

**trySetDirection() coverage**
- `TestTrySetDirection_BlocksReverse`
  - Setup: `g := New()` (direction right)
  - Input: `{X: -1, Y: 0}`
  - Expected: direction stays {1,0}

- `TestTrySetDirection_AllowsValidDirection`
  - Setup: `g := New()` (direction right)
  - Input: `{X: 0, Y: -1}`
  - Expected: direction becomes {0,-1}

### Package: `snake/internal/food`
- `TestFoodGenerateNotOnSnake`
- `TestFoodPosition`

### Package: `snake/internal/snake`
- `TestSnakeMove`
- `TestSnakeGrow`
- `TestSnakeReverseDirection`
- `TestSnakeWallCollision`
- `TestSnakeOccupies`
- `TestSnakeSelfCollision`

---

## Integration / E2E Scenarios

This project is a CLI terminal game using `termbox-go`. E2E testing requires a real terminal and is not suitable for automated CI. Manual E2E verification was performed in prior iterations (v1-v5) using `scripts/e2e-screenshot.sh`.

For v6, no E2E scenarios are required because:
- No functional behavior changed
- All changes are internal (mutex, tests)
- Playwright is not applicable (terminal CLI, not web)

---

## Environment / Setup Notes

```bash
cd /Users/wangyinneng/SaaS/Woali/workspace/room-ab68f9ee-f954-4ec9-a486-84e0d9d86bca/feature-814b09c0-0b8e-4133-b1bd-250edb0ec936/snake

# Dependencies
# github.com/nsf/termbox-go (already in go.mod)

# Run tests
go test ./... -v

# Coverage
go test ./... -cover

# Race detection
go test -race ./...

# Lint
go vet ./...
```

## Expected Results

| Metric | Expected |
|--------|----------|
| All unit tests | PASS |
| game package coverage | >= 40% |
| Race detector | No races |
| go vet | No issues |

## Playwright Coverage

N/A — this is a Go CLI application, not a web application. No Playwright tests exist or are required.
