# urlscan-cli

The official urlscan CLI.

## Installation

### Built Binaries

The built binaries can be found on [the releases page](https://github.com/urlscan/urlscan-cli/releases).

> [!NOTE]
> Apple's [Gatekeeper](https://support.apple.com/guide/security/gatekeeper-and-runtime-protection-sec5599b66df/web) may prevent the binary from running.
> Use `xattr -c urlscan` to unblock it or install via Homebrew Cask (see below).

### macOS/Linux

[Homebrew](https://brew.sh/) Cask is supported for macOS/Linux:

```sh
brew install --cask urlscan/tap/urlscan-cli
```

### Windows

> [!NOTE]
> There is a [WinGet package](https://winstall.app/apps/Urlscan.urlscan-cli) but it's not maintained by urlscan.
> Use at your own risk. For any issues regarding the WinGet package, please report them to the winget-pkgs repository and not here

### Manual Build

See [the docs](./docs/dev.md#build).

## Usage

### Configuring Your API Key

There are two ways for configuring your API key:

1. `URLSCAN_API_KEY` environment variable
2. Keyring (e.g., macOS Keychain, GNOME Keyring)

> [!NOTE]
> Ordered by the higher precedence.

If you want to use the keyring, you can set it via the terminal or via standard input:

```bash
$ urlscan key set
Enter a urlscan.io API key:
# or
$ echo "<api_key>" | urlscan key set -
```

> [!NOTE]
> Keyring suport for Linux depends on [GNOME Keyring](https://wiki.gnome.org/Projects/GnomeKeyring). See [troubleshooting](./docs/troubleshooting.md#keyring) for details.

### Basic Commands

#### Scan

```bash
urlscan scan submit <url>
urlscan scan result <uuid>
urlscan search <query>
```

Alternatively, you can pass an argument via the standard input by passing `-`.

```bash
echo "<uuid>" | urlscan scan result -
```

See `urlscan --help` and also [the document](docs/urlscan.md) for more details.

### Proxy

`HTTP_PROXY` and `HTTPS_PROXY` environment variables are respected by default. Additionally, you can set the proxy via the `--proxy` option:

```bash
urlscan --proxy http://proxy:1234 <command>
```
