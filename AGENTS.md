# Repository Guidelines

## Project Structure & Module Organization
- `cmd/ghex` is the main CLI entry; `cmd/ghe` is the compact binary alias. Build artifacts land in `build/` via Make targets.
- Domain logic lives under `internal/`: `account` (account management and validation), `config` (profiles and persistence), `git`/`ssh`/`platform` helpers, `ui` (Bubble Tea TUI flows), `update` (self-update), and `uninstall`.
- Shared download helpers sit in `pkg/download`. Installation helpers reside in `scripts/` (shell and PowerShell).

## Build, Test, and Development Commands
- `make build` (or `go build ./cmd/ghex`) produces the local binary; `make run` runs it directly.
- `make test` runs the full Go test suite; `make test-prop` targets property-based tests (look for `Prop`-prefixed tests).
- `make fmt` applies `gofmt`; `make lint` runs `golangci-lint` if installed. Use `make clean` to drop `build/` outputs.
- Cross-platform builds: `make build-all` or the OS/arch-specific targets (e.g., `make build-linux-amd64`). `make install` copies the binary to `/usr/local/bin`; `make release-dry` / `make release` invoke GoReleaser.

## Coding Style & Naming Conventions
- Go 1.21 codebase; run `make fmt` before pushing. Keep imports grouped and avoid unused dependencies; use `go mod tidy` when modules change.
- Follow Go naming: exported items use `CamelCase` with clear scopes; Cobra commands and flags stay lower-kebab (e.g., `ghex update --check`).
- Prefer small, composable functions; keep UI state updates predictable and avoid global mutation outside `internal/config`.

## Testing Guidelines
- Place tests alongside code in `*_test.go`; use table-driven cases where possible. Property tests live near their subject packages (`internal/update`, `internal/account`).
- Default: `make test`. For property suites: `make test-prop` or `go test -run "Prop" ./...`.
- When adding behaviors that touch network/file system, use fakes in `internal/*` or temporary dirs; keep tests hermetic so `make test` works offline.

## Commit & Pull Request Guidelines
- Commit messages follow a light Conventional Commits style seen in history (`feat:`, `fix:`, optional scopes like `fix(update): ...`). Use imperative voice and keep the first line under ~72 chars.
- PRs should include: a short summary of intent, linked issues, commands run (`make test`, `make lint`), and screenshots or terminal captures for TUI changes.
- Keep changesets focused; prefer separate PRs for features vs. release chores. Update `VERSION` via the `make bump-*` targets when altering release numbers.

## Security & Configuration Notes
- Never commit personal tokens, SSH keys, or contents of `~/.config/ghe`. Redact logs when sharing debug output.
- For update/install flows, verify permissions before touching `/usr/local/bin` or platform-specific key stores; tests should avoid writing outside temporary paths.
