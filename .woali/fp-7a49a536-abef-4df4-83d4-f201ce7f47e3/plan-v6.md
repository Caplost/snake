# 贪吃蛇 CLI 游戏 - 实现计划 v6

## 上下文继承

- **v2**: 原始架构设计 - Go + termbox-go，ASCII 边框，20x20 网格
- **v3**: 代码审查修复 - 数据竞态（mutex/waitgroup）、rand.Seed 移动到 reset()
- **v4**: 移除 input goroutine，简化 Run() 为单线程事件循环；**发现测试覆盖缺口**
- **v5**: 颜色实现完成，12/12 测试通过，竞态检测通过

## v4 代码审查发现的关键问题（遗留）

| 优先级 | 问题 | 当前状态 |
|--------|------|----------|
| HIGH | `reset()` 零测试覆盖 | 未修复 |
| HIGH | `handleKeyEvent()` 零测试覆盖 | 未修复 |
| MEDIUM | `trySetDirection()` 无直接测试 | 未修复 |
| LOW | `Render()` 无 mutex 保护读取 | 未修复 |

**v6 目标**：修复以上 4 个测试覆盖/代码问题，使 game 包测试覆盖率达到 40%+。

---

## 1. 需求分析与验收标准

### 功能（已在 v5 完成，v6 不变更功能）

- [x] 20x20 游戏网格，带 ASCII 边框（角 `+`，边 `-` `|`）
- [x] 蛇用 `█` 表示，**黄色**
- [x] 食物用 `*` 表示，**红色**
- [x] 方向键或 WASD 控制蛇移动方向
- [x] 蛇吃食物后长度+1，分数+10
- [x] 撞墙或撞自己时游戏结束，显示 "Game Over" 和最终分数
- [x] 按 R 键重新开始游戏
- [x] 游戏帧率约 10 FPS（100ms 间隔）

### v6 验收标准（测试驱动）

- [ ] `reset()` 函数有 3+ 个单元测试覆盖
- [ ] `handleKeyEvent()` 函数有 3+ 个单元测试覆盖
- [ ] `trySetDirection()` 有直接测试覆盖
- [ ] `Render()` 读取状态时持有 mutex（或通过其他方式消除数据竞争）
- [ ] `go test -cover ./...` game 包覆盖率 40%+
- [ ] 所有测试通过（`go test -race ./...`）

---

## 2. 技术架构

### 现有组件结构（无变化）

```
snake/
├── cmd/snake/main.go          # 入口
├── internal/
│   ├── constants/constants.go # 常量
│   ├── snake/snake.go          # 蛇逻辑
│   ├── food/food.go            # 食物逻辑
│   └── game/game.go           # 游戏主循环
└── go.mod
```

### `handleKeyEvent` 测试策略

`handleKeyEvent` 依赖 `termbox.Event`，已在 v4 重构为接收参数而非内部调用 `PollEvent()`，因此**可以直接单元测试**。

### `reset()` 测试策略

`reset()` 在 game 内部使用 mutex 保护，可通过构造 game 实例后调用（非导出方法通过反射或包内测试直接访问）。

### `Render()` mutex 问题

`Render()` 在 lines 133-138 读取 `g.snake.Body()` 和 `g.food.Position()` 时未持有 `g.mu`。解决方案：在 `Render()` 开头加读锁 `g.mu.RLock()`，或接受当前单线程事件循环下无实际问题（review 建议修复）。

---

## 3. 文件级变更列表

### `snake/internal/game/game.go`

| 位置 | 变更 | 原因 |
|------|------|------|
| `Render()` (line 113) | 添加 `g.mu.RLock()` / `defer g.mu.RUnlock()` | 修复 LOW 优先级数据竞争风险 |
| `reset()` (line 155) | 无代码变更 | 仅为测试目标 |

### `snake/internal/game/game_test.go`

| 新增测试 | 覆盖目标 |
|----------|----------|
| `TestResetClearsScore` | `reset()` 分数清零 |
| `TestResetClearsGameOver` | `reset()` 游戏结束标志清除 |
| `TestResetReinitializesSnake` | `reset()` 蛇重新初始化到中心 |
| `TestResetGeneratesNewFood` | `reset()` 生成不与蛇重叠的食物 |
| `TestHandleKeyEvent_RestartsGame` | R 键调用 reset() |
| `TestHandleKeyEvent_ExitOnEsc` | Esc 设置 gameOver=true |
| `TestHandleKeyEvent_ExitOnCtrlC` | Ctrl+C 设置 gameOver=true |
| `TestHandleKeyEvent_IgnoresNonKeyEvents` | 非键盘事件被忽略 |
| `TestHandleKeyEvent_WASDAndArrows` | WASD/方向键调用 trySetDirection |
| `TestTrySetDirection_BlocksReverse` | trySetDirection 阻止反向 |

**新增测试数量**：10 个（game 包从 4 个测试增加到 14 个）

---

## 4. 数据模型变更

无数据模型变更。

---

## 5. 测试策略

### 覆盖率目标

| 包 | 当前覆盖 | 目标覆盖 |
|----|----------|----------|
| snake | 87.0% | 87%+ (不变) |
| food | 100% | 100% (不变) |
| game | **15.6%** | **40%+** |

### 新增测试用例详解

#### `TestResetClearsScore`
```go
func TestResetClearsScore(t *testing.T) {
    g := New()
    g.score = 100
    g.reset()
    if g.score != 0 {
        t.Errorf("reset() should clear score, got %d", g.score)
    }
}
```

#### `TestResetClearsGameOver`
```go
func TestResetClearsGameOver(t *testing.T) {
    g := New()
    g.setGameOver(true)
    g.reset()
    if g.gameOver {
        t.Error("reset() should clear gameOver flag")
    }
}
```

