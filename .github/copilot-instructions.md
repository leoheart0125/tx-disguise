# Copilot Instructions for tx-disguise

This project is a Go-based terminal UI tool for disguising trading activity by displaying fake system info and real-time Taiwan Futures Exchange (TAIFEX) prices. It uses the Bubble Tea TUI framework and is organized for clarity and extensibility.

## Architecture Overview
- **Entry Point:** `cmd/tx-disguise/main.go` launches the TUI and orchestrates startup.
- **Terminal UI:** All TUI logic is in `internal/tui/teaui.go`, using Bubble Tea patterns (models, messages, update/view functions).
- **Futures Data:** `internal/futures/` contains API logic, DTOs, and utilities for fetching and processing futures/actuals prices.
- **Fake System Info:** `internal/fakeinfo/` generates and manages fake system/process data, mimicking `top` output.
- **Config:** `config/config.yaml` is reserved for future configuration; currently not used in logic.
- **Shared Utilities:** `internal/shared/ringbuffer.go` provides a ring buffer for history management.

## Developer Workflows
- **Run TUI in Dev Mode:**
  ```sh
  make dev-tui
  # Runs: go run ./cmd/tx-disguise
  ```
- **Linting:**
  ```sh
  make lint
  # Runs: golangci-lint (auto-installs if missing)
  ```
- **Testing:**
  Standard Go test files are in `internal/*/*_test.go`. Use:
  ```sh
  go test ./internal/...
  ```
- **Release:**
  GitHub Actions workflow in `.github/workflows/release.yml` builds and releases Darwin/arm64 binaries.

## Project-Specific Patterns
- **TUI Updates:**
  - System info refreshes every 0.5s; futures prices every 2s; history appends every minute.
  - Use Bubble Tea's message/update/view conventions for UI logic.
- **Futures Contracts:**
  - Supports 小台 (MXF), 微台 (TMF), and TXF. Symbol selection via CLI flags (`-y`, `-z`).
- **Fake Info Generation:**
  - See `internal/fakeinfo/pure_fake_system_info.go` for how fake process/system data is generated.
- **History Management:**
  - Price history is managed via a ring buffer (`internal/shared/ringbuffer.go`).

## Integration Points
- **External Data:**
  - Futures prices are fetched from TAIFEX APIs (see `internal/futures/service.go`).
- **No persistent config yet:**
  - `config/config.yaml` is a placeholder for future features.

## Conventions
- **Commit Messages:**
  - Follow Conventional Commits (see `.github/prompts/commit_instruction.md`) and give developer in markdown format.
- **Code Quality:**
  - Use `golangci-lint` for linting and formatting.
- **Directory Structure:**
  - Keep UI, futures logic, and fake info generation in separate internal packages for clarity.

## Example: Adding a New Futures Contract
1. Add contract logic to `internal/futures/dto.go` and `service.go`.
2. Update CLI flag handling in `cmd/tx-disguise/main.go`.
3. Update TUI display logic in `internal/tui/teaui.go`.

---
For questions about unclear patterns or missing documentation, ask the user for clarification or examples from the codebase.
