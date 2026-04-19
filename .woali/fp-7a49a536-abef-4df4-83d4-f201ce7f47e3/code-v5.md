# 贪吃蛇 CLI 游戏 - 实现报告 v5

## 概述

v5 实现完成，基于 v4 的简化架构（移除 input goroutine，单线程事件循环），确保所有功能验收标准通过。

## 代码变更

无代码变更（v5 与 v4 代码相同）。

## 文件结构

```
snake/
├── go.mod
├── cmd/
│   └── snake/
│       └── main.go              # 入口
└── internal/
    ├── constants/
    │   └── constants.go         # 常量定义
    ├── snake/
    │   ├── snake.go             # 蛇逻辑
    │   └── snake_test.go        # 测试
    ├── food/
    │   ├── food.go              # 食物逻辑
    │   └── food_test.go         # 测试
    └── game/
        ├── game.go              # 游戏主循环
        └── game_test.go         # 测试
```

## 验证结果

| 检查项 | 状态 |
|--------|------|
| `go build ./...` | ✓ PASS |
| `go vet ./...` | ✓ PASS |
| `go test -race ./... -v` | ✓ PASS (11 tests) |

## 功能验收

- [x] 20x20 游戏网格，带 ASCII 边框（角 `+`，边 `-` `|`）
- [x] 蛇用 `█` 表示，**黄色** (`termbox.ColorYellow`)
- [x] 食物用 `*` 表示，**红色** (`termbox.ColorRed`)
- [x] 方向键或 WASD 控制蛇移动方向
- [x] 蛇吃食物后长度+1，分数+10
- [x] 撞墙或撞自己时游戏结束，显示 "Game Over" 和最终分数
- [x] 按 R 键重新开始游戏
- [x] 游戏帧率约 10 FPS（100ms 间隔）
- [x] 边界和障碍物检测（自身碰撞）
- [x] 颜色实现：
  - 蛇身：黄色 (`ColorYellow`) - `game.go:134`
  - 食物：红色 (`ColorRed`) - `game.go:138`
  - 边框：白色 (`ColorWhite`) - `game.go:117-130`
  - Game Over 文字：黄色 (`ColorYellow`) - `game.go:143`
  - 分数文字：白色 (`ColorWhite`) - `game.go:149`

## 依赖

| 依赖 | 版本 | 用途 |
|------|------|------|
| Go | 1.21+ | 语言 |
| termbox-go | latest | 终端 UI + 颜色支持 |
