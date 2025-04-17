# sol

A simple Node.js version manager for Mac with Apple Silicon chips (M1/M2/M3).

## Description

`sol` is a minimalist tool for installing and switching between different Node.js versions on macOS systems with ARM64 architecture (Apple Silicon). It was created as a personal project and **is not intended to be actively maintained**.

## Limitations

- **Only works on Mac with Apple Silicon chips** (M1, M2, M3)
- Not compatible with Intel-based macOS or any other operating system
- Basic functionality without advanced features

## Usage

```bash
# Install a specific Node.js version
sol install 20.0.0

# Switch to an installed version
sol use 20.0.0

# List installed versions
sol list

# Remove a version
sol remove 20.0.0
```

## Installation

Build the binary with:

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o sol
```

Then move the binary to a location in your PATH.

## Disclaimer

This project is experimental and not guaranteed to work in all cases. Use at your own risk.