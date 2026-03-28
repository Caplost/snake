# Implementation Plan - Snake CLI Game v3

## 1. Overview

v3 addresses critical concurrency bugs identified in the v2 code review. The game currently has data races that can cause undefined behavior during actual gameplay.

## 2. Issues to Fix

### CRITICAL Issues

#### 1. Data Race on Snake Direction
- **File:** `snake/internal/game/game.go`
- **Problem:** Goroutine at lines 166-171 writes `g.snake.direction` via `HandleInput()` → `trySetDirection()`. Main loop reads it via `g.snake.Move()`.
- **Fix:** Add `sync.Mutex` to protect direction access.

#### 2. Data Race on gameOver Flag
- **File:** `snake/internal/game/game.go`
- **Problem:** Goroutine reads `g.gameOver` as loop condition, main loop writes it.
- **Fix:** Use mutex to protect gameOver reads/writes.

### HIGH Issues

#### 3. rand.Seed Not Called on Game Reset
- **File:** `snake/internal/game/game.go`
- **Problem:** `rand.Seed()` is called only in `New()`, not in `reset()`. After restart, food generation is predictable.
- **Fix:** Call `rand.Seed()` in `reset()` function.

### MEDIUM Issues

#### 4. Goroutine Never Terminates
- **File:** `snake/internal/game/game.go`
- **Problem:** Input goroutine runs until process exit, no cleanup.
- **Fix:** Use `sync.WaitGroup` to properly terminate goroutine on game exit.

### LOW Issues

#### 5. Unused Constant EmptySymbol
- **File:** `snake/internal/constants/constants.go`
- **Fix:** Remove unused constant.

## 3. Implementation Changes

### game.go Changes
1. Add `sync.Mutex` field to `Game` struct
2. Add `sync.WaitGroup` field to `Game` struct
3. Protect `gameOver` flag with mutex in all accesses
4. Protect `snake.direction` with mutex in all accesses
5. Move `rand.Seed()` from `New()` to `reset()`
6. Use WaitGroup to signal goroutine termination

### constants.go Changes
1. Remove `EmptySymbol` constant (unused)

## 4. Verification

```bash
# Build
go build ./...

# Run tests with race detector
go test -race ./... -v

# Lint
go vet ./...
```

## 5. Files to Modify

| File | Changes |
|------|---------|
| `snake/internal/game/game.go` | Add mutex/waitgroup, protect shared state, fix rand.Seed |
| `snake/internal/constants/constants.go` | Remove EmptySymbol |
| `snake/internal/game/game_test.go` | May need updates if internal state access changes |

## 6. Expected Result

- All CRITICAL and HIGH issues resolved
- `go test -race` passes without warnings
- Game restarts with fresh random sequence
- Goroutine properly cleaned up on exit