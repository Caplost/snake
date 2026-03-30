# 贪吃蛇 CLI 游戏 - 实现计划 v5

## 上下文继承

- **v2**: 原始架构设计 - Go + termbox-go，ASCII 边框，20x20 网格
- **v3**: 代码审查修复 - 数据竞态（mutex/waitgroup）、rand.Seed 移动到 reset()
- **v4**: 移除 input goroutine，简化 Run() 为单线程事件循环

## 1. 需求分析与验收标准

### 功能验收标准
- [x] 20x20 游戏网格，带 ASCII 边框（角 `+`，边 `-` `|`）
- [x] 蛇用 `█` 表示，**绿色**
- [x] 食物用 `*` 表示，**红色**
- [x] 方向键或 WASD 控制蛇移动方向
- [x] 蛇吃食物后长度+1，分数+10
- [x] 撞墙或撞自己时游戏结束，显示 "Game Over" 和最终分数
- [x] 按 R 键重新开始游戏
- [x] 游戏帧率约 10 FPS（100ms 间隔）
- [x] 边界和障碍物检测（自身碰撞）

### 当前颜色实现

| 元素 | 颜色 | 已在 Render() 中实现 |
|------|------|---------------------|
| 蛇身 | 绿色 (ColorGreen) | ✓ |
| 食物 | 红色 (ColorRed) | ✓ |
| 边框 | 白色 (ColorWhite) | ✓ |
| Game Over 文字 | 黄色 (ColorYellow) | ✓ |
| 分数文字 | 白色 (ColorWhite) | ✓ |

## 2. 当前代码状态

### 已实现文件
- `snake/cmd/snake/main.go` - 入口
- `snake/internal/constants/constants.go` - 常量定义
- `snake/internal/snake/snake.go` - 蛇逻辑
- `snake/internal/snake/snake_test.go` - 测试
- `snake/internal/food/food.go` - 食物逻辑
- `snake/internal/food/food_test.go` - 测试
- `snake/internal/game/game.go` - 游戏主循环
- `snake/internal/game/game_test.go` - 测试
- `snake/go.mod` - 模块定义

## 3. 验证清单

```bash
cd /Users/wangyinneng/SaaS/Woali/workspace/room-ab68f9ee-f954-4ec9-a486-84e0d9d86bca/snake

# 构建
go build ./...

# 竞态检测
go test -race ./... -v

# 代码检查
go vet ./...
```

### 手动测试检查项
- [ ] 蛇身显示为**绿色**
- [ ] 食物显示为**红色**
- [ ] 边框显示为**白色**
- [ ] Game Over 文字显示为**黄色**
- [ ] 方向键/WASD 控制正常
- [ ] 吃食物加分 (+10)
- [ ] 撞墙/撞自己游戏结束
- [ ] R 键重开正常

## 4. 依赖

| 依赖 | 版本 | 用途 |
|------|------|------|
| Go | 1.21+ | 语言 |
| termbox-go | latest | 终端 UI + 颜色支持 |
