# 贪吃蛇 CLI 游戏 - 实现计划

## 1. 需求分析与验收标准

### 需求概述
在终端运行的贪吃蛇游戏，支持方向键控制、计分、游戏结束/重新开始。

### 功能验收标准
- [ ] 20x20 游戏网格，蛇用 `█` 表示，食物用 `*` 表示
- [ ] 方向键或 WASD 控制蛇移动方向
- [ ] 蛇吃食物后长度+1，分数+10
- [ ] 撞墙或撞自己时游戏结束，显示 "Game Over" 和最终分数
- [ ] 按 R 键重新开始游戏
- [ ] 游戏帧率约 10 FPS（100ms 间隔）
- [ ] 边界和障碍物检测（自身碰撞）

## 2. 技术架构

### 技术选型
- **语言**: Go 1.21+
- **终端UI**: 标准库 `termbox-go`（跨平台，支持键盘事件和ANSI转义码）
- **编译**: `go build -o snake ./cmd/snake`

### 组件设计

```
snake/
├── cmd/
│   └── snake/
│       └── main.go      # 入口，初始化 termbox，主循环
├── internal/
│   ├── game/
│   │   └── game.go      # Game 结构体：状态管理、渲染
│   ├── snake/
│   │   └── snake.go     # Snake 结构体：移动、增长、碰撞检测
│   ├── food/
│   │   └── food.go      # Food 结构体：随机位置生成
│   └── constants/
│       └── constants.go # 常量：网格大小、速度、符号
└── internal/game/*_test.go  # 单元测试
```

### 数据流
```
键盘输入 → Game.handleInput() → Snake.direction
                                        ↓
                              Game.update() → Snake.move()
                                        ↓
                              Game.render() → termbox 绘制
```

## 3. 文件变更清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `snake/cmd/snake/main.go` | 新建 | 入口，termbox 初始化，main() |
| `snake/internal/constants/constants.go` | 新建 | 常量：GridSize(20)、Tick(100ms)、Symbols |
| `snake/internal/snake/snake.go` | 新建 | Snake 结构体：Move/Grow/CheckCollision |
| `snake/internal/food/food.go` | 新建 | Food 结构体：Generate |
| `snake/internal/game/game.go` | 新建 | Game 结构体：主循环、状态、渲染 |
| `snake/internal/game/game_test.go` | 新建 | 单元测试 |
| `snake/internal/snake/snake_test.go` | 新建 | 单元测试 |
| `snake/internal/food/food_test.go` | 新建 | 单元测试 |
| `snake/go.mod` | 新建 | go.mod 文件 |

## 4. 数据模型

### Snake
```go
type Point struct { X, Y int }

type Snake struct {
    body      []Point      // 蛇身坐标列表，头在 body[0]
    direction Point        // 方向向量，如 (1,0) 表示向右
}
```

### Food
```go
type Food struct {
    position Point         // 食物坐标
}
```

### Game 状态
```go
type Game struct {
    score    int
    gameOver bool
    snake    Snake
    food     Food
}
```

## 5. 测试策略

### 单元测试
- `snake_test.go`: 测试蛇移动、增长、180度反向拒绝、自身碰撞
- `food_test.go`: 测试食物生成位置不与蛇身重叠
- `game_test.go`: 测试分数计算(吃食物+10)、游戏结束条件

### 运行测试
```bash
go test ./internal/... -v
```

## 6. 边缘情况与约束

### 约束
- 仅支持支持 ANSI 转义码的终端（Linux/macOS/Windows Terminal）
- Go 1.21+

### 边缘情况处理
| 情况 | 处理方式 |
|------|----------|
| 蛇不能 180 度反向移动 | 检测新方向与当前方向，拒绝对向移动 |
| 食物不能生成在蛇身上 | 重新随机生成直到不在蛇身上 |
| 游戏结束需按键才退出 | 显示 "Press R to restart" 等待输入 |
| 窗口 resize | 固定 20x20 网格，不响应窗口大小 |

## 7. 交互流程

```
启动 → 初始化 termbox → 显示游戏网格和初始蛇 → 游戏循环:
  → 捕获键盘输入
  → 更新蛇位置
  → 检测碰撞(墙/自身)
  → 检测吃到食物(增长/加分)
  → 刷新显示
  → 循环直到游戏结束
游戏结束 → 显示 Game Over 和分数 → 等待 R 键 → 重新开始
```

## 8. Go vs Python 权衡

| 因素 | Go (本计划) | Python (原计划) |
|------|-------------|-----------------|
| 编译/运行 | 需编译，交叉编译简单 | 解释执行，无需编译 |
| 依赖 | termbox-go 一个依赖 | 标准库 curses，无需依赖 |
| 性能 | 更快，CPU 密集无压力 | 足够，CLI 游戏不敏感 |
| 分发 | 单一二进制文件 | 需 Python 环境 |
| **推荐** | **更适合 CLI 工具** | 适合原型快速验证 |
