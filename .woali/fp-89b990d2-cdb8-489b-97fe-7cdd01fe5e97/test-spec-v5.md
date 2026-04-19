# live-verdict-screenshot-check - Test Specification v5

## 1. Overview

This test specification covers the live-verdict-screenshot-check sub-feature for the snake game, ensuring:
1. Tests are not falsely marked as failed due to infrastructure timeouts
2. E2E screenshots are included in PRs to prove game functionality

## 2. Unit Test Coverage

### 2.1 snake/internal/food Package

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestFoodGenerateNotOnSnake` | Food should not spawn on snake body | PASS - 100 iterations, no food on snake |
| `TestFoodPosition` | Food position getter should return correct coords | PASS - position returned correctly |

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
| `TestSnakeMove` | Snake should move in set direction | PASS - head moves correctly |
| `TestSnakeGrow` | Snake should grow when eating food | PASS - body length increases by 1 |
| `TestSnakeReverseDirection` | Reverse direction should be blocked | PASS - reverse detected, 90° turns allowed |
| `TestSnakeWallCollision` | Wall collision should be detected | PASS - collision detected |
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
| `e2e-screenshot.sh` exists | Script file should exist in `snake/scripts/` | PASS |
| `asciinema.sh` exists | Alternative script should exist | PASS |
| `screenshots/` directory | Screenshots directory should exist | PASS |

### 4.2 E2E Screenshot Execution & Verification

```bash
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
| `.cast` file git-tracked | `git ls-files snake/screenshots/*.cast` shows file | PASS |

### 4.3 Git Tracking Verification

```bash
# Verify .cast files are tracked by git
git ls-files snake/screenshots/*.cast

# Expected output: list of .cast files including recent recording
```

## 5. Anti-Misjudgment Mechanisms

| Mechanism | Purpose | Verification |
|-----------|---------|--------------|
| `-timeout 5m` on tests | Prevents default timeout misjudgment | PASS - explicit timeout set |
| `-race` flag | Detects data races causing flaky failures | PASS - no races found |
| `\|\| true` on timeout | Ensures script exits 0 when asciinema killed | PASS - exit 0 ensured |
| Valid .cast format check | Ensures .cast file is asciinema, not empty typescript | PASS - JSON header validation |
| Git tracking | Ensures E2E screenshots included in PR | PASS - files git-tracked |

## 6. Environment & Setup Notes

### Prerequisites
- Go 1.x installed
- **asciinema** installed (required for proper .cast format)
  - macOS: `brew install asciinema`
  - Linux: `pip install asciinema`
- Terminal with color support

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

# Verify git tracking
git ls-files snake/screenshots/*.cast
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
8. `.cast` file is git-tracked (visible in `git ls-files`)

### Any test FAIL when:
1. Build fails
2. Vet reports issues
3. Any test fails
4. Race detector finds data races
5. `.cast` file is empty or not in asciinema format
6. `.cast` file is not git-tracked

## 8. Test Execution Summary

- **Total Unit Tests**: 12
- **Build Checks**: 2 (build, vet)
- **Race Detection**: 1
- **E2E Scripts**: 2 (e2e-screenshot.sh, asciinema.sh)
- **Git Tracking**: Verified for all .cast files

All tests must PASS for this feature to be considered complete.

## 9. Known Limitations

### termbox-go in PTY
- termbox-go may not initialize properly in all PTY environments
- If termbox.Init() fails, the game will panic with terminal error
- This is a known limitation and may cause E2E tests to fail in some environments

### asciinema Python Bug
- asciinema Python version 2.4.0 has a bug where it crashes with "ValueError: list.remove(x): x not in list" when terminated by timeout
- The .cast file is still written correctly before the crash
- stderr is suppressed in the script to avoid ugly traceback output

### Headless CI
- Cannot run snake game E2E tests in headless CI environments (requires TTY)
- Unit tests can run in any environment