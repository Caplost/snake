# v4 Code Review - 贪吃蛇命令游戏

**Reviewer**: Agent #1
**Date**: 2026-03-28
**Branch**: feature/sub-feature

---

## Overall Assessment

**VERDICT**: REQUEST_CHANGES

The v4 implementation successfully fixes the CRITICAL input goroutine restart bug and converts polling input to event-driven. The core architectural change is sound — removing the background goroutine eliminates the goroutine lifecycle management problem entirely. However, **critical test coverage gaps** exist for the most important game logic paths.

---

## Test Results Summary

| Check | Result |
|-------|--------|
| `go build ./...` | ✅ PASS |
| `go vet ./...` | ✅ PASS |
| `go test -race ./...` | ✅ PASS (no races) |
| Unit tests | ✅ 12/12 PASS |
| Coverage (snake) | ✅ 87.0% |
| Coverage (food) | ✅ 100.0% |
| Coverage (game) | ⚠️ **15.6%** — too low |

---

## Issues Found

### 1. [HIGH] `reset()` function has zero test coverage

**File**: `snake/internal/game/game.go:155-163`

`reset()` is a critical function — it reinitializes the entire game state on restart. It is completely untested. If `reset()` has a bug, the game will silently break on restart.

**Missing test cases**:
- `reset()` after game over restores initial snake position
- `reset()` clears score to 0
- `reset()` clears gameOver flag
- `reset()` generates new food not overlapping snake
- `reset()` is idempotent (calling twice is safe)

---

### 2. [HIGH] `handleKeyEvent()` has zero test coverage

**File**: `snake/internal/game/game.go:50-81`

`handleKeyEvent()` handles ALL user input: direction keys, WASD, restart (R), and exit (Esc/Ctrl+C). This is the most security-and-correctness-critical path in the game, yet it is completely untested via unit tests.

**Missing test cases**:
- R key calls `reset()` when game over (game.go:76-79)
- R key does nothing when game is NOT over
- Esc/Ctrl+C sets `gameOver = true`
- Arrow keys call `trySetDirection()` with correct direction
- WASD keys call `trySetDirection()` with correct direction
- Non-key events are ignored (ev.Type != termbox.EventKey)

---

### 3. [MEDIUM] `trySetDirection()` has no direct test coverage

**File**: `snake/internal/game/game.go:83-90`

While indirectly covered by snake's `TestSnakeReverseDirection`, the `Game`-level `trySetDirection()` (which wraps with mutex) is not tested. Specifically:
- Mutex protection is not verified
- Integration with `CheckReverseDirection` at Game level not tested

---

### 4. [MEDIUM] Manual test checklist completely unchecked

**File**: `.woali/fp-.../test-v4-1.md:96-111`

All 12 manual test checklist items are unchecked. The test report correctly identifies that termbox-go is not browser-testable, but this makes the missing unit tests for `handleKeyEvent` and `reset` even more critical.

---

### 5. [LOW] `Render()` accesses game state without mutex lock

**File**: `snake/internal/game/game.go:113-153`

`Render()` reads `g.snake.Body()` and `g.food.Position()` without holding `g.mu`, while `Update()` modifies `g.snake` and `g.food` with the lock held. In the current single-threaded event loop this causes no visible bug, but it violates the locking discipline and could become a race condition if the architecture changes.

**Note**: `Render()` does call `g.GameOver()` which locks correctly, but the snake/food reads are unprotected.

---

## What Was Done Well

1. **Architectural fix is correct**: Removing the input goroutine eliminates the "goroutine doesn't restart on reset" bug entirely — a cleaner solution than trying to manage goroutine lifecycle.

2. **Event-driven input**: Using `termbox.PollEvent()` blocking is more efficient than 16ms polling loop.

3. **Mutex protection preserved**: `trySetDirection()`, `setGameOver()`, and `GameOver()` correctly use mutex — no data races detected.

4. **rand.Seed in reset()**: Correctly moved from `New()` to `reset()` so each game starts with a different random seed.

5. **Clean separation of concerns**: `handleKeyEvent()` receives event as parameter (testable interface), unlike the old `HandleInput()` which called `PollEvent()` internally.

---

## Required Changes Before Approval

1. **Add unit tests for `reset()`** — at minimum 3 tests covering: score reset, gameOver flag reset, snake re-initialization.

2. **Add unit tests for `handleKeyEvent()`** — at minimum test the R-key restart path and Esc/Ctrl+C exit path with mocked termbox events.

3. **Increase game package coverage** from 15.6% to at least 40% — current coverage is too low for a package containing critical game logic.

---

## Suggested Test Cases to Add

```go
// In game_test.go

func TestResetClearsScore(t *testing.T) { ... }
func TestResetClearsGameOver(t *testing.T) { ... }
func TestResetReinitializesSnake(t *testing.T) { ... }

func TestHandleKeyEvent_RestartsGame(t *testing.T) {
    ev := termbox.Event{Type: termbox.EventKey, Ch: 'r'}
    // setup game in gameOver state
    g.handleKeyEvent(ev)
    // verify reset was called
}

func TestHandleKeyEvent_ExitOnEsc(t *testing.T) {
    ev := termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEsc}
    g.handleKeyEvent(ev)
    if !g.GameOver() { t.Error("Esc should trigger game over") }
}

func TestHandleKeyEvent_IgnoresNonKeyEvents(t *testing.T) {
    ev := termbox.Event{Type: termbox.EventResize}
    g.handleKeyEvent(ev)
    // game state should be unchanged
}
```

---

## Conclusion

The v4 code is **structurally sound** — the architectural fix for the goroutine bug is correct and well-reasoned. However, the **complete absence of tests for `reset()` and `handleKeyEvent()`** leaves the most critical user journeys (restart, exit, direction input) verified only by manual testing. Given the game cannot be e2e-tested in CI, unit tests for these functions are essential before merge.
