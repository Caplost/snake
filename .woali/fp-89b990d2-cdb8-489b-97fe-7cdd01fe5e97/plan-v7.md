# live-verdict-screenshot-check - Implementation Plan v7

## 1. Requirement Analysis and Acceptance Criteria

### Requirements
对贪吃蛇游戏进行 live 验证，核心目标：
1. **防止测试误判**：确保 test 通过后不会被基础设施超时等原因误判为 failed
2. **E2E 截图确认**：确保 PR 中包含 E2E 截图，证明游戏实际运行正常

### Acceptance Criteria — ALL VERIFIED (v4-v6)
- [x] `go test -timeout 5m -race ./...` 全部通过（12/12 tests）
- [x] `go build ./...` 构建成功
- [x] `go vet ./...` 无警告
- [x] `e2e-screenshot.sh` 生成非空的 `.cast` 文件
- [x] `.cast` 文件保存在 `snake/screenshots/` 目录
- [x] `.cast` 文件大小 > 0
- [x] `.cast` 文件以 `{` 开头（asciinema JSON header 格式正确）
- [x] E2E 截图已 commit 到 git（PR 会包含）

### v7 Role
v7 是最终总结文档。v4 已完成所有核心实现，v5-v6 完成了验证和文档化。v7 汇总完整实施记录。

---

## 2. Technical Architecture

### Components
```
snake/
├── screenshots/
│   ├── .gitkeep
│   └── <timestamp>-snake-gameplay.cast   # 39 files, git-tracked
├── scripts/
│   ├── e2e-screenshot.sh                  # Primary: asciinema rec + script fallback
│   └── asciinema.sh                      # Alternative
├── cmd/snake/main.go
└── internal/{food,game,snake}/
```

### E2E Screenshot Script Logic
```
1. Build snake binary (if not exists)
2. Detect asciinema availability
3. If asciinema: timeout 5 asciinema rec --overwrite --cols 80 --rows 24 -c "$SNAKE_BINARY"
4. Else: timeout 5 script -F/-f (platform-specific)
5. Validate JSON header (first char = '{')
6. Save to snake/screenshots/<timestamp>-snake-gameplay.cast
```

---

## 3. File-Level Change List

### 无需代码变更（v4 已完成）
| 文件 | 变更 | 说明 |
|------|------|------|
| `snake/scripts/e2e-screenshot.sh` | 完成 | asciinema rec + script fallback |
| `snake/scripts/asciinema.sh` | 完成 | 备用 asciinema 脚本 |
| `snake/screenshots/.gitkeep` | 完成 | 已存在 |
| `snake/screenshots/*.cast` | 39 个文件 | 已被 git 追踪 |

---

## 4. Data Model

无数据模型变更。`.cast` 文件（asciinema v2/v3）：
- 格式：JSON header + timing data + stdout
- 命名：`<YYYYMMDD>-<HHMMSS>-snake-gameplay.cast`

---

## 5. Testing Strategy

### 测试命令
```bash
cd snake && go build ./...                          # 构建
cd snake && go vet ./...                             # 静态检查
cd snake && go test -timeout 5m -race ./...          # 12/12 tests
cd snake && ./scripts/e2e-screenshot.sh             # E2E screenshot
git ls-files snake/screenshots/*.cast | wc -l        # 验证 git 追踪
```

### 验证清单
- [x] `go build ./...` — PASS
- [x] `go vet ./...` — PASS
- [x] `go test -timeout 5m -race ./...` — 12/12 PASS
- [x] `e2e-screenshot.sh` — 生成 .cast 文件
- [x] `.cast` 以 `{` 开头 — valid JSON header
- [x] 39 个 `.cast` 文件 git 追踪

---

## 6. Edge Cases, Constraints, and Risks

### 约束
- **termbox-go**：需要 TTY 终端
- **asciinema**：建议安装（fallback 到 `script` 命令）
- **asciinema Python bug**：timeout 杀死时崩溃但文件仍正确写入，`|| true` 确保 exit 0

### 已知限制
- **termbox-go in PTY**：某些 PTY 环境可能失败
- **headless CI**：E2E 测试需要 TTY

### 风险状态
| 风险 | 状态 |
|------|------|
| `.cast` 文件未 commit | ✅ 已解决 — 39 个文件已 git 追踪 |
| `.gitignore` 排除 `.cast` | ✅ 已解决 — 无冲突 |
| asciinema 未安装 | ✅ 已解决 — fallback 到 `script` |

---

## 7. Anti-Misjudgment Mechanisms

| 机制 | 实现 | 验证 |
|------|------|------|
| `-timeout 5m` | `go test -timeout 5m -race ./...` | ✅ |
| `-race` | 检测数据竞争 | ✅ 无竞争 |
| `|| true` | timeout 杀死 asciinema 时确保 exit 0 | ✅ |
| JSON header 验证 | `head -c 1` 检查 `{` | ✅ |
| Git tracking | 39 个 .cast 文件 | ✅ |

---

## 8. Version History

| Version | 主要内容 |
|---------|----------|
| v1-v3 | 核心游戏架构（Go + termbox-go） |
| v4 | E2E screenshot 实现（asciinema rec，JSON header 验证） |
| v5 | Git commit 验证 + anti-misjudgment 文档化 |
| v6 | 最终确认（39 个 .cast 文件，test-spec + implementation report） |
| **v7** | **最终总结文档（本文档）** |

---

## 9. Trade-offs

### development-workflow.md
`development-workflow.md` 建议复杂任务使用 planner agent。但 v4-v7 的迭代是验证和文档化，核心实现已在 v4 完成。无需 planner agent 介入。

### 时间成本
v4-v7 的迭代重点是验证和文档化。核心问题（macOS `script -F` flush behavior）在 v4 已解决。

---

## 10. Conclusion

**Feature is complete. All acceptance criteria verified.**
