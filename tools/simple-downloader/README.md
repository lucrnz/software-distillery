# simple-downloader
A lightweight Go program for downloading files from URLs with optional SHA256 integrity verification.
Designed for simplicity and easy embedding in Docker containers or CI/CD pipelines.

It supports progress reporting, human-readable byte sizes, and quiet mode for non-interactive environments.

## Features

- **Download with Progress**: Real-time progress bar showing percentage and human-readable bytes (e.g., "1.2 MB / 5.0 GB"), updated every 500ms to prevent output spam.
- **SHA256 Verification**: Optional hash check against the downloaded file—exits with code 1 on mismatch for easy CI integration.
- **Redirect Handling**: Automatically follows HTTP redirects.
- **Quiet Mode**: Suppress all non-error output for scripts or logs.
- **Flexible Output**: Write to file (default: URL basename) or stdout (`--output -`).
- **No External Dependencies**: Pure standard library + `crypto/sha256`.

## Usage
```sh
simple-downloader [flags]
```

Run `simple-downloader --help` for full options.

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--url` | `-U` | **Required**: The URL to download (e.g., `https://example.com/file.zip`). | None |
| `--output` | `-O` | Output file path. Use `-` for stdout. Defaults to the URL's basename (or `download` if none). | URL basename |
| `--quiet` | `-q` | Suppress progress and final messages (ideal for CI/CD). Errors still printed to stderr. | `false` |
| `--hash` | `-H` | Expected SHA256 hex digest (64 chars). Verifies file integrity; exits 1 on mismatch. In quiet mode, no success message. | None |

## License
MIT License. See [LICENSE](https://github.com/lucrnz/software-distillery/blob/main/LICENSE) for details.
