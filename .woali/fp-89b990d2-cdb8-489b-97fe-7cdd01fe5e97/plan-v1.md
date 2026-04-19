# live-verdict-screenshot-check - 实现计划 v1

## 1. 需求分析与验收标准

### 需求概述
对贪吃蛇游戏进行 live 验证，核心目标：
1. **防止测试误判**：确保 test 通过后不会被基础设施超时等原因误判为 failed
2. **E2E 截图确认**：确保 PR 中包含 E2E 截图，证明游戏实际运行正常

### 功能验收标准
- [ ] `go test ./...` 全部通过（12/12 tests）
- [ ] `go build ./...` 构建成功
- [ ] `go test -race ./...` 无数据竞态
- [ ] 生成至少 1 张 E2E 运行截图（游戏界面截图）
- [ ] 截图文件保存到 `snake/screenshots/` 目录
- [ ] 截图在 PR 描述中可见或作为 PR comment 上传
- [ ] 消除 session stale 超时导致的测试误判（使用子进程隔离 + 超时控制）

## 2. 技术架构

### 组件设计

```
snake/
├── screenshots/              # E2E 截图存放目录（新建）
│   └── <timestamp>-snake-gameplay.png
├── scripts/
│   └── e2e-screenshot.sh    # 截图脚本（新建）
└── internal/...              # 现有代码不变
```

### E2E 截图方案
- 使用 **script** 命令（macOS/Linux）录制终端会话，生成 timing + typescript 文件
- 从 typescript 中提取游戏启动后的关键帧截图
- 或使用 **asciinema** 录制 asciicast 格式（跨平台）
- 截图需展示：蛇（黄色）、食物（红色）、边框（白色）、分数显示

### 防误判策略
- `go test` 运行时使用 `-timeout 5m` 显式设置超时，避免默认超时过短
- 使用 `go test -v ./...` 输出详细日志，便于诊断
- 使用 `script` 或 `asciinema` 录制时设置合理的录制时长（10-15 秒）
- 所有测试命令通过 `run_in_background=false` 同步执行，避免 session stale

## 3. 文件变更清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `snake/screenshots/.gitkeep` | 新建 | 截图目录占位，纳入版本控制 |
| `snake/scripts/e2e-screenshot.sh` | 新建 | E2E 截图生成脚本（script 命令） |
| `snake/scripts/asciinema.sh` | 新建 | asciinema 备选方案（跨平台） |

**无现有文件被修改** — 所有变更均为新增文件。

## 4. 数据模型

无数据模型变更。截图文件为二进制/文本输出。

截图文件命名规范：
```
<YYYYMMDD>-<HHMMSS>-snake-gameplay.png
```
或 asciinema 格式：
```
<YYYYMMDD>-<HHMMSS>-snake-gameplay.cast
```

## 5. 测试策略

### 自动化测试（CI 友好）

```bash
# 1. 构建验证
cd snake && go build ./...

# 2. 静态检查
cd snake && go vet ./...

# 3. 单元测试（设置较长超时，避免误判）
cd snake && go test -timeout 5m -race ./...

# 4. E2E 截图（headless terminal 录制）
cd snake && ./scripts/e2e-screenshot.sh
```

### E2E 截图验证检查项
- [ ] 截图显示 20x20 网格边框（ASCII 边框）
- [ ] 蛇身显示为**黄色**（非黑白色）
- [ ] 食物显示为**红色**（非黑白色）
- [ ] 分数显示在界面上
- [ ] 截图文件名符合命名规范

### 防误判机制
- 使用 `-timeout 5m` 替代默认 10m（显式控制，减少歧义）
- `go test -race` 必须通过（无数据竞态）
- 所有测试命令退出码必须为 0
- session stale 问题根源：agent session 超时导致测试进程被中断；通过 `-timeout` + 同步执行（`run_in_background=false`）缓解

## 6. 边缘情况与约束

### 约束
- `script` 命令仅 macOS/Linux 可用，Windows 需要 WSL 或 asciinema
- termbox-go 需要真实终端，CI 中可能需要 `TERM=dumb` 或 headless 方案
- E2E 截图依赖人工验证截图内容是否正确着色

### 边缘情况
- **asciinema 方案**：如 `script` 不可用，改用 asciinema录制
- **CI headless**：如 CI 无 TTY，使用 asciinema + `asciinema cat` 导出帧
- **截图验证**：着色验证依赖人工，自动化仅验证文件存在

### 风险
- termbox-go 在非交互式终端可能无法正常运行（作为风险记录，不影响主流程）
- E2E 截图无法通过单元测试自动化验证颜色正确性（需人工确认）
