# Installation

This document describes how to install **esplogbridge** on a Linux system.

## Requirements

- Linux host
- Go ≥ 1.25.1
- Running Grafana Loki instance

## Recommended: install using `go install`

This method:

- installs a single static binary,
- does not require root,
- is easy to update.

```bash
go install github.com/Miklakapi/esplogbridge/cmd/esplogbridge@latest
```

Make sure your Go binary directory is in `$PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

> Note: this command affects only the current shell session.<br>
> To make the change permanent, add it to your shell configuration file<br>
> (e.g. `~/.bashrc`, `~/.zshrc`).

Verify installation:

```bash
esplogbridge --version
```

## Updating

Updating esplogbridge replaces only the binary. The configuration file is not modified.

### Update to the latest version

Re-run `go install` with the `@latest` version tag:

```bash
go install github.com/Miklakapi/esplogbridge/cmd/esplogbridge@latest
```

This downloads, builds, and installs the newest released version into the Go binary directory ($(go env GOPATH)/bin).

Verify the updated binary:

```bash
esplogbridge --help
```

### Install a specific version

To install a specific version:

```bash
go install github.com/Miklakapi/esplogbridge/cmd/esplogbridge@v1.0.2
```

## Optional: system-wide installation

For production systems, you may want the binary available system-wide.

```bash
sudo cp "$(go env GOPATH)/bin/esplogbridge" /usr/local/bin/esplogbridge
sudo chmod +x /usr/local/bin/esplogbridge
```

Verify:

```bash
which esplogbridge
sudo esplogbridge --help
```
