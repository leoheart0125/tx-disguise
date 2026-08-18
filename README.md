[繁體中文](README.zh.md)
# tx-disguise

A terminal-based tool for disguising trading activity by displaying fake system information alongside real-time Taiwan Futures Exchange (TAIFEX) futures and actuals prices. Built with Go and Bubble Tea TUI framework.

## Features
- **Terminal UI**: Shows live system process info (from `top`) and real-time futures/actuals prices.
- **Futures Support**: Supports 小台 (MXF), 微台 (TMF), and default TXF contracts.
- **Configurable**: Uses a simple YAML config (currently empty, for future use).
- **Linting**: Integrated with golangci-lint for code quality.
- **Release Workflow**: GitHub Actions workflow for building and releasing Darwin/arm64 (Apple Silicon) binaries.

## Usage


```
Usage: txd [-v] [-h] [ -y | -z ] [-f type]
    -v: show version
    -h: show this help
Symbol Options:
    -y: 小台 (MXF)
    -z: 微台 (TMF)
Fake Info Options:
    -f: fake info service type (top or pure)
        Default is "pure" (fully fake system info)
        "top" uses system's top command (Apple Silicon only)
Example:
    ./txd -y -f top
```

- Download the release asset for your platform (e.g., `tx-disguise-vX.Y.Z-arm64.tar.gz`).
- Unpack it:
  ```sh
  tar -xzvf tx-disguise-vX.Y.Z-arm64.tar.gz
  ```
- Run the binary:
  ```sh
  ./txd [flags]
  ```

### Terminal UI Usage

When you run `txd`, you'll see a terminal interface with three main views:

- **System Info View**: Displays fake system process info (from `top`) and current futures/actuals prices.
- **History View**: Shows a scrollable history of futures/actuals prices over the past hour.
- **Chart View**: Displays real-time line charts of futures and actuals price trends (updated every minute).

**Controls:**
- `[q]` Quit the program
- `[tab]` Switch between System Info, History and Chart views
- `[up/down]` Scroll through history (in History view)

The UI updates system info every 0.5 seconds, futures prices every 2 seconds, appends to history every minute, and updates the chart every minute.

## Development

### Prerequisites
- Go 1.24+

### Run in Dev Mode

```
make dev-tui
```
This runs the TUI directly with `go run ./cmd/tx-disguise`.

### Lint

```
make lint
```
This will auto-install `golangci-lint` if not present, then run formatting and lint checks.


## Project Structure
- `cmd/tx-disguise/main.go` — Entry point, CLI flag handling, TUI startup
- `internal/tui/teaui.go` — Terminal UI logic (Bubble Tea)
- `internal/futures/` — TAIFEX API, contract logic, DTOs, utils
- `internal/fakeinfo/` — Fake system/process info generation
- `internal/shared/ringbuffer.go` — Price history ring buffer
- `config/config.yaml` — Reserved for future config
- `Makefile` — Dev and lint tasks
- `.github/workflows/release.yml` — Release automation
- `.air.toml` — Air config for hot reload

## License
MIT
