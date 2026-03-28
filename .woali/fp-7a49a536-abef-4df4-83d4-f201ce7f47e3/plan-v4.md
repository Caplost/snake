# 贪吃蛇 CLI 游戏 - 实现计划 v4

## 上下文继承

- **v2**: 原始架构设计 - Go + termbox-go，ASCII 边框，20x20 网格
- **v3**: 代码审查修复 - 数据竞态（mutex/waitgroup）、rand.Seed 移动到 reset()

## 1. 需求分析与验收标准

### 功能验收标准
- [x] 20x20 游戏网格，带 ASCII 边框（角 `+`，边 `-` `|`）
- [x] 蛇用 `█` 表示，食物用 `*` 表示
- [x] 方向键或 WASD 控制蛇移动方向
- [x] 蛇吃食物后长度+1，分数+10
- [x] 撞墙或撞自己时游戏结束，显示 "Game Over" 和最终分数
- [x] 按 R 键重新开始游戏
- [x] 游戏帧率约 10 FPS（100ms 间隔）
- [x] 边界和障碍物检测（自身碰撞）

## 2. 当前代码状态

### 已实现
- `snake/cmd/snake/main.go` - 入口
- `snake/internal/constants/constants.go` - 常量定义
- `snake/internal/snake/snake.go` - 蛇逻辑
- `snake/internal/food/food.go` - 食物逻辑
- `snake/internal/game/game.go` - 游戏主循环、mutex/waitgroup
- 测试文件已创建

### 待修复问题（v4）

#### CRITICAL - Input Goroutine 不重启

**问题**: `Run()` 中启动的 input goroutine（game.go:181-188）在 `reset()` 后不会重启。第一次游戏结束后，再次按 R 重开时，goroutine 已退出。

**代码位置**: `game.go:184-188`

```go
// 问题：goroutine 看到 gameOver=true 后退出，reset() 设置 gameOver=false
// 但 goroutine 已经终止，不会重新启动
for !g.GameOver() {  // gameOver=true 时条件为 false，循环退出
    g.HandleInput()
    time.Sleep(16 * time.Millisecond)
}
```

**修复方案**:
1. 将 input goroutine 的生命周期与 `Run()` 循环绑定，而不是与 gameOver 绑定
2. 使用 channel 传递输入事件，goroutine 只负责读取输入并发送
3. 主循环统一处理所有事件

#### HIGH - Polling 输入效率低

**问题**: goroutine 每 16ms 轮询一次 termbox，应改用 event-driven

**修复方案**: 使用 termbox.PollEvent() 阻塞等待输入，消除不必要的 CPU 轮询

## 3. 架构改进

### 改进后数据流
```
termbox.PollEvent() ──(channel)──> 主循环处理
                                    │
                    ┌───────────────┼───────────────┐
                    ▼               ▼               ▼
              handleInput()    Update()        Render()
                    │               │
                    ▼               ▼
              snake.direction   碰撞检测
```

### 关键变更
1. **移除 input goroutine**：主循环直接调用 `termbox.PollEvent()`
2. **简化并发模型**：消除不必要的 goroutine，gameOver 和 direction 仍需 mutex 保护
3. **消除 16ms sleep**：termbox.PollEvent() 本身阻塞，无需轮询

## 4. 文件变更清单

| 文件 | 操作 | 变更说明 |
|------|------|----------|
| `snake/internal/game/game.go` | 修改 | 移除 input goroutine，简化 Run() 循环 |

### game.go 变更详情

```go
// 移除 wg 字段（不再需要）
type Game struct {
    score    int
    gameOver bool
    snake    *snake.Snake
    food     *food.Food
    mu       sync.Mutex
}

// Run() 简化为主循环
func (g *Game) Run() {
    if err := termbox.Init(); err != nil {
        panic(err)
    }
    defer termbox.Close()
    termbox.SetInputMode(termbox.InputEsc)
    g.reset()

    ticker := time.NewTicker(time.Duration(constants.TickMs) * time.Millisecond)
    defer ticker.Stop()

    for {
        g.Render()

        // 事件驱动输入（阻塞）
        ev := termbox.PollEvent()
        if ev.Type == termbox.EventKey {
            g.handleKeyEvent(ev)
        }

        // 游戏逻辑更新
        g.Update()

        <-ticker.C  // 控制帧率
    }
}
```

## 5. 验证清单

```bash
# 构建
go build ./...

# 竞态检测
go test -race ./... -v

# 代码检查
go vet ./...

# 手动测试
./snake
# - 验证方向键/WASD 控制
# - 验证吃食物加分
# - 验证撞墙/撞自己游戏结束
# - 验证 R 键重开
```

## 6. 边缘情况

| 情况 | 处理 |
|------|------|
| 180度反向 | `CheckReverseDirection()` 拒绝 |
| 食物在蛇身 | `Generate()` 重新随机直到不重叠 |
| Ctrl+C/Esc | `setGameOver(true)`，主循环检测到后退出 |
| 快速连续按键 | mutex 保护 direction，串行化处理 |

## 7. 风险

| 风险 | 影响 | 缓解 |
|------|------|------|
| termbox-go 维护不活跃 | 长期兼容性问题 | 可替换为 tcell（API 相似） |
| 移除 goroutine 后输入响应 | 可能略有延迟 | 帧率 100ms，用户感知不明显 |

## 8. Trade-off 说明

**原方案（v3）的 goroutine 设计**:
- 优点：输入处理与游戏逻辑分离
- 缺点：goroutine 生命周期管理复杂，重开时需重启 goroutine

**本方案（v4）的单线程事件循环**:
- 优点：简单可靠，无 goroutine 生命周期问题
- 缺点：输入处理阻塞（但 termbox.PollEvent 本身已阻塞，体感无差异）
