# 贪吃蛇 CLI 游戏 - 实现计划 v1

## 1. 需求分析与验收标准

### 需求概述
在终端运行的贪吃蛇游戏，支持方向键控制、计分、游戏结束/重新开始。

### 功能验收标准
- [x] 20x20 游戏网格，蛇用 `@` 或 `█` 表示，食物用 `*` 表示
- [x] 方向键控制蛇移动方向（WASD 或方向键）
- [x] 蛇吃食物后长度+1，分数+10
- [x] 撞墙或撞自己时游戏结束，显示 "Game Over" 和最终分数
- [x] 按 R 键重新开始游戏
- [x] 游戏帧率约 10 FPS（100ms 间隔）
- [x] 边界和障碍物检测

## 2. 技术架构

### 技术选型
- **语言**: Go（使用 termbox-go 库进行终端处理）
- **输入处理**: termbox-go 捕获键盘事件
- **渲染**: termbox-go 字符画布，刷新终端显示

### 组件设计

```
snake/
├── cmd/snake/main.go          # 游戏入口
├── internal/
│   ├── constants/constants.go  # 常量定义
│   ├── game/game.go           # Game 类：主循环、状态管理
│   ├── snake/snake.go         # Snake 类：位置、方向、增长逻辑
│   └── food/food.go           # Food 类：食物生成
```

### 数据流
```
键盘输入 → Game.handleKeyEvent() → Snake.direction
                                        ↓
                                  Game.Update() → Snake.Move()/Grow()
                                        ↓
                                  Game.Render() → termbox 绘制
```

## 3. 文件变更清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `cmd/snake/main.go` | 已有 | 入口，初始化游戏 |
| `internal/game/game.go` | 已有 | 主循环、状态管理 |
| `internal/snake/snake.go` | 已有 | 蛇逻辑 |
| `internal/food/food.go` | 已有 | 食物逻辑 |
| `internal/constants/constants.go` | 已有 | 常量定义 |

## 4. 数据模型

### Snake
```go
type Point struct { X, Y int }
type Snake struct {
    body      []Point
    direction Point
}
```

### Food
```go
type Food struct {
    position snake.Point
}
```

### Game State
```go
type Game struct {
    score    int
    gameOver bool
    snake    *snake.Snake
    food     *food.Food
    mu       sync.Mutex
}
```

## 5. 测试策略

### 单元测试
- `snake_test.go`: 测试蛇移动、增长、碰撞检测
- `food_test.go`: 测试食物生成（不与蛇身重叠）
- `game_test.go`: 测试分数计算、游戏结束条件

### 手动测试
- 在终端运行 `go run cmd/snake/main.go` 验证完整流程

## 6. 边缘情况与约束

### 约束
- 仅支持支持 ANSI 转义码的终端（Linux/macOS/Windows Terminal）
- 需要 termbox-go 库

### 边缘情况
- 蛇不能 180 度反向移动
- 食物不能生成在蛇身上
- 游戏结束前需按键才退出
- 窗口大小调整不敏感（固定 20x20）
