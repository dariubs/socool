# Contributing to socool

Thanks for your interest in contributing. Here's everything you need to get started.

## Requirements

- Go 1.21 or later
- `make` (optional but handy)

## Getting the code

```sh
git clone https://github.com/dariubs/socool
cd socool
```

## Building and running

```sh
make build   # compile → ./socool
make run     # go run .
make install # install to $GOPATH/bin
```

## Before opening a pull request

```sh
make fmt     # format code
make lint    # run go vet
make test    # run tests
```

Please make sure all three pass cleanly before submitting. CI runs the same checks automatically on every PR.

## Project layout

```
scanner/          # pure scan logic — no UI imports
  scanner.go      # shared types (FileEntry, DupGroup), FormatSize, TopN, shouldSkip
  bigfiles.go     # FindLargestFiles
  largedirs.go    # FindLargestDirs
  dupfiles.go     # FindDuplicateFiles
  *_test.go       # tests live alongside the code they test
styles.go         # lipgloss colours and styles
model.go          # Bubble Tea model: types, Init, Update, tea.Cmd wrappers
view.go           # View() and all render* helpers
logo.go           # animated glasses logo
main.go           # entry point: main()
```

## Adding a new scanner

1. Add a `Find*` function in `scanner/` — pure Go, no UI imports.
2. Write tests for it in `scanner/*_test.go` using `t.TempDir()`.
3. Add the result message type and a `scan*()` `tea.Cmd` wrapper in `model.go`.
4. Add a `menuItem` entry to the `menuItems` slice in `model.go`.
5. Handle the result message in `Update()` and dispatch `startScan()` for the new cursor index.
6. Add a `render*` function in `view.go` and route to it from `View()`.

## Code style

- Standard Go formatting (`gofmt`). No exceptions.
- No unnecessary comments — code should speak for itself.
- Keep the UI consistent: orange highlights, white text, dim hints.
- Scanner functions are pure: they take a `root` and `n` and return data. No side effects, no UI.

## Reporting issues

Open a GitHub issue with steps to reproduce, your OS, terminal emulator, and Go version.

## License

By contributing you agree that your code will be released under the [MIT License](LICENSE).
