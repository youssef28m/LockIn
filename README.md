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


## Usage

| Key | Action |
|-----|--------|
| `↑`/`k`, `↓`/`j` | Navigate menu |
| `enter` | Select / Start |
| `esc` | Back to home |
| `?` | Toggle help |
| `q` / `ctrl+c` | Quit |

1. **Set Timer** — configure hours/minutes, press enter to start a focus session
2. **Add website to block list** — enter a domain (e.g. `facebook.com`) to add to the block list
3. **Show Block list** — view all blocked sites, press `x` or `backspace` to remove one

When a session is active, all sites on the block list are redirected to `127.0.0.1` via your system hosts file. They are automatically unblocked when the session expires or is stopped.

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

## Cross-Platform

- **Linux/macOS**: modifies `/etc/hosts` (requires root)
- **Windows**: modifies `C:\Windows\System32\drivers\etc\hosts` (requires Administrator)
