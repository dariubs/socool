# socool

[![CI](https://github.com/dariubs/socool/actions/workflows/ci.yml/badge.svg)](https://github.com/dariubs/socool/actions/workflows/ci.yml)

A terminal UI toolkit for system tasks, written in Go.

<img src="logo.svg" alt="socool logo" width="480"/>

## Features

- **Biggest Files** — scan your system and list the top 20 largest files by size, with interactive navigation and one-key deletion
- **Largest Dirs** — scan your system and list the top 20 directories by cumulative size
- **Duplicate Files** — find duplicate files and show wasted space, ranked by bytes recoverable

## Install

**From source:**

```sh
git clone https://github.com/dariubs/socool
cd socool
go install .
```

Requires Go 1.21+.

## Usage

```sh
socool
```

Navigate with arrow keys or `j`/`k`. Press `Enter` to select.

### Biggest Files

Scans from `/` (skipping virtual filesystems like `/proc`, `/sys`, `/dev`) and displays the top 20 files ranked by size.

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Open delete confirmation |
| `y` | Confirm delete |
| `n` / `Esc` | Cancel |
| `q` / `Esc` | Back to menu |

## Built with

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — terminal styling

## License

MIT
