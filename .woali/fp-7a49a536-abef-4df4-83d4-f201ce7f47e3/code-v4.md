# 贪吃蛇 CLI 游戏 - v4 实现报告

## 概述

v4 版本主要解决了 input goroutine 在游戏重开时不会重启的 CRITICAL 问题，同时将轮询输入改为事件驱动输入，提高了效率。

## 变更内容

### 1. 移除 `wg` 字段 (game.go:22)

**变更前:**
```go
type Game struct {
    score    int
    gameOver bool
    snake    *snake.Snake
    food     *food.Food
    mu       sync.Mutex
    wg       sync.WaitGroup  // 已移除
}
```

**变更后:** 移除了 `wg sync.WaitGroup` 字段，因为不再需要等待 input goroutine。

### 2. `HandleInput()` 改为 `handleKeyEvent()` (game.go:50-81)

**变更前:** `HandleInput()` 自己调用 `termbox.PollEvent()` 并处理按键

**变更后:** `handleKeyEvent(ev termbox.Event)` 接收已获取的事件作为参数，主循环统一管理事件获取

**关键改进:**
- 移除了内部的 `time.Sleep(16 * time.Millisecond)` 轮询
- 移除了对 `g.GameOver()` 的循环检查

### 3. 简化 `Run()` 主循环 (game.go:165-203)

**变更前:**
```go
// 启动独立 goroutine 处理输入
g.wg.Add(1)
go func() {
    defer g.wg.Done()
    for !g.GameOver() {
        g.HandleInput()  // 内部调用 PollEvent + 16ms sleep
        time.Sleep(16 * time.Millisecond)
    }
}()

for {
    g.Render()
    if g.GameOver() {
        // game over 时等待用户按键
        for { ev := termbox.PollEvent() ... }
    }
    <-ticker.C
    g.Update()
}
```

**变更后:**
```go
for {
    g.Render()

    // 事件驱动输入（阻塞）
    ev := termbox.PollEvent()
    g.handleKeyEvent(ev)

    // 游戏结束时等待用户重启或退出
    if g.GameOver() {
        for {
            ev := termbox.PollEvent()
            if ev.Type == termbox.EventKey {
                if ev.Ch == 'r' || ev.Ch == 'R' {
                    g.reset()
                    break
                }
                if ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyEsc {
                    return
                }
            }
        }
    }

    <-ticker.C
    g.Update()
}
```

## 数据流对比

**v3 数据流:**
```
独立 goroutine: PollEvent() ──(内存)──> 主循环
                                          │
                              ┌───────────┼───────────┐
                              ▼           ▼           ▼
                        HandleInput  Update      Render
```

**v4 数据流:**
```
主循环: PollEvent() ──> handleKeyEvent()
        │
        ├──> trySetDirection ──> snake.direction
        │
        └──> Update ──> 碰撞检测、蛇移动、吃食物
```

## 修复的问题

| 问题 | 严重性 | 修复状态 |
|------|--------|----------|
| input goroutine 在游戏结束时退出，重开时不重启 | CRITICAL | ✅ 已修复 |
| 每 16ms 轮询 termbox，效率低 | HIGH | ✅ 已修复 |

## 验证结果

```bash
$ go build ./...   # 通过
$ go vet ./...     # 通过
$ go test -race ./... -v  # 通过 (无竞态)
```

所有测试通过，无竞态条件检测到。

## 保留的并发保护

虽然移除了 input goroutine，但以下并发保护仍然必要:
- `GameOver()` - mutex 保护 `gameOver` 字段
- `setGameOver()` - mutex 保护 `gameOver` 字段
- `trySetDirection()` - mutex 保护 `snake.direction`

这些保护在 `Update()` 和 `handleKeyEvent()` 并发访问时仍需要。
