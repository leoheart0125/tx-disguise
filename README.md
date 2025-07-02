# tx-disguise

A terminal-based tool for disguising trading activity by displaying fake system information alongside real-time Taiwan Futures Exchange (TAIFEX) futures and actuals prices. Built with Go and Bubble Tea TUI framework.

## Features
- **Terminal UI**: Shows live system process info (from `top`) and real-time futures/actuals prices.
- **Futures Support**: Supports 小台 (MXF), 微台 (TMF), and default TXF contracts.
- **Configurable**: Uses a simple YAML config (currently empty, for future use).
- **Hot reload**: Supports live development with [air](https://github.com/air-verse/air).
- **Linting**: Integrated with golangci-lint for code quality.
- **Release Workflow**: GitHub Actions workflow for building and releasing Linux/arm64 binaries.

## Usage

```
Usage: tx-disguise [-v] [-h] [ -y | -z ]
    -v: show version  
    -h: show this help
Symbol Options:
    -y: 小台 (MXF)  
    -z: 微台 (TMF)
Example: 
    tx-disguise -y
```

## Development

### Prerequisites
- Go 1.24+

### Run in Dev Mode (with hot reload)
```
make dev
```

### Lint
```
make lint
```

## Project Structure
- `cmd/tx-disguise/main.go` — Entry point
- `internal/future/` — Futures logic, API, and utils
- `internal/tui/` — Terminal UI (Bubble Tea)
- `config/config.yaml` — Config file (reserved for future use)
- `Makefile` — Dev and lint tasks
- `.air.toml` — Air config for hot reload
- `.github/workflows/release.yml` — Release automation

## License
MIT
