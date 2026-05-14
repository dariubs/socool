# socool

[![CI](https://github.com/dariubs/socool/actions/workflows/ci.yml/badge.svg)](https://github.com/dariubs/socool/actions/workflows/ci.yml)

A terminal UI toolkit for system housekeeping, written in Go.

<img src="logo.svg" alt="socool logo" width="480"/>

## Features

| Tool | What it does |
|------|-------------|
| **Biggest Files** | Top 20 largest files on your system, with one-key deletion |
| **Largest Dirs** | Top 20 directories ranked by cumulative size |
| **Duplicate Files** | Groups of identical files ranked by recoverable space |

## Install

```sh
go install github.com/dariubs/socool@latest
```

Or build from source:

```sh
git clone https://github.com/dariubs/socool
cd socool
make install
```

Requires Go 1.21+.

## Usage

```sh
socool
```

Use `↑`/`↓` (or `k`/`j`) to navigate the menu, `Enter` to run a scan, `q` to quit.

### Biggest Files

Walks from `/`, skipping virtual filesystems (`/proc`, `/sys`, `/dev`) and top-level hidden directories. Shows the top 20 files by size.

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Prompt to delete selected file |
| `y` | Confirm deletion |
| `n` / `Esc` | Cancel |
| `q` / `Esc` | Back to menu |

### Largest Dirs

Same walk as above. Each directory's size is the cumulative total of everything inside it.

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `q` / `Esc` | Back to menu |

### Duplicate Files

Groups files by content (size pre-filter, then MD5). Results are sorted by wasted space — `size × (copies − 1)`.

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `q` / `Esc` | Back to menu |

## Built with

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — terminal styling

## License

[MIT](LICENSE)
