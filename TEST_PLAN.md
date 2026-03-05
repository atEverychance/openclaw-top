# Test Plan: openclaw-top v1.0

**Date:** 2026-03-04  
**Version:** 1.0  
**Status:** In Progress

---

## Automated Tests

### Unit Tests

| Package | Tests | Status |
|---------|-------|--------|
| `pkg/gateway` | JSON parsing, client methods | ✅ Passing |
| `pkg/models` | Model creation, state management | ✅ Passing |
| `pkg/ui` | Table, StatusBar, Help, Confirm, LogViewer | ✅ Passing |

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test -v ./pkg/ui
```

---

## Manual Test Scenarios

### Feature: Kill Agent ('k' key)

| Test ID | Scenario | Steps | Expected Result | Status |
|---------|----------|-------|-----------------|--------|
| KILL-01 | Kill running agent | 1. Select RUNNING agent<br>2. Press 'k'<br>3. Press 'y' | Agent killed, success message shown | ⏳ Pending |
| KILL-02 | Cancel kill | 1. Select RUNNING agent<br>2. Press 'k'<br>3. Press 'n' | Modal closes, no action taken | ⏳ Pending |
| KILL-03 | Kill idle agent | 1. Select IDLE agent<br>2. Press 'k' | Error message shown, cannot kill | ⏳ Pending |
| KILL-04 | Kill with Esc | 1. Select RUNNING agent<br>2. Press 'k'<br>3. Press Esc | Modal closes, no action taken | ⏳ Pending |

### Feature: View Logs ('l' key)

| Test ID | Scenario | Steps | Expected Result | Status |
|---------|----------|-------|-----------------|--------|
| LOGS-01 | View logs | 1. Select agent<br>2. Press 'l' | Log viewer opens with content | ⏳ Pending |
| LOGS-02 | Scroll logs | 1. Open logs<br>2. Press 'j' or '↓' | Logs scroll down | ⏳ Pending |
| LOGS-03 | Page scroll | 1. Open logs<br>2. Press PgDn | Page down in logs | ⏳ Pending |
| LOGS-04 | Exit logs | 1. Open logs<br>2. Press 'q' | Return to table view | ⏳ Pending |
| LOGS-05 | Empty logs | 1. Select agent with no logs<br>2. Press 'l' | Shows "no logs available" | ⏳ Pending |

### Feature: Navigation

| Test ID | Scenario | Steps | Expected Result | Status |
|---------|----------|-------|-----------------|--------|
| NAV-01 | Move up/down | Press 'j'/'k' or '↑'/'↓' | Selection moves | ✅ Verified |
| NAV-02 | Sort columns | Press '1', '2', '3', '4' | Table sorts by column | ✅ Verified |
| NAV-03 | Help toggle | Press '?' | Help overlay shows | ✅ Verified |
| NAV-04 | Refresh | Press 'r' | Data refreshes | ✅ Verified |
| NAV-05 | Quit | Press 'q' | Application exits | ✅ Verified |

### Feature: Build & Distribution

| Test ID | Scenario | Steps | Expected Result | Status |
|---------|----------|-------|-----------------|--------|
| BUILD-01 | go build | `go build ./cmd/openclaw-top` | Binary created | ✅ Verified |
| BUILD-02 | make build | `make build` | Binary with version info | ✅ Verified |
| BUILD-03 | Version flag | `./openclaw-top --version` | Version info displayed | ✅ Verified |
| BUILD-04 | go install | `go install` | Binary in GOPATH/bin | ⏳ Pending |

---

## Test Coverage Goals

- [x] Unit tests for all UI components
- [x] Unit tests for gateway client
- [ ] Integration tests for full workflow
- [ ] Manual testing on macOS
- [ ] Manual testing on Linux

## Known Issues

1. Gateway `KillSession` and `GetLogs` are stubs — need real OpenClaw CLI commands
2. Log streaming not yet implemented (static view only)
3. No attach mode ('a' key) yet

## Sign-off

- [ ] All automated tests passing
- [ ] Manual test scenarios completed
- [ ] Documentation updated
- [ ] Ready for release

---

*Test plan created by QA (Artemis) | openclaw-top v1.0*
