# Implementation Report - Snake CLI Game v2

## Summary

V2 implements the ASCII border enhancement for the Snake CLI game, replacing the single `#` border character with proper ASCII border characters: `+` for corners, `-` for horizontal edges, and `|` for vertical edges.

## Changes Made

### 1. `snake/internal/constants/constants.go`

**Before:**
```go
BorderSymbol = '#'
```

**After:**
```go
CornerSymbol      = '+'
HorizontalBorder  = '-'
VerticalBorder    = '|'
```

**Rationale:** The plan specified proper ASCII border characters. CornerSymbol, HorizontalBorder, and VerticalBorder are now separate constants to support distinct rendering of border corners and edges.

### 2. `snake/internal/game/game.go` - Render()

**Before:**
```go
for x := 0; x <= constants.GridSize; x++ {
    termbox.SetCell(x, 0, constants.BorderSymbol, ...)
    termbox.SetCell(x, constants.GridSize, constants.BorderSymbol, ...)
}
for y := 0; y <= constants.GridSize; y++ {
    termbox.SetCell(0, y, constants.BorderSymbol, ...)
    termbox.SetCell(constants.GridSize, y, constants.BorderSymbol, ...)
}
```

**After:**
```go
// Draw corners
termbox.SetCell(0, 0, constants.CornerSymbol, ...)
termbox.SetCell(constants.GridSize, 0, constants.CornerSymbol, ...)
termbox.SetCell(0, constants.GridSize, constants.CornerSymbol, ...)
termbox.SetCell(constants.GridSize, constants.GridSize, constants.CornerSymbol, ...)

// Draw horizontal borders (top and bottom)
for x := 1; x < constants.GridSize; x++ {
    termbox.SetCell(x, 0, constants.HorizontalBorder, ...)
    termbox.SetCell(x, constants.GridSize, constants.HorizontalBorder, ...)
}
// Draw vertical borders (left and right)
for y := 1; y < constants.GridSize; y++ {
    termbox.SetCell(0, y, constants.VerticalBorder, ...)
    termbox.SetCell(constants.GridSize, y, constants.VerticalBorder, ...)
}
```

**Rationale:** Separates corner and edge rendering to apply correct symbols at each position. Corners are drawn first at the four corners, then horizontal edges fill between corners on top/bottom, and vertical edges fill between corners on left/right.

### 3. `snake/internal/snake/snake_test.go`

**Added:** `TestSnakeSelfCollision` test

```go
func TestSnakeSelfCollision(t *testing.T) {
    s := New(5, 5)
    s.SetDirection(Point{X: 1, Y: 0})
    s.Grow()
    s.Grow()
    s.body = []Point{{X: 5, Y: 5}, {X: 5, Y: 5}, {X: 4, Y: 5}}
    if !s.CheckSelfCollision() {
        t.Error("Expected self collision when head overlaps with body")
    }
}
```

**Rationale:** The plan specified testing self-collision detection, which was missing from v1.

### 4. `snake/internal/game/game_test.go`

**Added:** `TestScoreIncreaseOnEatingFood` test

```go
func TestScoreIncreaseOnEatingFood(t *testing.T) {
    g := New()
    initialScore := g.score
    head := g.snake.Head()
    g.food = food.New(head.X, head.Y)
    g.Update()
    if g.score != initialScore+10 {
        t.Errorf("Score should increase by 10 when eating food, got %d", g.score)
    }
}
```

**Rationale:** The plan specified testing score calculation when eating food (+10), which the original `TestScoreIncrease` did not cover.

## Verification

### Build
```bash
$ go build ./...
# No output = success
```

### Tests
```bash
$ go test ./... -v
ok  snake/internal/food    0.00s
ok  snake/internal/game    1.341s
ok  snake/internal/snake   1.841s
```

### Lint
```bash
$ go vet ./...
# No output = success
```

## File Summary

| File | Lines Changed | Description |
|------|---------------|-------------|
| `constants.go` | ~4 | Added CornerSymbol, HorizontalBorder, VerticalBorder |
| `game.go` | ~15 | Updated Render() for proper ASCII border |
| `snake_test.go` | ~12 | Added TestSnakeSelfCollision |
| `game_test.go` | ~8 | Added TestScoreIncreaseOnEatingFood |

## Remaining Items

- None - all plan items implemented

## Notes

- The `rand.Seed` deprecation warning (Go 1.20+) is present but does not affect functionality
- This is a CLI game requiring a terminal with ANSI escape code support
- No E2E/Playwright tests needed as this is a terminal-based application
