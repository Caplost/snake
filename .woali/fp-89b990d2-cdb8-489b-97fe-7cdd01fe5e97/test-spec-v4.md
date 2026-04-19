# live-verdict-screenshot-check - Test Specification v4

## 1. Overview

This test specification covers the live-verdict-screenshot-check sub-feature for the snake game, ensuring tests don't get falsely marked as failed due to infrastructure timeouts, and providing E2E screenshot proof of game functionality.

## 2. Unit Test Coverage

### 2.1 snake/internal/food Package

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestFoodGenerateNotOnSnake` | Food should not spawn on snake body | PASS - 100 iterations, no food on snake |
| `TestFoodPosition` | Food position getter should return correct coords | PASS - position (3,4) returned correctly |

### 2.2 snake/internal/game Package

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestScoreIncrease` | Initial score should be 0 | PASS - score equals 0 |
| `TestScoreIncreaseOnEatingFood` | Score should increase by 10 when eating food | PASS - score increases correctly |
| `TestGameOverOnWallCollision` | Snake should trigger game over on wall collision | PASS - collision detected |
| `TestNewGameInitialization` | New game should initialize all fields | PASS - score=0, gameOver=false, snake/food initialized |

### 2.3 snake/internal/snake Package

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestSnakeMove` | Snake should move in set direction | PASS - head moves right correctly |
| `TestSnakeGrow` | Snake should grow when eating food | PASS - body length increases by 1 |
| `TestSnakeReverseDirection` | Reverse direction should be blocked | PASS - reverse detected, 90° turns allowed |
| `TestSnakeWallCollision` | Wall collision should be detected | PASS - collision at x=0 moving left |
| `TestSnakeOccupies` | Occupies() should check body positions | PASS - correct occupancy detection |
| `TestSnakeSelfCollision` | Self-collision should be detected | PASS - overlapping body detected |

## 3. Build & Static Analysis

| Check | Command | Expected Result |
|-------|---------|-----------------|
| Build | `go build ./...` | PASS - no errors |
| Vet | `go vet ./...` | PASS - no warnings |
| Race Detection | `go test -race ./...` | PASS - no data races |

## 4. Integration / E2E Scenarios

### 4.1 E2E Screenshot Script Validation

| Check | Description | Expected Result |
|-------|-------------|-----------------|
| `e2e-screenshot.sh` exists | Script file should exist in `snake/scripts/` | PASS - file exists and is executable |
| `asciinema.sh` exists | Alternative script should exist | PASS - file exists and is executable |
| `screenshots/` directory | Screenshots directory should exist | PASS - directory created with .gitkeep |

### 4.2 E2E Screenshot Execution & Verification

```bash
# Ensure asciinema is available (install if needed)
# macOS: brew install asciinema
# Linux: pip install asciinema

# Run E2E screenshot script
cd snake && ./scripts/e2e-screenshot.sh

# Verify .cast file is non-empty AND has valid asciinema format
test -s snake/screenshots/*.cast
# AND first character should be '{' (JSON header)
head -c 1 snake/screenshots/*.cast | grep -q '{'
```

| Check | Description | Expected Result |
|-------|-------------|-----------------|
| E2E script exits 0 | Script completes successfully | PASS |
| `.cast` file created | File exists in `snake/screenshots/` | PASS |
| `.cast` file non-empty | `test -s` returns true (size > 0) | PASS |
| `.cast` file is asciinema format | First character is `{` (JSON header) | PASS |
| `.cast` file has JSON header | File starts with `{"version": 2,` | PASS |

**v4 Fix**: The `e2e-screenshot.sh` now uses `asciinema rec` to produce proper `.cast` files (asciinema v2 format) instead of `script` which produces typescript format. The script also validates that the output file starts with a JSON header (`{`) to ensure it's valid asciinema format.

**Prerequisite**: `asciinema` must be installed. Install via:
- macOS: `brew install asciinema`
- Linux: `pip install asciinema`

If asciinema is not installed, the script falls back to `script` command (typescript format) but will warn that the output is not asciinema format.

## 5. Anti-Misjudgment Mechanisms

| Mechanism | Purpose | Verification |
|-----------|---------|--------------|
| `-timeout 5m` on tests | Prevents default timeout misjudgment | PASS - explicit timeout set |
| `-race` flag | Detects data races that could cause flaky failures | PASS - no races found |
| Event-driven input | No goroutine lifecycle issues on restart | Verified via code review |
| Valid .cast format check | Ensures .cast file is asciinema, not empty typescript | PASS - JSON header validation |

## 6. Environment & Setup Notes

### Prerequisites
- Go 1.x installed
- **asciinema** installed (required for proper .cast format)
  - macOS: `brew install asciinema`
  - Linux: `pip install asciinema`
- Terminal with color support (24-bit color ideal for snake game)

### Running Tests

```bash
# Full test suite with race detection
cd snake && go test -timeout 5m -race ./...

# Build verification
cd snake && go build ./...

# Static analysis
cd snake && go vet ./...

# E2E screenshot (requires asciinema)
cd snake && ./scripts/e2e-screenshot.sh
```

## 7. Pass/Fail Criteria

### All tests PASS when:
1. `go build ./...` completes without errors
2. `go vet ./...` completes without warnings
3. `go test -timeout 5m -race ./...` all pass (12/12 tests)
4. E2E scripts exist and are executable
5. `screenshots/` directory exists with `.gitkeep`
6. `.cast` file exists and is non-empty
7. `.cast` file starts with `{` (valid asciinema JSON header)

### Any test FAIL when:
1. Build fails
2. Vet reports issues
3. Any test fails
4. Race detector finds data races
5. `.cast` file is empty or not in asciinema format

## 8. Test Execution Summary

- **Total Unit Tests**: 12
- **Build Checks**: 2 (build, vet)
- **Race Detection**: 1
- **E2E Scripts**: 2 (e2e-screenshot.sh, asciinema.sh)

All tests must PASS for this feature to be considered complete.

## 9. Event-Driven Architecture Validation

### Input Handling
- Uses `termbox.PollEvent()` for blocking event-driven input
- No background goroutine for input handling
- Main loop handles all events synchronously

### Game Restart Flow
1. Game over detected in `Update()`
2. `setGameOver(true)` called
3. Main loop enters restart-wait loop (lines 186-197)
4. User presses R → `reset()` called → game restarts
5. User presses Esc/Ctrl+C → function returns, game exits

This flow has no goroutine lifecycle issues because no input goroutine exists.

## 10. Known Limitations

### termbox-go in PTY
- termbox-go may not initialize properly in all PTY environments
- If termbox.Init() fails, the game will panic with "panic: error attempting to initialize terminal"
- This is a known limitation and may cause E2E tests to fail in some environments

### asciinema Python Bug
- asciinema Python version 2.4.0 has a bug where it crashes with "ValueError: list.remove(x): x not in list" when terminated by timeout
- The .cast file is still written correctly before the crash
- stderr is suppressed in the script to avoid ugly traceback output
