# live-verdict-screenshot-check - Implementation Plan v4

## 1. Requirement Analysis and Acceptance Criteria

### Requirements
对贪吃蛇游戏进行 live 验证，核心目标：
1. **防止测试误判**：确保 test 通过后不会被基础设施超时等原因误判为 failed
2. **E2E 截图确认**：确保 PR 中包含 E2E 截图，证明游戏实际运行正常

### Acceptance Criteria
- [ ] `go test -timeout 5m -race ./...` 全部通过（12/12 tests）
- [ ] `go build ./...` 构建成功
- [ ] `go vet ./...` 无警告
- [ ] `e2e-screenshot.sh` 在 macOS 上能生成非空的 `.cast` 文件
- [ ] 生成的 `.cast` 文件保存在 `snake/screenshots/` 目录
- [ ] `.cast` 文件大小 > 0（验证实际有内容）

### v4 Key Changes from v1
**Problem from v1 test**: The `e2e-screenshot.sh` script was backgrounded and killed via `kill $!`, which on macOS creates an empty `.cast` file (script -F only flushes on session end).

**Fix**: Run script in foreground with `timeout` command, and explicitly invoke the snake binary inside the script session.

---

## 2. Technical Architecture

### Components
```
snake/
├── screenshots/
│   └── <timestamp>-snake-gameplay.cast   # E2E recording (asciinema format)
├── scripts/
│   ├── e2e-screenshot.sh                  # Fixed: foreground timeout approach
│   └── asciinema.sh                       # Alternative via asciinema
├── cmd/snake/main.go                      # Entry point (unchanged)
└── internal/...                          # Game logic (unchanged)
```

### Data Flow
1. `e2e-screenshot.sh` detects platform (Darwin/Linux)
2. On Darwin: `timeout 5 script -F "$OUTPUT_FILE" "$SNAKE_BINARY"` (foreground)
3. On Linux: `timeout 5 script -f "$OUTPUT_FILE" "$SNAKE_BINARY"` (foreground)
4. Snake binary runs for 5 seconds inside PTY, output recorded to .cast
5. Script exits when snake binary exits (or timeout kills it)
6. .cast file is flushed with actual content

### Why Foreground + Timeout instead of Background + Kill
- Background + kill: macOS `script -F` creates file immediately but content only flushes on session end. `kill` terminates before flush → empty file.
- Foreground + timeout: `script` runs until process exits or timeout fires. Session end triggers flush → non-empty file.

---

## 3. File-Level Change List

### Modify: `snake/scripts/e2e-screenshot.sh`

**Old approach** (broken on macOS):
```bash
script $SCRIPT_FLAGS -q "$OUTPUT_FILE" &   # background
PID=$!
sleep 3
kill $PID 2>/dev/null || true              # kill before flush
```

**New approach** (cross-platform):
```bash
# Build snake binary first
SNAKE_BINARY="$PROJECT_DIR/snake"

# Run snake inside script session with timeout (foreground)
if [[ "$(uname)" == "Darwin" ]]; then
    timeout 5 script -F "$OUTPUT_FILE" "$SNAKE_BINARY" || true
else
    timeout 5 script -f "$OUTPUT_FILE" "$SNAKE_BINARY" || true
fi
```

**Key changes**:
1. Build snake binary path and run it explicitly inside `script` session
2. Use `timeout 5 script ...` in foreground (not background + sleep + kill)
3. Keep platform detection for `-F` (macOS) vs `-f` (Linux)
4. Add `|| true` to suppress timeout exit code

### No Other File Changes
- No Go files modified
- `snake/internal/*` unchanged (tests already pass)
- `snake/screenshots/.gitkeep` already exists

---

## 4. Data Model

No data model changes. The `.cast` file (asciinema format) is a text-based recording file:
- Format: JSON header + timing data + stdout output
- Schema: asciinema v2 format (well-documented)
- Naming: `<YYYYMMDD>-<HHMMSS>-snake-gameplay.cast`

---

## 5. Testing Strategy

### Automated Tests
```bash
# 1. Build
cd snake && go build -o snake ./cmd/snake

# 2. Vet
cd snake && go vet ./...

# 3. Unit tests with timeout and race detection
cd snake && go test -timeout 5m -race ./...

# 4. E2E screenshot (fixed for macOS)
cd snake && ./scripts/e2e-screenshot.sh

# 5. Verify .cast file is non-empty
test -s snake/screenshots/*.cast && echo "Screenshot captured"
```

### Verification Checklist
- [ ] `go build -o snake ./cmd/snake` succeeds
- [ ] `go vet ./...` passes
- [ ] `go test -timeout 5m -race ./...` all 12 tests pass
- [ ] `e2e-screenshot.sh` exits 0
- [ ] `snake/screenshots/*.cast` exists and `test -s` is true (size > 0)

### Anti-Misjudgment Mechanisms
1. `-timeout 5m` on `go test` — explicit timeout, no ambiguity
2. `-race` flag — detects data races that cause flaky failures
3. `|| true` on timeout — ensures script exits 0 even if timeout kills snake

---

## 6. Edge Cases, Constraints, and Risks

### Constraints
- **termbox-go**: Requires terminal with color support. Headless CI may not render colors.
- **script command**: macOS vs Linux flag differences (-F vs -f)
- **PTYX**: Snake game needs a PTY for termbox-go to work

### Edge Cases
1. **asciinema available**: Prefer `asciinema rec` if installed (better cross-platform)
2. **script fails**: Fall back to asciinema or exit with clear error message
3. **No terminal (headless CI)**: termbox-go will panic. Accept as known limitation.
4. **Empty .cast file**: If file size is 0, test should FAIL and report the issue

### Risks
| Risk | Impact | Mitigation |
|------|--------|------------|
| termbox-go needs real TTY | E2E script may fail in headless CI | Accept as known; CI should run with pseudo-TTY |
| macOS `script -F` behavior | Previous approach created empty files | Fixed: foreground timeout ensures flush |
| Snake binary not built | script has nothing to run | Script builds binary first via `go build` |

### Trade-offs (Repo Rules vs Feature Requirement)
The `development-workflow.md` suggests using **planner agent** for complex tasks. This is a focused fix (1-file change) with clear root cause from v1 testing, so direct implementation is appropriate here rather than re-planning extensively.

---

## 7. Summary of v4 Changes

| File | Change | Reason |
|------|--------|--------|
| `snake/scripts/e2e-screenshot.sh` | Rewrite background+kill to foreground+timeout | macOS `script -F` only flushes on session end; kill prevented flush |
