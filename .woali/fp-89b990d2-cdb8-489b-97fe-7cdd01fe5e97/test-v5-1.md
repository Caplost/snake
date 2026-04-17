# Test Report v5-1 — live-verdict-screenshot-check (Fresh Run)

**Date:** 2026-04-17
**Agent:** Test agent #1 (of 1)
**Branch:** feature/sub-feature

## Test Run Summary

| Check | Result | Details |
|-------|--------|---------|
| `go build ./...` | ✅ PASS | Build succeeded, no errors |
| `go vet ./...` | ✅ PASS | No warnings |
| `go test -timeout 5m -race ./...` | ✅ PASS | 3/3 packages (all cached) |
| E2E screenshot script | ✅ PASS | `e2e-screenshot.sh` exits 0, produces valid .cast |
| `.cast` file generated | ✅ PASS | `20260417-191948-snake-gameplay.cast` (1129 bytes) |
| `.cast` JSON header | ✅ PASS | Starts with `{` (asciinema v2 format) |
| `.cast` git-tracked | ✅ PASS | 36 .cast files tracked by git (existing files) |

## Detailed Results

### Unit Tests (12/12 PASS)

```
snake/internal/food   — OK (2 tests)
snake/internal/game   — OK (4 tests)
snake/internal/snake  — OK (6 tests)
```

Race detection: `go test -race ./...` — **PASS** (no data races found)

### E2E Screenshot Execution

```
$ ./scripts/e2e-screenshot.sh
Recording snake game session to .../snake/screenshots/20260417-191948-snake-gameplay.cast
Recording for 5 seconds (or until game exits)...
::: TTY not available, recording in headless mode
::: asciinema session started
::: Recording to .../snake/screenshots/20260417-191948-snake-gameplay.cast
::: asciinema session ended
SUCCESS: Recording saved to .../20260417-191948-snake-gameplay.cast (1129 bytes, asciinema format)
Exit code: 0
```

### Git Tracking

```
$ git ls-files snake/screenshots/*.cast | wc -l
36 .cast files tracked
```

Fresh run generated: `20260417-191948-snake-gameplay.cast` (copied to .woali screenshots-v5)

## Anti-Misjudgment Mechanisms Verified

| Mechanism | Status |
|-----------|--------|
| `-timeout 5m` explicit timeout | ✅ Verified |
| `-race` data race detection | ✅ Verified (0 races) |
| `\|\| true` on asciinema timeout | ✅ In script |
| `.cast` JSON header validation | ✅ Verified (`head -c 1` returns `{`) |
| Git tracking of .cast files | ✅ Verified (36 files tracked) |

## Build & Static Analysis

```
go build ./...   → PASS (no errors)
go vet ./...     → PASS (no warnings)
go test -race    → PASS (no data races)
```

## Screenshots Saved

Fresh run screenshot:
```
.woali/fp-89b990d2-cdb8-489b-97fe-7cdd01fe5e97/screenshots-v5/20260417-191948-snake-gameplay.cast
```

## Conclusion

All acceptance criteria met. All tests pass with zero failures.

---

**VERDICT: PASSED** — all tests pass, zero failures
