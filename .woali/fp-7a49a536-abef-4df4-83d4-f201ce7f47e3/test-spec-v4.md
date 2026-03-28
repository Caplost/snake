# 贪吃蛇 CLI 游戏 - v4 测试规格

## 概述

v4 版本简化了事件循环，移除 input goroutine，改为事件驱动输入。本测试规格涵盖单元测试、集成场景和手动测试清单。

## 测试环境

- **平台**: macOS (darwin)
- **语言**: Go 1.21+
- **终端**: 支持 termbox-go 的终端（xterm-256color 或类似）
- **依赖**:
  - `github.com/nsf/termbox-go`

## 1. 单元测试

### 1.1 snake 包 (`snake/internal/snake/snake_test.go`)

| 测试用例 | 描述 | 预期结果 |
|----------|------|----------|
| `TestSnakeMove` | 测试蛇向当前方向移动 | 蛇头坐标正确更新 |
| `TestSnakeGrow` | 测试蛇吃食物后增长 | 蛇身长度+1 |
| `TestSnakeReverseDirection` | 测试180度反向检测 | 90度转弯允许，180度拒绝 |
| `TestSnakeWallCollision` | 测试撞墙检测 | 超出边界返回 true |
| `TestSnakeOccupies` | 测试位置占用检测 | 蛇身位置返回 true |
| `TestSnakeSelfCollision` | 测试自撞检测 | 蛇头与身体重叠返回 true |

### 1.2 food 包 (`snake/internal/food/food_test.go`)

| 测试用例 | 描述 | 预期结果 |
|----------|------|----------|
| `TestFoodGenerateNotOnSnake` | 测试食物不在蛇身上生成 | 100次生成均不与蛇身重叠 |
| `TestFoodPosition` | 测试食物位置获取 | 返回正确的 X,Y 坐标 |

### 1.3 game 包 (`snake/internal/game/game_test.go`)

| 测试用例 | 描述 | 预期结果 |
|----------|------|----------|
| `TestScoreIncrease` | 测试初始分数为0 | score == 0 |
| `TestScoreIncreaseOnEatingFood` | 测试吃食物后分数+10 | score == initialScore + 10 |
| `TestGameOverOnWallCollision` | 测试撞墙触发游戏结束 | CheckWallCollision 返回 true |
| `TestNewGameInitialization` | 测试新游戏初始化状态 | score=0, gameOver=false, snake/food 已初始化 |

## 2. 集成测试

### 2.1 游戏流程集成测试

| 场景 | 测试方法 | 预期结果 |
|------|----------|----------|
| 游戏初始化 | 创建 Game 并调用 reset() | snake 在中心，food 不与 snake 重叠 |
| 方向控制 | 调用 handleKeyEvent() 模拟方向键 | snake.direction 更新为正确方向 |
| 180度反向拒绝 | 先设置方向向上，再发送向下事件 | direction 保持向上 |
| 吃食物加分 | 将 food 放在蛇头位置，调用 Update() | score + 10，蛇身增长 |
| 撞墙游戏结束 | 将蛇头移动到边界外，调用 Update() | gameOver == true |
| 自撞游戏结束 | 创建重叠的蛇身，调用 Update() | gameOver == true |
| 重置游戏 | 游戏结束后调用 reset() | score=0, gameOver=false, snake 重新初始化 |

### 2.2 并发安全测试

```bash
go test -race ./... -v
```

**检测内容:**
- `GameOver()` 与 `setGameOver()` 的 mutex 保护
- `trySetDirection()` 对 `snake.direction` 的 mutex 保护
- `Update()` 与 `handleKeyEvent()` 的并发访问

## 3. 手动测试清单

### 3.1 基本功能测试

| 功能 | 操作 | 预期结果 |
|------|------|----------|
| 方向键控制 | 按上/下/左/右 | 蛇朝相应方向移动 |
| WASD 控制 | 按 w/a/s/d | 蛇朝相应方向移动 |
| 吃食物 | 操控蛇吃 * | 分数+10，蛇身增长 |
| 撞墙结束 | 操控蛇撞墙 | 显示 "Game Over! Score: X" |
| 撞自己结束 | 操控蛇撞自己 | 显示 "Game Over! Score: X" |
| 重新开始 | 游戏结束时按 R | 游戏重置，分数清零 |
| 退出游戏 | 按 Esc 或 Ctrl+C | 游戏退出 |

### 3.2 边界情况测试

| 情况 | 操作 | 预期结果 |
|------|------|----------|
| 180度反向 | 向上时按向下 | 蛇不反向，继续向上 |
| 快速按键 | 快速连按方向键 | 蛇平滑转向，无抖动 |
| 重开时输入 | 游戏结束后先按其他键再按 R | R 正确重启游戏 |
| 多次重开 | 游戏结束后按 R 多次 | 每次都正确重置 |

### 3.3 帧率测试

| 指标 | 验证方法 | 预期结果 |
|------|----------|----------|
| 帧率 | 观察蛇移动速度 | 约 10 FPS (100ms/帧) |
| 输入响应 | 按方向键后蛇立即响应 | 下一帧开始朝新方向 |

## 4. Playwright E2E 测试

**注意**: termbox-go 是 TUI 库，不适用于浏览器 Playwright 测试。对于 CLI 游戏的 E2E 测试，建议使用以下方法之一：

### 4.1 推荐：使用 `go test` + mock termbox

```go
// game_test.go
func TestGameRestart(t *testing.T) {
    // 创建 Game 实例
    g := New()

    // 模拟撞墙
    g.snake.SetDirection(Point{X: -1, Y: 0})
    for i := 0; i < 15; i++ {
        g.Update()
    }

    // 验证游戏结束
    if !g.GameOver() {
        t.Error("Expected game over after wall collision")
    }

    // 模拟按 R 重开
    g.reset()

    // 验证重置
    if g.GameOver() || g.Score() != 0 {
        t.Error("Expected game to be reset")
    }
}
```

### 4.2 手动终端测试脚本

```bash
# 启动游戏
./snake &

# 发送方向键测试
echo -n "w" | ./snake
sleep 0.2

# 验证进程退出
pkill -f "snake"
```

## 5. 测试覆盖率目标

| 包 | 覆盖率目标 | 当前状态 |
|----|------------|----------|
| snake | 90%+ | ✅ 已有单元测试 |
| food | 90%+ | ✅ 已有单元测试 |
| game | 80%+ | ✅ 已有单元测试 |

## 6. 运行测试

```bash
# 所有测试
cd snake && go test ./...

# 带竞态检测
go test -race ./... -v

# 覆盖率报告
go test -cover ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# 特定包测试
go test ./internal/game/... -v
```

## 7. 测试失败处理

| 失败类型 | 可能原因 | 解决方法 |
|----------|----------|----------|
| `TestFoodGenerateNotOnSnake` | 随机数生成问题 | 检查 rand 是否正确初始化 |
| `TestScoreIncreaseOnEatingFood` | 食物位置不对 | 检查 food.Generate() |
| 竞态检测失败 | mutex 保护缺失 | 检查 Game 结构体的 mutex 使用 |
| `go build` 失败 | 缺少依赖 | 运行 `go mod tidy` |
