# snake-game - Test Specification v1

## 1. Overview

**Project:** snake-game (贪吃蛇 CLI 游戏)
**Type:** Go CLI game
**Version:** v1

This test specification covers unit testing, build verification, and E2E screenshot validation for the snake game.

## 2. Unit Test Coverage

### 2.1 snake/internal/food Package

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestFoodGenerateNotOnSnake` | Food should not spawn on snake body after 100 generations | PASS |
| `TestFoodPosition` | Food position getter returns correct coordinates | PASS |

### 2.2 snake/internal/game Package

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestScoreIncrease` | Initial score should be 0 | PASS |
| `TestScoreIncreaseOnEatingFood` | Score increases by 10 when snake eats food | PASS |
| `TestGameOverOnWallCollision` | Snake triggers game over when hitting wall | PASS |
| `TestNewGameInitialization` | New game initializes score=0, gameOver=false, snake!=nil, food!=nil | PASS |

### 2.3 snake/internal/snake Package

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestSnakeMove` | Snake moves in set direction (right: x+1, y unchanged) | PASS |
| `TestSnakeGrow` | Snake body length increases by 1 when growing | PASS |
| `TestSnakeReverseDirection` | 180° reverse direction is blocked; 90° turns allowed | PASS |
| `TestSnakeWallCollision` | Wall collision detected after moving left from x=0 | PASS |
| `TestSnakeOccupies` | Occupies() correctly identifies snake body positions | PASS |
| `TestSnakeSelfCollision` | Self-collision detected when head overlaps body | PASS |

**Total: 12 unit tests**

## 3. Build & Static Analysis

| Check | Command | Expected Result |
|-------|---------|-----------------|
| Build | `go build ./...` | PASS (no errors) |
| Vet | `go vet ./...` | PASS (no warnings) |
| Race Detection | `go test -race ./...` | PASS (no data races) |

## 4. Integration / E2E Scenarios

### 4.1 E2E Screenshot Script Validation

| Check | Description | Expected Result |
|-------|-------------|-----------------|
| `e2e-screenshot.sh` exists | Script exists in `snake/scripts/` | PASS |
| `asciinema.sh` exists | Alternative script exists | PASS |
| `screenshots/` directory | Screenshots directory exists | PASS |
| `screenshots/.gitkeep` | Gitkeep file exists | PASS |

### 4.2 E2E Screenshot Execution

```bash
cd snake && ./scripts/e2e-screenshot.sh
```

| Check | Description | Expected Result |
|-------|-------------|-----------------|
| Script exits 0 | Script completes successfully | PASS |
| `.cast` file created | File exists in `snake/screenshots/` | PASS |
| `.cast` file non-empty | Size > 0 bytes | PASS |
| `.cast` file valid format | First char = `{` (JSON header) | PASS |

### 4.3 Git Tracking Verification

```bash
git ls-files snake/screenshots/*.cast | wc -l
```

Expected: Files are tracked by git.

## 5. Anti-Misjudgment Mechanisms

| Mechanism | Purpose | Verification |
|-----------|---------|--------------|
| `-timeout 5m` | Explicit timeout prevents infrastructure misjudgment | Required for `go test` |
| `-race` flag | Detects data races causing flaky failures | Required for `go test` |
| `|| true` on timeout | Ensures script exits 0 when asciinema killed | In shell scripts |
| JSON header validation | Ensures .cast is valid asciinema format | `head -c 1` |
| Git tracking | Ensures E2E screenshots included in PR | `git ls-files` |

## 6. Environment & Setup Notes

### Prerequisites
- Go 1.x
- **asciinema** (recommended for .cast format)
  - macOS: `brew install asciinema`
  - Linux: `pip install asciinema`
- Terminal with TTY support (termbox-go requirement)

### Running Tests

```bash
# Full test suite
cd snake && go test -timeout 5m -race ./...

# Build verification
cd snake && go build ./...

# Vet
cd snake && go vet ./...

# E2E screenshot
cd snake && ./scripts/e2e-screenshot.sh
```

### Platform Notes
- **macOS**: Fully supported — asciinema rec works, script fallback works with `-F`
- **Linux**: Use `script -f` instead of `script -F`
- **Headless CI**: Unit tests work; E2E requires TTY

## 7. Pass/Fail Criteria

### All PASS when:
1. `go build ./...` completes without errors
2. `go vet ./...` completes without warnings
3. `go test -timeout 5m -race ./...` all pass
4. E2E scripts exist and are executable
5. `.cast` file exists, non-empty, starts with `{`
6. `.cast` file is git-tracked

### Any FAIL when:
1. Build fails
2. Vet reports issues
3. Any test fails or race detector finds issues
4. `.cast` file is empty or not valid asciinema format
5. `.cast` file is not git-tracked

## 8. Playwright Coverage

**Not applicable.** This project uses asciinema for terminal session recording, not Playwright.

## 9. Expected Test Results

| Category | Count | Expected |
|----------|-------|----------|
| Unit Tests | 12 | All PASS |
| Build Checks | 2 | All PASS |
| Race Detection | 1 | PASS (0 races) |
| E2E Scripts | 2 | All PASS |

## 10. Known Limitations

| Limitation | Description | Workaround |
|-----------|-------------|------------|
| termbox-go in PTY | May fail in some PTY environments | Run in real TTY |
| asciinema Python bug | Crash on timeout but file still written | `|| true` ensures exit 0 |
| Headless CI | Cannot run E2E tests (requires TTY) | Unit tests work in any env |