# urlscan-cli

The official urlscan CLI.

## Installation

### Built Binaries

The built binaries can be found on [the releases page](https://github.com/urlscan/urlscan-cli/releases).

### macOS

```sh
brew install urlscan-cli
```

## Usage

### Configuring Your API Key

There are two ways for configuring your API key:

1. `URLSCAN_API_KEY` environment variable
2. Keyring (e.g., macOS Keychain, GNOME Keyring)

NOTE: Ordered by the higher precedence.

If you want to use the keyring, you can set it via the terminal or via standard input:

```bash
$ urlscan key set
Enter a urlscan.io API key:
# or
$ echo "<api_key>" | urlscan key set -
```

### Basic Commands

#### Scan

```bash
urlscan scan submit <url>
urlscan scan result <uuid>
urlscan scan search <query>
```

Alternatively, you can pass an argument via the standard input by passing `-`.

```bash
echo "<uuid>" | urlscan scan result -
```

See `urlscan --help` and also [the document](docs/urlscan.md) for more details.

### Proxy

`HTTP_RPOXY` and `HTTPS_PROXY` environment variables are respected by default. Additionally, you can set the proxy via `--proxy` option:

```bash
urlscan --proxy http://proxy:1234 <command>
```
