# Code Review - Snake CLI Game v2 (Review Agent #1)

## Overall Assessment: REQUEST_CHANGES

The v2 implementation correctly adds the ASCII border rendering (corners `+`, horizontal `-`, vertical `|`) and adds two new tests (`TestSnakeSelfCollision`, `TestScoreIncreaseOnEatingFood`) as specified. However, **all CRITICAL and HIGH concurrency/reset bugs from the v1 review remain unfixed**. These are blocking issues.

---

## CRITICAL Issues

### 1. Data Race on Snake Direction (UNRESOLVED from v1)

**Severity:** CRITICAL
**File:** `snake/internal/game/game.go`
**Lines:** 76-81 (trySetDirection writes), 35-38 (Move reads via s.direction), 166-171 (goroutine)

The goroutine at lines 166-171 continuously calls `HandleInput()` → `trySetDirection()` which **writes** `g.snake.direction`. The main loop at line 100 calls `g.snake.Move()` which **reads** `g.snake.direction` via `s.direction.X/Y`. These happen concurrently without any mutex or synchronization.

`go test -race` does NOT catch this because the test suite never exercises the actual concurrent execution path (goroutine + main loop running simultaneously).

**Fix:** Use a `sync.Mutex` protecting `g.snake.direction`, or restructure `Run()` to handle input synchronously in the main loop (remove the goroutine).

---

### 2. Data Race on gameOver Flag (UNRESOLVED from v1)

**Severity:** CRITICAL
**File:** `snake/internal/game/game.go`
**Lines:** 89 (write), 167 (read in goroutine loop condition)

The goroutine at line 167 reads `g.gameOver` as its loop condition `for !g.gameOver`. The main loop sets `g.gameOver = true` at line 89. No synchronization.

**Fix:** Use the same mutex that protects direction (or a separate one) to protect `g.gameOver` reads/writes.

---

## HIGH Issues

### 3. rand.Seed Not Called on Game Reset (UNRESOLVED from v1)

**Severity:** HIGH
**File:** `snake/internal/game/game.go`
**Line:** 146 (reset function)

`rand.Seed(time.Now().UnixNano())` is called only in `New()` (line 23). The `reset()` function (line 146) does NOT call `rand.Seed()`. After pressing R to restart, the food will be generated using the same random sequence that was already exhausted, making the restart predictable and the game less replayable.

**Fix:** Move `rand.Seed()` into `reset()`, or call `New()` entirely from `reset()` instead of duplicating its logic.

---

### 4. Goroutine Never Terminates (UNRESOLVED from v1)

**Severity:** MEDIUM
**File:** `snake/internal/game/game.go`
**Lines:** 166-171

The input goroutine runs `for !g.gameOver { g.HandleInput(); time.Sleep(16 * time.Millisecond) }`. When the game ends, the goroutine keeps running until the process exits — it never receives a stop signal. No `sync.WaitGroup`, `context.Context`, or stop channel is used.

**Fix:** Pass a stop channel or use `sync.WaitGroup` to properly signal and wait for goroutine termination on game exit.

---

## MINOR Issues

### 5. Unused Constant EmptySymbol

**Severity:** LOW
**File:** `snake/internal/constants/constants.go`
**Line:** 8

`EmptySymbol` is defined but never used anywhere in rendering. Either remove it or document why it exists.

---

## What Was Fixed from v1

| Issue | Status |
|-------|--------|
| `TestSnakeSelfCollision` added | FIXED |
| `TestScoreIncreaseOnEatingFood` added | FIXED |
| Data race on direction | UNRESOLVED |
| Data race on gameOver | UNRESOLVED |
| rand.Seed on reset | UNRESOLVED |
| Goroutine never terminates | UNRESOLVED |

---

## Test Coverage Note

`go test -race ./...` passes without warnings. However, this does **not** mean the data races are fixed — it means the test suite does not exercise the concurrent code path (goroutine + main loop together). The race condition only manifests at runtime during actual gameplay when `Run()` is called.

---

## Summary

| Issue | Severity | Status |
|-------|----------|--------|
| Data race: direction | CRITICAL | UNRESOLVED |
| Data race: gameOver | CRITICAL | UNRESOLVED |
| rand.Seed on reset | HIGH | UNRESOLVED |
| Goroutine never terminates | MEDIUM | UNRESOLVED |
| Unused EmptySymbol | LOW | UNRESOLVED |

**Recommendation:** REQUEST_CHANGES. The two CRITICAL data races and HIGH reset bug are functional correctness issues that can cause undefined behavior in production. The v2 border rendering feature is correctly implemented and tests pass, but the concurrent code structure needs a mutex or architectural fix before this code is safe to run.
