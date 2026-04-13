# 贪吃蛇 CLI 游戏 - 测试规格 v5

## 测试环境

- Go 1.21+
- termbox-go (终端 UI 库)
- macOS/Linux 终端

## 验证命令

```bash
cd /Users/wangyinneng/SaaS/Woali/workspace/room-ab68f9ee-f954-4ec9-a486-84e0d9d86bca/snake

# 构建验证
go build ./...

# 竞态检测
go test -race ./... -v

# 代码检查
go vet ./...
```

## 单元测试覆盖

### snake 包 (snake_test.go)

| 测试 | 描述 | 预期结果 |
|------|------|----------|
| TestSnakeMove | 蛇移动后位置正确 | PASS |
| TestSnakeGrow | 蛇吃食物后长度+1 | PASS |
| TestSnakeReverseDirection | 不能反向移动 | PASS |
| TestSnakeWallCollision | 撞墙检测正确 | PASS |
| TestSnakeOccupies | 位置占用检测正确 | PASS |
| TestSnakeSelfCollision | 自身碰撞检测正确 | PASS |

### food 包 (food_test.go)

| 测试 | 描述 | 预期结果 |
|------|------|----------|
| TestFoodGenerateNotOnSnake | 食物不生成在蛇身上 | PASS |
| TestFoodPosition | 食物位置正确返回 | PASS |

### game 包 (game_test.go)

| 测试 | 描述 | 预期结果 |
|------|------|----------|
| TestNewGameInitialization | 新游戏初始化正确 | PASS |
| TestScoreIncrease | 分数递增正确 | PASS |
| TestScoreIncreaseOnEatingFood | 吃食物后分数+10 | PASS |
| TestGameOverOnWallCollision | 撞墙游戏结束 | PASS |

## 手动测试检查项

由于 termbox-go 需要真实终端，以下检查项需要手动验证：

### 颜色验证
- [ ] 蛇身显示为**黄色** (`ColorYellow`)
- [ ] 食物显示为**红色** (`ColorRed`)
- [ ] 边框显示为**白色** (`ColorWhite`)
- [ ] Game Over 文字显示为**黄色** (`ColorYellow`)

### 游戏逻辑验证
- [ ] 方向键/WASD 控制正常
- [ ] 吃食物加分 (+10)
- [ ] 撞墙/撞自己游戏结束
- [ ] R 键重开正常

### 手动测试步骤

1. 构建并运行游戏：
   ```bash
   cd /Users/wangyinneng/SaaS/Woali/workspace/room-ab68f9ee-f954-4ec9-a486-84e0d9d86bca/snake
   go build -o snake ./cmd/snake
   ./snake
   ```

2. 验证初始状态：
   - 20x20 网格，带 ASCII 边框
   - 蛇（黄色 `█`）在中央
   - 食物（红色 `*`）随机位置

3. 测试控制：
   - 按方向键或 WASD 移动蛇
   - 蛇应朝指定方向移动

4. 测试吃食物：
   - 引导蛇吃食物 `*`
   - 分数应增加 10
   - 蛇身长度应增加 1

5. 测试游戏结束：
   - 撞墙或撞自己
   - 应显示 "Game Over! Score: X"
   - 按 R 键应重新开始

## E2E 测试说明

此项目使用 termbox-go 进行终端渲染，不适合自动化 E2E 测试。所有核心逻辑通过单元测试覆盖。

## 已知限制

- termbox-go 需要真实终端，无法在无头环境测试渲染
- 手动测试需要交互式终端会话
