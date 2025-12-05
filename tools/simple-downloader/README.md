# simple-downloader
A lightweight Go program for downloading files from URLs with optional SHA256 integrity verification and archive extraction.
Designed for simplicity and easy embedding in Docker containers or CI/CD pipelines.

It supports progress reporting, human-readable byte sizes, and quiet mode for non-interactive environments.

## Features

- **Download with Progress**: Real-time progress bar showing percentage and human-readable bytes (e.g., "1.2 MB / 5.0 GB"), updated every 500ms to prevent output spam.
- **SHA256 Verification**: Optional hash check against the downloaded file—exits with code 1 on mismatch for easy CI integration. When outputting to stdout (`--output -`) with hash verification, the file is buffered in memory, verified, and only written to stdout if the hash matches.
- **Archive Extraction**: Extract downloaded archives automatically. Supports zip, tar, tar.gz, tar.bz2, tar.xz, and tar.zstd formats.
- **Magic Byte Detection**: Archive format detection uses file magic bytes, not extensions, for reliable format identification.
- **Zip Slip Protection**: Production-ready security against path traversal attacks in archives.
- **Redirect Handling**: Automatically follows HTTP redirects.
- **Quiet Mode**: Suppress all non-error output for scripts or logs.
- **Flexible Output**: Write to file (default: URL basename) or stdout (`--output -`).
- **Clean Piping**: All status messages (progress, hash verification, final messages) are written to stderr, keeping stdout clean for data piping.
- **Working Directory**: Change to a specific directory before any operation with `--chdir`.

## Usage
```sh
simple-downloader [flags]
```

Run `simple-downloader --help` for full options.

### Flags

#### Global

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--chdir` | `-C` | Change working directory before any operation. Panics if directory doesn't exist. | None |
| `--quiet` | `-q` | Suppress progress and final messages (ideal for CI/CD). Errors still printed to stderr. | `false` |

#### Downloader

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--url` | `-U` | **Required**: The URL to download (e.g., `https://example.com/file.zip`). | None |
| `--output` | `-O` | Output file path. Use `-` for stdout. Defaults to the URL's basename (or `download` if none). | URL basename |
| `--hash` | `-H` | Expected SHA256 hex digest (64 chars). Verifies file integrity; exits 1 on mismatch. In quiet mode, no success message. When used with `--output -`, the file is buffered in memory and only written to stdout after successful verification. | None |
| `--connect-timeout` | | Maximum time for connection establishment. | `300s` |
| `--max-time` | `-m` | Maximum total time for the entire operation (0 = unlimited). | `0` |

#### Archive Extractor

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--extract-archive` | `-x` | Extract the downloaded archive. Format auto-detected via magic bytes. | `false` |
| `--remove-archive` | | Delete archive file after successful extraction. | `true` |
| `--extract-strip-components` | | Strip N leading components from file names during extraction. | `0` |

### Supported Archive Formats

- ZIP
- TAR
- GZIP (tar.gz) 
- BZIP2 (tar.bz2)
- XZ (tar.xz)
- ZSTD (tar.zstd)

### Examples

Download and extract a tarball:
```sh
simple-downloader -U https://example.com/archive.tar.gz -x
```

Download to a specific directory and extract:
```sh
simple-downloader -U https://example.com/release.zip -C /opt/app -x
```

Download with hash verification and quiet mode:
```sh
simple-downloader -U https://example.com/file.tar.xz -H abc123... -x -q
```

Keep the archive after extraction:
```sh
simple-downloader -U https://example.com/data.tar.gz -x --remove-archive=false
```

Download to stdout with hash verification (buffered):
```sh
simple-downloader -U https://example.com/file.bin -O - -H abc123... | process-file
```

## Output Behavior

### Stdout vs Stderr
- **stdout**: Contains only the downloaded file data (when using `--output -`)
- **stderr**: Contains all status messages (progress, hash verification results, final messages, archive extraction logs)

This design ensures clean piping: `simple-downloader -U url -O - | other-tool` will only pass file data to the next command.

### Hash Verification with Stdout
When using `--output -` (stdout) together with `--hash`, the download behavior changes:
1. The file is downloaded to an in-memory buffer
2. SHA256 hash is computed and verified
3. If hash matches: buffer is written to stdout
4. If hash fails: error is written to stderr (unless `--quiet`), program exits with code 1

**Note**: Large files with hash verification to stdout will consume memory proportional to file size. For very large files, consider downloading to a file first, then verifying and piping.

## License
MIT License. See [LICENSE](https://github.com/lucrnz/software-distillery/blob/main/LICENSE) for details.
