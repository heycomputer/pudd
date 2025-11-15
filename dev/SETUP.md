# Development Environment Setup

This guide will help you set up a development environment for pudding.

## Prerequisites

You'll need a version manager to ensure you're using the correct Go version. We support both **asdf** and **mise** (both read `.tool-versions`).

## Option 1: Using mise (recommended)

mise is a faster, Rust-based alternative to asdf with the same `.tool-versions` compatibility.

### Install mise

**macOS (Homebrew):**

```bash
brew install mise
echo 'eval "$(mise activate zsh)"' >> ~/.zshrc
source ~/.zshrc
```

**macOS/Linux (curl):**

```bash
curl https://mise.run | sh
echo 'eval "$(~/.local/bin/mise activate zsh)"' >> ~/.zshrc
source ~/.zshrc
```

## Option 2: Using asdf

### Install asdf

**macOS (Homebrew):**

```bash
brew install asdf
echo -e "\n. $(brew --prefix asdf)/libexec/asdf.sh" >> ~/.zshrc
source ~/.zshrc
```

**macOS/Linux (Git):**

```bash
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.14.0
echo -e '\n. "$HOME/.asdf/asdf.sh"' >> ~/.zshrc
source ~/.zshrc
```

### Install Go plugin

```bash
# Add the Go plugin
asdf plugin add golang
```

## Clone and Setup

```bash
# Clone the repository
git clone https://github.com/heycomputer/pudding.git
cd pudding

# Use your version manager install command to install the correct Go version
mise install
# or for asdf:
# asdf install

# Check Go version matches .tool-versions
bat .tool-versions
go version

# Install Go dependencies
go mod download

# Install development tools
make tools

# Build the project
make build

# run the binary
./pd --help

# Run tests
make test
```

## Next Steps

- See [RELEASE.md](RELEASE.md) for information about the release process
- See [TODO.md](TODO.md) for current development tasks
- Run `make` to see all available targets

## Troubleshooting

**Go version mismatch:**

- Run `asdf current golang` or `mise current golang` to see which version is active
- Run `asdf install` or `mise install` to install the version from `.tool-versions`

**Command not found after install:**

- Ensure `$GOPATH/bin` is in your `PATH`
- Default location is `~/go/bin`
- Add to `.zshrc`: `export PATH="$HOME/go/bin:$PATH"`

**Tests failing:**

- Ensure you're using the correct Go version
- Run `go mod download` to fetch dependencies
- Run `make clean` then `make test`
