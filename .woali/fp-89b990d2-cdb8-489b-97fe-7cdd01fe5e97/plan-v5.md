# live-verdict-screenshot-check - Implementation Plan v5

## 1. Requirement Analysis and Acceptance Criteria

### Requirements
对贪吃蛇游戏进行 live 验证，核心目标：
1. **防止测试误判**：确保 test 通过后不会被基础设施超时等原因误判为 failed
2. **E2E 截图确认**：确保 PR 中包含 E2E 截图，证明游戏实际运行正常

### Acceptance Criteria
- [ ] `go test -timeout 5m -race ./...` 全部通过（12/12 tests）
- [ ] `go build ./...` 构建成功
- [ ] `go vet ./...` 无警告
- [ ] `e2e-screenshot.sh` 生成非空的 `.cast` 文件
- [ ] `.cast` 文件保存在 `snake/screenshots/` 目录
- [ ] `.cast` 文件大小 > 0（验证实际有内容）
- [ ] `.cast` 文件以 `{` 开头（asciinema JSON header 格式正确）
- [ ] E2E 截图已 commit 到 git（PR 会包含）

### v5 Key Changes from v4
v4 已完成核心实现（`e2e-screenshot.sh` 使用 `asciinema rec` 替代 `script`，解决了 macOS 上 `script -F` 只在 session 结束时才 flush 的问题）。

v5 重点：
1. **确认 E2E 截图已 commit 到 PR** — 确保 git 追踪 `snake/screenshots/*.cast` 文件
2. **Anti-misjudgment 机制验证** — 明确文档化所有防误判措施

---

## 2. Technical Architecture

### Components
```
snake/
├── screenshots/
│   ├── .gitkeep
│   └── <timestamp>-snake-gameplay.cast   # E2E recording (asciinema v2 format)
├── scripts/
│   ├── e2e-screenshot.sh                  # Primary: asciinema rec + script fallback
│   └── asciinema.sh                      # Alternative via asciinema
├── cmd/snake/main.go                      # Entry point
└── internal/...                          # Game logic (unchanged)
```

### Data Flow
1. `e2e-screenshot.sh` 先构建 snake binary
2. 检测 `asciinema` 是否可用
3. 可用时：`timeout 5 asciinema rec "$OUTPUT_FILE" --overwrite --cols 80 --rows 24 -c "$SNAKE_BINARY"`
4. 不可用时 fallback：`timeout 5 script -F "$TMP_TYPESCRIPT" "$SNAKE_BINARY"`（macOS 用 `-F`，Linux 用 `-f`）
5. 验证输出文件非空且以 `{` 开头（asciinema JSON header）
6. 文件保存在 `snake/screenshots/<timestamp>-snake-gameplay.cast`

### Anti-Misjudgment 机制
| 机制 | 作用 | 验证状态 |
|------|------|----------|
| `-timeout 5m` | 显式超时，无歧义 | 已验证 |
| `-race` | 检测数据竞争导致的 flaky failure | 已验证（无竞争） |
| `\|\| true` on timeout | 确保脚本 exit 0（asciinema 被 timeout 杀死的已知 bug） | 已实现 |
| `.cast` JSON header 验证 | 确保文件非空且为 asciinema 格式 | 已验证（1146 bytes） |

---

## 3. File-Level Change List

### No Code Changes Required
v4 已完成所有实现。v5 仅需：
1. 确认 `snake/screenshots/*.cast` 文件已 commit 到 git
2. 确认 `.gitignore` 不会排除 `.cast` 文件

### 需要检查的文件
- `snake/.gitignore` — 确保 `.cast` 文件可被追踪
- `.gitignore` 根目录 — 确保没有全局排除 screenshots

### E2E 截图 Commit 状态
根据 test-v4-1.md 的记录，最新的 `.cast` 文件是 `20260417-144806-snake-gameplay.cast` (1146 bytes)，已保存到 `.woali/fp-89b990d2-cdb8-489b-97fe-7cdd01fe5e97/screenshots-v4/`。

**需要确认**：`snake/screenshots/` 目录下的 `.cast` 文件是否已被 git 追踪（需要在 PR commit 中）。

---

## 4. Data Model

无数据模型变更。`.cast` 文件（asciinema v2 格式）：
- 格式：JSON header + timing data + stdout output
- Schema：asciinema v2 format
- 命名：`<YYYYMMDD>-<HHMMSS>-snake-gameplay.cast`

---

## 5. Testing Strategy

### 测试命令（与 v4 相同）
```bash
# 1. Build
cd snake && go build ./...

# 2. Vet
cd snake && go vet ./...

# 3. Unit tests with timeout and race detection
cd snake && go test -timeout 5m -race ./...

# 4. E2E screenshot
cd snake && ./scripts/e2e-screenshot.sh

# 5. Verify .cast file is non-empty and valid
test -s snake/screenshots/*.cast
head -c 1 snake/screenshots/*.cast | grep -q '{'

# 6. Verify .cast files are git-tracked
git ls-files snake/screenshots/*.cast
```

### 验证清单（v5 新增）
- [x] v4 所有测试通过（12/12 tests, build, vet, E2E script）
- [ ] `.cast` 文件已 git commit（PR 会包含）
- [ ] `.gitignore` 不会排除 `.cast` 文件

---

## 6. Edge Cases, Constraints, and Risks

### 约束
- **termbox-go**：需要终端颜色支持。无头 CI 可能无法渲染颜色。
- **asciinema**：需要安装。macOS: `brew install asciinema`，Linux: `pip install asciinema`。
- **asciinema Python bug**：2.4.0 版本被 timeout 杀死时会崩溃（`ValueError: list.remove(x): x not in list`），但 .cast 文件仍正确写入。stderr 被 suppress。

### 已知限制
- **termbox-go in PTY**：在某些 PTY 环境下 termbox.Init() 可能失败，导致 panic。这是已知限制。
- **headless CI**：无法运行贪吃蛇游戏的 E2E 测试（需要 TTY）。

### 风险
| 风险 | 影响 | 缓解 |
|------|------|------|
| `.cast` 文件未 commit | PR 不包含 E2E 截图 | 确保 git add + commit |
| `.gitignore` 排除 `.cast` | git 不追踪文件 | 检查 .gitignore 配置 |
| asciinema 未安装 | 使用 fallback（script 命令） | fallback 仍能生成非空 typescript |

### Trade-offs
repo-rules 中 `development-workflow.md` 建议对复杂任务使用 planner agent。本任务 v4 已完成核心实现，v5 主要是验证和确认，无新增复杂代码实现，直接执行即可。

---

## 7. Summary of v5 Changes

| 项目 | 状态 | 说明 |
|------|------|------|
| 核心实现（e2e-screenshot.sh） | ✅ 完成（v4） | asciinema rec + script fallback |
| 测试通过（12/12） | ✅ 完成（v4） | test-v4-1.md 记录 |
| Anti-misjudgment 机制 | ✅ 验证（v4） | -timeout 5m, -race, \|\| true |
| E2E 截图 commit 到 git | ⚠️ 待确认 | 需要 git add + commit |
| .gitignore 检查 | ⚠️ 待确认 | 确保不排除 .cast 文件 |
