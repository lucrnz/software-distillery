# Node-NVM Dog-fooding: Python 3.14 Pyenv and Clang 19 Integration

**Date:** 2024-12-09

## Overview

Updated `node-nvm/Dockerfile.alpine` and `node-nvm/Dockerfile.debian` to dog-food software-distillery's own builds for Python 3.14 (pyenv) and Clang 19, replacing external dependencies with internal release artifacts.

## Changes Made

### Dockerfile.alpine

#### Python 3.14 Pyenv Integration
- **Removed**: System `python3` and `py3-pip` packages from apk dependencies
- **Added**: `PYENV_ROOT=/opt/pyenv` environment variable
- **Added**: New build stage that downloads and extracts python-3.14-pyenv from software-distillery releases
  - Uses release tag: `released-20251209-0200`
  - Pyenv version: `2.6.15`
  - URL pattern: `https://github.com/lucrnz/software-distillery/releases/download/${RELEASE_TAG}/testbuild_donotuseinproduction_python-3.14-pyenv-${PYENV_VERSION}-alpine-${ALPINE_VERSION}-${RELEASE_TAG}-${ARCH}.tar.gz`
  - Implements architecture detection and normalization (`x86_64` → `amd64`, `aarch64` remains)
  - Uses `ripvex` with `--extract-archive --extract-strip-components=2`
  - Verifies installation with `python --version`
- **Modified**: NVM build step now sources `$PYENV_ROOT/env.sh` before building Node.js

#### Clang 19 Update
- **Updated**: Clang 19 download URLs from deprecated `lucrnz/test-builds` repository to `lucrnz/software-distillery`
- **Changed**: Release tag from `released-20251204-1700` to `released-20251204-2252`
- **Changed**: URL pattern to match new naming convention: `https://github.com/lucrnz/software-distillery/releases/download/${CLANG_RELEASE_TAG}/testbuild_donotuseinproduction_clang-19-${CLANG_VERSION}-alpine-${ALPINE_VERSION}-${CLANG_RELEASE_TAG}-${ARCH}.tar.gz`
- **Added**: Architecture normalization for Clang downloads (same `x86_64` → `amd64` logic)
- **Kept**: `--extract-strip-components=2` flag (required for proper tarball structure)
- **Added**: `git` package to apk dependencies (user modification, likely needed for NVM operations)

### Dockerfile.debian

#### Python 3.14 Pyenv Integration
- **Removed**: System `python3` package from apt dependencies
- **Added**: `PYENV_ROOT=/opt/pyenv` environment variable
- **Added**: `DEBIAN_DISTRO` build argument to builder stage (required for URL construction)
- **Added**: New build stage that downloads and extracts python-3.14-pyenv from software-distillery releases
  - Uses release tag: `released-20251209-0200`
  - Pyenv version: `2.6.15`
  - URL pattern: `https://github.com/lucrnz/software-distillery/releases/download/${RELEASE_TAG}/testbuild_donotuseinproduction_python-3.14-pyenv-${PYENV_VERSION}-${DEBIAN_DISTRO}-${DEBIAN_VERSION}-${RELEASE_TAG}-${ARCH}.tar.gz`
  - Implements same architecture detection and normalization as Alpine variant
  - Uses `ripvex` with `--extract-archive --extract-strip-components=2`
  - Uses quoted variable expansion: `"$PYENV_ROOT"` (user modification for shell safety)
  - Verifies installation with `python --version`
- **Modified**: Clang 20 build step now sources `$PYENV_ROOT/env.sh` before building Node.js
- **Added**: `curl` and `git` packages to apt dependencies (user modifications)

## Technical Details

### Architecture Normalization
Both Dockerfiles implement consistent architecture handling:
```bash
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
fi
```
This converts `x86_64` to `amd64` to match GitHub release tarball naming conventions, while `aarch64` is used as-is.

### Tarball Structure
- **Python 3.14 Pyenv**: Uses `--extract-strip-components=2` to skip two leading directory components
- **Clang 19**: Also uses `--extract-strip-components=2` for consistent extraction
- Both verify successful extraction by checking for `env.sh` file presence

### Build Flow Integration
The builder stage now:
1. Downloads and extracts pyenv to `/opt/pyenv`
2. Downloads Clang (Alpine) or Clang 20 (Debian) if needed
3. Sources pyenv environment via `. $PYENV_ROOT/env.sh`
4. Sources NVM environment via `. $NVM_DIR/nvm.sh`
5. Builds Node.js with yarn using the dog-fooded toolchain

## Benefits

1. **Self-reliance**: Node-NVM builds now use software-distillery's own Python and Clang builds
2. **Consistency**: Same toolchain versioning across all software-distillery projects
3. **Maintainability**: Single source of truth for Python 3.14 and Clang versions
4. **Security**: Leverages ripvex for secure, verified downloads from controlled releases
5. **Validation**: Dog-fooding ensures python-3.14-pyenv builds are production-quality

## Release Dependencies

- Python 3.14 Pyenv: `released-20251209-0200` (version 2.6.15)
- Clang 19 (Alpine only): `released-20251204-2252` (version 19.1.7)
- Ripvex container: `ghcr.io/lucrnz/ripvex:dev-20251209-c34d36e`

## Notes

- Debian variant uses Clang 20 from LLVM apt repositories (not dog-fooded)
- Alpine variant may still use system `clang19` package if available in repos
- User added `git` package dependency (likely required for NVM git operations)
- User improved shell safety with quoted variable expansion in Debian variant
