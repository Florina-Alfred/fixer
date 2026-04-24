# fixer

Terminal TUI for running Docker-based training labs via Bubble Tea.

## Quick start

```
go run .                          # run with labs/ relative to cwd
go run . --labs /path/to/labs     # specify lab directory
FIXER_LABS_DIR=/path go run .     # or via env var
```

## Architecture

- `main.go` — entrypoint; resolves lab dir, loads YAML, runs TUI loop. Shell sessions exit the TUI and use `pty.ExecuteShell` for a full interactive `docker exec`. After shell exits (Ctrl+C / exit), TUI restarts.
- `labs/` — lab YAML configs and `LoadAll`/`Load` functions. Supports both flat YAML files and folder-based structure (`category/lab/lab.yml`). Expects `name`, `image`, `goal`, `validate`, `setup`, `hints`, `category`, `level`, `description`.
- `tui/` — Bubble Tea model (`model.go`), view (`view.go`). Two-level navigation: tools (categories) vertically, labs horizontally. Pressing `e` quits the TUI and hands the terminal to an interactive shell.
- `docker/` — `Client` wrapping Docker CLI calls: `Start`, `Stop`, `Remove`, `CleanUp`, `Validate`, `ContainerExists`, `IsRunning`, `Setup`.
- `pty/` — wraps `github.com/creack/pty`. `StartDockerExec` returns a `*Session`; `ExecuteShell` runs a full interactive shell (raw mode, stdin/stdout forwarding).

Lab directory resolution order: CLI args (`--labs`, `-d`) > `FIXER_LABS_DIR` env > relative to executable > relative to cwd.

## Navigation

- **↑/↓** — Navigate between tools (categories: grep, find, cut, awk, sed)
- **←/→** — Navigate between labs within a tool
- **Enter** — Start/launch container for selected lab
- **e** — Enter shell in running container
- **s** — Stop and remove container
- **r** — Reset container (stop and remove)
- **v** — Run validation checks
- **t** — Show task details and hints
- **l** — Toggle log view
- **q** — Quit

## Container States

- **● (green)** — Active: Container is running and you're in the shell
- **○ (yellow)** — Idle: Container is running but you're in the TUI
- **□** — Stopped: Container is not running

## Key details

- Container names are prefixed with `fixer-lab-` and lab names are normalized: lowercase alphanumeric + hyphens.
- **Containers must use `tail -f /dev/null` as their command** (set in `docker.go:Start`). Alpine's default `/bin/sh` exits immediately with `-d`, leaving the container in `Exited` state.
- Setup commands run automatically when a container starts (defined in lab `setup` field).
- Validation runs each `validate` entry via `docker exec`; exit code 0 = pass. Checks execute in parallel as Bubble Tea commands.
- All dependencies in `go.mod` are marked `// indirect` — they are runtime deps for a `main` module, not transitive.
- Shell mode drops out of the TUI entirely. The terminal is put into raw mode via `golang.org/x/term`, stdin is forwarded to the PTY, and PTY output is forwarded to stdout via `io.Copy`. Terminal state is restored on exit.

## Testing

```
go test ./...
```

Tests use `t.TempDir()` for temp filesystems. Docker-dependent tests (`docker/`) require Docker installed but test structural correctness — they verify return types and non-panics rather than real container operations.

## Lab YAML format

```yaml
name: Lab Title
image: python:3.11-slim
category: grep
level: intermediate
goal: Description of what to accomplish
description: Detailed description of the lab
setup:
  - "chmod +x /setup/setup.sh && /setup/setup.sh"
validate:
  - "chmod +x /validate/validate.sh && /validate/validate.sh"
hints:
  - "helpful tip"
```

`name`, `image`, and `category` are required. `goal`, `validate`, `setup`, `hints`, `level`, `description` are optional.

## Folder-based lab structure

```
labs/
  grep/
    secret-hunter/
      lab.yml
      setup/
        setup.sh
      plays/
        01-find-auth-code.md
        02-extract-secrets.md
      validate/
        validate.sh
  find/
    zombie-processes/
      ...
```

Each lab has:
- `lab.yml` — Lab configuration
- `setup/` — Setup scripts to prepare the lab environment
- `plays/` — Individual challenge files with instructions
- `validate/` — Validation scripts to check completion
