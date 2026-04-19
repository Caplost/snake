# Woali smoke2 - Test Specification v1

## 1. Overview

This test specification covers the Woali smoke2 smoke test for the snake game.

## 2. Unit Test Coverage

### 2.1 snake/internal/food Package

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestFoodGenerateNotOnSnake` | Food should not spawn on snake body | PASS |
| `TestFoodPosition` | Food position getter should return correct coords | PASS |

### 2.2 snake/internal/game Package

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestScoreIncrease` | Initial score should be 0 | PASS |
| `TestScoreIncreaseOnEatingFood` | Score should increase by 10 when eating food | PASS |
| `TestGameOverOnWallCollision` | Snake should trigger game over on wall collision | PASS |
| `TestNewGameInitialization` | New game should initialize all fields | PASS |

### 2.3 snake/internal/snake Package

| Test | Description | Expected Result |
|------|-------------|-----------------|
| `TestSnakeMove` | Snake should move in set direction | PASS |
| `TestSnakeGrow` | Snake should grow when eating food | PASS |
| `TestSnakeReverseDirection` | Reverse direction should be blocked | PASS |
| `TestSnakeWallCollision` | Wall collision should be detected | PASS |
| `TestSnakeOccupies` | Occupies() should check body positions | PASS |
| `TestSnakeSelfCollision` | Self-collision should be detected | PASS |

## 3. Build & Static Analysis

| Check | Command | Expected Result |
|-------|---------|----------------|
| Build | `go build ./...` | PASS |
| Vet | `go vet ./...` | PASS |
| Race Detection | `go test -race ./...` | PASS |

## 4. Integration / E2E Scenarios

### 4.1 E2E Screenshot Script Validation

| Check | Description | Expected Result |
|-------|-------------|-----------------|
| `e2e-screenshot.sh` exists | Script file should exist in `snake/scripts/` | PASS |
| `asciinema.sh` exists | Alternative script should exist | PASS |
| `screenshots/` directory | Screenshots directory should exist | PASS |

### 4.2 E2E Screenshot Verification (Manual Review)

The following should be verified visually from the generated screenshot:

- [ ] 20x20 grid border visible (ASCII borders)
- [ ] Snake displayed in yellow color
- [ ] Food displayed in red color
- [ ] Score displayed in top-left area
- [ ] Screenshot filename matches pattern: `YYYYMMDD-HHMMSS-snake-gameplay.cast`

## 5. Anti-Misjudgment Mechanisms

| Mechanism | Purpose | Verification |
|-----------|---------|--------------|
| `-timeout 5m` on tests | Prevents default timeout misjudgment | PASS |
| `-race` flag | Detects data races that could cause flaky failures | PASS |
| Synchronous test execution | Prevents session stale from interrupting tests | Verified |

## 6. Environment & Setup Notes

### Prerequisites
- Go 1.x installed
- For E2E screenshots: `asciinema` (recommended) or `script` command available
- Terminal with color support (24-bit color ideal for snake game)

### Running Tests

```bash
# Full test suite with race detection
cd snake && go test -timeout 5m -race ./...

# Build verification
cd snake && go build ./...

# Static analysis
cd snake && go vet ./...

# E2E screenshot (requires asciinema or script)
cd snake && ./scripts/e2e-screenshot.sh
```

## 7. Playwright Coverage

This project does not use Playwright. Playwright is a browser automation tool and is not applicable for a terminal-based game.

## 8. Pass/Fail Criteria

### All tests PASS when:
1. `go build ./...` completes without errors
2. `go vet ./...` completes without warnings
3. `go test -timeout 5m -race ./...` all pass (12/12 tests)
4. E2E scripts exist and are executable
5. `screenshots/` directory exists with `.gitkeep`

### Any test FAIL when:
1. Build fails
2. Vet reports issues
3. Any test fails
4. Race detector finds data races

## 9. Test Execution Summary

- **Total Unit Tests**: 12
- **Build Checks**: 2 (build, vet)
- **Race Detection**: 1
- **E2E Scripts**: 2 (e2e-screenshot.sh, asciinema.sh)

All tests must PASS for this feature to be considered complete.
