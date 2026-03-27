# Code Review - Snake CLI Game (v1)

## Overall Assessment: Request Changes

The implementation is functional and tests pass, but there are **critical concurrency bugs** and several high/medium issues that must be addressed before approval.

---

## CRITICAL Issues

### 1. Data Race on Snake Direction (game.go:158-163 + game.go:83-102)

**Severity:** CRITICAL
**File:** `snake/internal/game/game.go`
**Lines:** 41-81 (HandleInput/trySetDirection), 158-163 (goroutine), 83-102 (Update)

A goroutine continuously calls `HandleInput()` which calls `trySetDirection()` that modifies `g.snake.direction`. The main loop's `Update()` reads `g.snake.direction` via `s.direction.X/Y` in `Move()`. These happen concurrently without any synchronization (no mutex).

**Race scenario:** `Update()` reads direction while `trySetDirection()` writes it → data race → undefined behavior.

**Fix:** Use a mutex to protect direction changes, or restructure so input handling happens synchronously in the main loop (remove the goroutine).

---

### 2. Data Race on gameOver Flag (game.go:158-163 + game.go:88-90)

**Severity:** CRITICAL
**File:** `snake/internal/game/game.go`
**Lines:** 89 (write), 159 (read in goroutine loop condition)

The goroutine reads `g.gameOver` at line 159 while the main loop writes it at line 89. No synchronization.

---

### 3. No Test for Self-Collision Detection

**Severity:** HIGH
**File:** `snake/internal/snake/snake_test.go`
**Lines:** (missing)

`CheckSelfCollision()` is implemented but never tested. This is a core game-over condition per the plan. Without a test, we cannot verify the snake dies when it hits itself.

**Fix:** Add `TestSnakeSelfCollision` test case.

---

## HIGH Issues

### 4. rand.Seed Not Called on Game Reset (game.go:138-143)

**Severity:** HIGH
**File:** `snake/internal/game/game.go`
**Line:** 139

`rand.Seed(time.Now().UnixNano())` is called only in `New()` (line 23), but `reset()` (called on restart) does NOT call `New()`. It directly creates `snake.New()` and calls `food.Generate()`. After restart, the random sequence will repeat the same pattern.

**Fix:** Move `rand.Seed()` into `reset()`, or call `New()` instead of inline reset logic.

---

### 5. Goroutine Never Terminates (game.go:158-163)

**Severity:** MEDIUM
**File:** `snake/internal/game/game.go`
**Lines:** 158-163

The input goroutine runs `for !g.gameOver` and is never stopped. When the game ends, the goroutine continues looping (calling `HandleInput()` → `PollEvent()`) until the process exits. No `sync.WaitGroup` or channel stop mechanism.

**Fix:** Use a stop channel or `sync.WaitGroup` to properly terminate the goroutine.

---

### 6. Input Goroutine Timing Doesn't Match Ticker (game.go:158-163)

**Severity:** MEDIUM
**File:** `snake/internal/game/game.go`
**Lines:** 155-163

The input goroutine calls `HandleInput()` then sleeps 16ms, regardless of the 100ms game tick. Input is checked ~6 times per game update, which is wasteful. More importantly, if `PollEvent()` blocks for longer than expected (e.g., terminal issues), the goroutine starves.

---

### 7. Unused Constant (constants.go:8)

**Severity:** LOW
**File:** `snake/internal/constants/constants.go`
**Line:** 8

`EmptySymbol` is defined but never used in rendering or anywhere else.

**Fix:** Remove or use it for empty grid cells if that's the intended semantics.

---

### 8. Test Names Don't Match Test Content (game_test.go:24-31)

**Severity:** LOW
**File:** `snake/internal/game/game_test.go`
**Lines:** 24-31

`TestGameOverOnWallCollision` creates a standalone snake, moves it, then tests `s.CheckWallCollision()`. It never creates a `Game` struct or tests `Game.GameOver`. The name is misleading.

---

## Summary

| Issue | Severity | Status |
|-------|----------|--------|
| Data race: direction | CRITICAL | Must fix |
| Data race: gameOver | CRITICAL | Must fix |
| Missing self-collision test | HIGH | Must fix |
| rand.Seed on reset | HIGH | Must fix |
| Goroutine lifecycle | MEDIUM | Should fix |
| Goroutine timing | MEDIUM | Should fix |
| Unused EmptySymbol | LOW | Nice to fix |
| Misleading test name | LOW | Nice to fix |

**Recommendation:** Request changes. The concurrency bugs are blockers. All tests pass but they don't cover the critical self-collision path or the reset randomness issue.
