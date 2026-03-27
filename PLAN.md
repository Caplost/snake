# 贪吃蛇 CLI 游戏 - 实现计划

## 1. 需求分析与验收标准

### 需求概述
在终端运行的贪吃蛇游戏，支持方向键控制、计分、游戏结束/重新开始。

### 功能验收标准
- [ ] 20x20 游戏网格，蛇用 `@` 或 `█` 表示，食物用 `*` 表示
- [ ] 方向键控制蛇移动方向（WASD 或方向键）
- [ ] 蛇吃食物后长度+1，分数+10
- [ ] 撞墙或撞自己时游戏结束，显示 "Game Over" 和最终分数
- [ ] 按 R 键重新开始游戏
- [ ] 游戏帧率约 10 FPS（100ms 间隔）
- [ ] 边界和障碍物检测

## 2. 技术架构

### 技术选型
- **语言**: Python 3（跨平台，支持标准库 `curses`）
- **输入处理**: `curses` 模块捕获键盘事件（非阻塞）
- **渲染**: 字符画布，刷新终端显示

### 组件设计

```
snake_game/
├── main.py          # 游戏入口，初始化 curses 窗口
├── game.py          # Game 类：主循环、状态管理
├── snake.py         # Snake 类：位置、方向、增长逻辑
├── food.py          # Food 类：食物生成
└── constants.py     # 常量：网格大小、速度、符号定义
```

### 数据流
```
键盘输入 → Game.handle_input() → Snake.direction
                                        ↓
                              Game.update() → Snake.move()
                                        ↓
                              Game.render() → 终端绘制
```

## 3. 文件变更清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `snake_game/main.py` | 新建 | 入口，curses 初始化 |
| `snake_game/game.py` | 新建 | 主循环、状态管理 |
| `snake_game/snake.py` | 新建 | 蛇逻辑 |
| `snake_game/food.py` | 新建 | 食物逻辑 |
| `snake_game/constants.py` | 新建 | 常量定义 |

## 4. 数据模型

### Snake
```python
body: List[Tuple[int, int]]  # [(x,y), ...] 蛇身坐标列表
direction: Tuple[int, int]   # (dx, dy) 方向向量
```

### Food
```python
position: Tuple[int, int]    # (x, y) 食物坐标
```

### Game State
```python
score: int
game_over: bool
paused: bool
```

## 5. 测试策略

### 单元测试
- `test_snake.py`: 测试蛇移动、增长、碰撞检测
- `test_food.py`: 测试食物生成（不与蛇身重叠）
- `test_game.py`: 测试分数计算、游戏结束条件

### 手动测试
- 在终端运行 `python -m snake_game.main` 验证完整流程

## 6. 边缘情况与约束

### 约束
- 仅支持支持 ANSI 转义码的终端（Linux/macOS/Windows Terminal）
- Python 3.6+

### 边缘情况
- 蛇不能 180 度反向移动（立即死亡）
- 食物不能生成在蛇身上
- 游戏结束前需按键才退出
- 窗口大小调整不敏感（固定 20x20）
