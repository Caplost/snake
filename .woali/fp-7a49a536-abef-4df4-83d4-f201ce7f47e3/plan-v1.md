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
- **语言**: Python 3（跨平台，支持标准库 `curses`）
- **输入处理**: `curses` 模块捕获键盘事件（非阻塞）
- **渲染**: 字符画布，刷新终端显示

### 组件设计

```
snake_game/
├── __init__.py       # 包标识
├── main.py           # 游戏入口，初始化 curses 窗口
├── game.py           # Game 类：主循环、状态管理
├── snake.py          # Snake 类：位置、方向、增长逻辑
├── food.py           # Food 类：食物生成
├── constants.py      # 常量：网格大小、速度、符号定义
└── tests/
    ├── __init__.py
    ├── test_snake.py
    ├── test_food.py
    └── test_game.py
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
| `snake_game/__init__.py` | 新建 | 包标识文件 |
| `snake_game/main.py` | 新建 | 入口，curses 初始化，main() 函数 |
| `snake_game/game.py` | 新建 | Game 类：主循环、状态管理、渲染 |
| `snake_game/snake.py` | 新建 | Snake 类：移动、增长、碰撞检测 |
| `snake_game/food.py` | 新建 | Food 类：随机位置生成 |
| `snake_game/constants.py` | 新建 | 常量：网格大小(20)、速度(100ms)、符号 |
| `snake_game/tests/__init__.py` | 新建 | 测试包标识 |
| `snake_game/tests/test_snake.py` | 新建 | 蛇移动/增长/碰撞测试 |
| `snake_game/tests/test_food.py` | 新建 | 食物生成测试 |
| `snake_game/tests/test_game.py` | 新建 | 分数/游戏结束测试 |

## 4. 数据模型

### Snake 类
```python
body: List[Tuple[int, int]]   # [(x,y), ...] 蛇身坐标列表，头在前
direction: Tuple[int, int]    # (dx, dy) 方向向量，如 (1,0) 表示向右
```

### Food 类
```python
position: Tuple[int, int]     # (x, y) 食物坐标
```

### Game 状态
```python
score: int       # 当前分数
game_over: bool  # 游戏是否结束
snake: Snake     # 蛇对象
food: Food       # 食物对象
```

## 5. 测试策略

### 单元测试
- `test_snake.py`: 测试蛇移动(head位置更新)、增长(吃食物后长度+1)、180度反向移动拒绝、自身碰撞检测
- `test_food.py`: 测试食物生成位置不与蛇身重叠
- `test_game.py`: 测试分数计算(吃食物+10)、游戏结束条件(撞墙/撞自己)

### 运行测试
```bash
python -m pytest snake_game/tests/ -v
```

## 6. 边缘情况与约束

### 约束
- 仅支持支持 ANSI 转义码的终端（Linux/macOS/Windows Terminal）
- Python 3.6+

### 边缘情况处理
| 情况 | 处理方式 |
|------|----------|
| 蛇不能 180 度反向移动 | 检测新方向与当前方向，拒绝对向移动 |
| 食物不能生成在蛇身上 | 重新随机生成直到不在蛇身上 |
| 游戏结束需按键才退出 | 显示 "Press R to restart" 等待输入 |
| 窗口 resize | 固定 20x20 网格，不响应窗口大小 |

## 7. 交互流程

```
启动 → 显示游戏网格和初始蛇 → 游戏循环:
  → 捕获键盘输入
  → 更新蛇位置
  → 检测碰撞(墙/自身)
  → 检测吃到食物(增长/加分)
  → 刷新显示
  → 循环直到游戏结束
游戏结束 → 显示 Game Over 和分数 → 等待 R 键 → 重新开始
```
