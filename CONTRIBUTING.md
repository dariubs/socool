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

Please make sure all three pass cleanly before submitting.

## How to contribute

1. Fork the repo and create a branch from `main`.
2. Make your changes. Keep commits focused — one logical change per commit.
3. Add or update tests if your change affects behaviour.
4. Open a pull request with a clear description of what you changed and why.

## Adding a new command

Each command lives in its own file (e.g. `bigfiles.go`). Follow the existing pattern:

- Define a `tea.Cmd` function that does the work in a goroutine and returns a result message.
- Add a `menuItem` entry in `main.go`.
- Handle the result message in `Update()` and render it in `View()`.

## Code style

- Standard Go formatting (`gofmt`). No exceptions.
- No unnecessary comments — code should speak for itself.
- Keep the UI consistent: orange highlights, white text, dim hints.

## Reporting issues

Open a GitHub issue with steps to reproduce, your OS, terminal emulator, and Go version.

## License

By contributing you agree that your code will be released under the [MIT License](LICENSE).