#### `TestResetReinitializesSnake`
```go
func TestResetReinitializesSnake(t *testing.T) {
    g := New()
    // Move snake away from center
    g.snake.Move()
    g.reset()
    head := g.snake.Head()
    expectedX := constants.GridSize / 2
    if head.X != expectedX || head.Y != expectedX {
        t.Errorf("reset() should reinit snake at center (%d,%d), got (%d,%d)",
            expectedX, expectedX, head.X, head.Y)
    }
}
```

#### `TestResetGeneratesNewFood`
```go
func TestResetGeneratesNewFood(t *testing.T) {
    g := New()
    g.reset()
    snakeBody := g.snake.Body()
    foodPos := g.food.Position()
    for _, p := range snakeBody {
        if p.X == foodPos.X && p.Y == foodPos.Y {
            t.Error("reset() should generate food not overlapping snake")
        }
    }
}
```

#### `TestHandleKeyEvent_RestartsGame`
```go
func TestHandleKeyEvent_RestartsGame(t *testing.T) {
    g := New()
    g.setGameOver(true)
    g.score = 999
    ev := termbox.Event{Type: termbox.EventKey, Ch: 'r'}
    g.handleKeyEvent(ev)
    if g.gameOver {
        t.Error("R key should restart game (clear gameOver)")
    }
    if g.score != 0 {
        t.Errorf("R key should reset score, got %d", g.score)
    }
}
```

#### `TestHandleKeyEvent_ExitOnEsc`
```go
func TestHandleKeyEvent_ExitOnEsc(t *testing.T) {
    g := New()
    ev := termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEsc}
    g.handleKeyEvent(ev)
    if !g.GameOver() {
        t.Error("Esc should trigger game over")
    }
}
```

#### `TestHandleKeyEvent_ExitOnCtrlC`
```go
func TestHandleKeyEvent_ExitOnCtrlC(t *testing.T) {
    g := New()
    ev := termbox.Event{Type: termbox.EventKey, Key: termbox.KeyCtrlC}
    g.handleKeyEvent(ev)
    if !g.GameOver() {
        t.Error("Ctrl+C should trigger game over")
    }
}
```

#### `TestHandleKeyEvent_IgnoresNonKeyEvents`
```go
func TestHandleKeyEvent_IgnoresNonKeyEvents(t *testing.T) {
    g := New()
    g.score = 123
    ev := termbox.Event{Type: termbox.EventResize}
    g.handleKeyEvent(ev)
    if g.GameOver() || g.score != 123 {
        t.Error("Non-key events should be ignored")
    }
}
```

#### `TestHandleKeyEvent_WASDAndArrows`
```go
func TestHandleKeyEvent_WASDAndArrows(t *testing.T) {
    g := New()
    tests := []struct {
        ev       termbox.Event
        expected snake.Point
    }{
        {termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowUp}, snake.Point{X: 0, Y: -1}},
        {termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowDown}, snake.Point{X: 0, Y: 1}},
        {termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowLeft}, snake.Point{X: -1, Y: 0}},
        {termbox.Event{Type: termbox.EventKey, Key: termbox.KeyArrowRight}, snake.Point{X: 1, Y: 0}},
        {termbox.Event{Type: termbox.EventKey, Ch: 'w'}, snake.Point{X: 0, Y: -1}},
        {termbox.Event{Type: termbox.EventKey, Ch: 's'}, snake.Point{X: 0, Y: 1}},
        {termbox.Event{Type: termbox.EventKey, Ch: 'a'}, snake.Point{X: -1, Y: 0}},
        {termbox.Event{Type: termbox.EventKey, Ch: 'd'}, snake.Point{X: 1, Y: 0}},
    }
    for _, tt := range tests {
        g.handleKeyEvent(tt.ev)
        dir := g.snake.Direction()
        if dir != tt.expected {
            t.Errorf("key event should set direction to %v, got %v", tt.expected, dir)
        }
    }
}
```

#### `TestTrySetDirection_BlocksReverse`
```go
func TestTrySetDirection_BlocksReverse(t *testing.T) {
    g := New()
    // Current direction is (1,0) - moving right
    // Try to go left (reverse) - should be blocked
    g.trySetDirection(snake.Point{X: -1, Y: 0})
    dir := g.snake.Direction()
    if dir.X != 1 || dir.Y != 0 {
        t.Error("trySetDirection should block reverse direction")
    }
}
```

### 验证命令

```bash
cd /Users/wangyinneng/SaaS/Woali/workspace/room-ab68f9ee-f954-4ec9-a486-84e0d9d86bca/snake

# 构建
go build ./...

# 运行新测试
go test ./internal/game/... -v -run "TestReset|TestHandleKeyEvent|TestTrySetDirection"

# 覆盖率检查
go test ./... -cover

# 竞态检测
go test -race ./...
```

---

## 6. 边缘情况、约束与风险

### 边缘情况

| 场景 | 处理 |
|------|------|
| R 键在游戏未结束时按下 | `handleKeyEvent` 检查 `g.GameOver()` 才调用 reset()，安全 |
| 快速连续按键 | 单线程事件循环，每帧处理一个事件，无竞态 |
| reset() 多次调用 | 安全，state 完全重新初始化 |
| 食物生成在蛇身上 | `food.Generate()` 已处理此情况 |

### 约束

- `termbox-go` 需要真实终端，`Render()` 无法在单元测试中验证视觉效果
- `Run()` 函数需要终端输入，单元测试不调用 `Run()`

### 风险

- **无新风险**：v6 仅添加测试和极小代码变更（mutex 添加到 Render），不改变游戏逻辑
