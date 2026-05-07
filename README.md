# LockIn

A terminal-based focus app that blocks distracting websites during timed focus sessions.

## Requirements

- Go 1.25+
- GCC (required by `go-sqlite3` — install `build-essential` on Linux, `mingw-w64` on Windows)
- **Root/admin privileges** to modify the hosts file:
  - Linux/macOS: run with `sudo`
  - Windows: run as Administrator

## Build & Run

```bash
go build -o lockin ./cmd
sudo ./lockin        # Linux/macOS
.\lockin.exe           # Windows (run as Admin)
```

## Project Structure

```
cmd/main.go              Entry point
internal/
├── blocker/             Hosts file manipulation (block/unblock domains)
├── core/                Background session scheduler
├── models/              Data types (Session, BlockedSite, BlockedApp)
├── service/             Business logic layer
├── storage/             SQLite persistence
├── ui/                  Bubble Tea TUI (terminal UI)
└── validator/           Domain and duration validation
```

