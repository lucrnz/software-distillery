# AGENTS.md / CLAUDE.md

This file provides guidance to AI agents when working with code in this repository.

## Repository Overview

Software Distillery is a hobby project that provides binary builds for software, focusing on Linux software builds for containers. The repository builds multiple software variants using Docker multi-stage builds with multi-architecture support (amd64, arm64).

## Security Considerations

While this is a hobby project with no warranties or assurances on security updates, agents should be mindful of security when building features and suggesting changes:

- **Encourage secure practices**: When applicable, suggest security improvements such as hash verification, minimal attack surfaces, and secure defaults.
- **Legacy/EOL software**: For end-of-life platforms (e.g., Python 2.7, older Alpine/Debian versions), security is a nice-to-have but cannot be guaranteed. Focus on functionality while noting any known limitations.
- **Modern software**: For actively maintained software versions, prioritize security best practices and flag potential vulnerabilities.
- **User awareness**: Help the developer understand security trade-offs when building with older or unsupported software versions.

The goal is to build the most secure software possible within the constraints of each project, while being pragmatic about what's achievable for legacy systems.

## Project Structure

The repository contains several independent build projects:

- `python-2.7/` - Python 2.7.18 with Ubuntu security patches + OpenSSL 1.1.1zd
- `python-3.14-pyenv/` - Python 3.14 built with pyenv
- `node-nvm/` - Node.js with NVM (Node Version Manager)
- `clang-19/` - Clang 19 compiler
- `tools/scripts/` - Shared utility scripts (e.g., `apt-get-safe.sh`)

Each project directory contains:
- `Dockerfile.alpine` - Alpine Linux variant
- `Dockerfile.debian` - Debian/Ubuntu variant (uses build args to select distro)

## Build System Architecture

### Docker Build Patterns

All Dockerfiles follow a consistent multi-stage pattern:

1. **ripvex/downloader stage**: Secure downloads using `ripvex` (replaces curl/wget with hash verification)
2. **builder stage**: Compile software with all build dependencies
3. **target stage**: Minimal runtime image with only necessary runtime dependencies

### Key Build Details

#### Python 2.7 Build

- Downloads patched source from Ubuntu (version 2.7.18-13ubuntu1.5)
- Builds with custom OpenSSL 1.1.1zd (from kzalewski/openssl-1.1.1)
- Installs to `/opt/python2.7`
- Environment setup via `/opt/python2.7/env.sh` (linked to `/etc/profile.d/python2.7.sh`)
- Uses `ripvex` for secure downloads with SHA256 verification

#### Node-NVM Build

- Installs NVM to `/opt/nvm`
- May use Clang 19 (from Alpine repos or software-distillery releases)
- Environment setup via `/opt/nvm/nvm.sh` (linked to `/etc/profile.d/nvm.sh`)
- Currently builds Node 24.11.1 (other versions commented out in workflow)

## Common Development Commands

### Running builds locally

Build a specific project variant:
```bash
docker build -f python-2.7/Dockerfile.alpine --build-arg ALPINE_VERSION=3.22 -t test-python2.7 .
docker build -f node-nvm/Dockerfile.alpine --build-arg ALPINE_VERSION=3.22 --build-arg NVM_VERSION=0.40.3 --build-arg NODE_VERSION=24.11.1 -t test-node-nvm .
```

For Debian/Ubuntu variants:
```bash
docker build -f python-2.7/Dockerfile.debian --build-arg DEBIAN_VERSION=bookworm -t test-python2.7 .
docker build -f python-2.7/Dockerfile.debian --build-arg DEBIAN_VERSION=noble --build-arg DEBIAN_DISTRO=ubuntu -t test-python2.7 .
```

### Testing builds

Run the built container:
```bash
docker run --rm -it test-python2.7
# Inside container: python2.7 -c "import ssl; print(ssl.OPENSSL_VERSION)"
```

For Node builds:
```bash
docker run --rm -it test-node-nvm
# Inside container: node -v && npm -v && yarn --version
```

## Supported Platforms

Standard build matrix includes:
- Alpine: 3.15, 3.16, 3.17, 3.18, 3.19, 3.20, 3.21, 3.22
- Debian: bullseye, bookworm, trixie
- Ubuntu: focal, jammy, noble
- Architectures: linux/amd64, linux/arm64 (linux/386 commented out)

## Important Conventions

1. **ripvex for downloads**: Always use `ripvex` (from ghcr.io/lucrnz/ripvex) instead of curl/wget for security.
  - ripvex documentation can be grabbed from here: [ripvex documentation](https://raw.githubusercontent.com/lucrnz/ripvex/refs/heads/main/README.md)
  - When using ripvex containers, always tag an specific version, do not tag `latest` or `latest-dev`.
2. **Multi-stage builds**: Keep builder and target stages separate to minimize final image size
3. **Environment scripts**: Install to `/opt/<software>` and provide `env.sh` for environment setup
4. **Shell strictness**: Use `set -euxo pipefail` in Dockerfiles for proper error handling
5. **Custom runners**: Workflows use custom runner labels (tenki-standard-*, ubicloud-standard-*)

# Workflow requirements

## Triggering and inputs

- Workflows live in `.github/workflows/` and run only via `workflow_dispatch`.
- Inputs: `create_release` decides whether to publish a GitHub release; `test_build` marks the run as test-only and applies a warning prefix to every tag and artifact.

## Pipeline stages

- **prepare**: Generates the timestamped version `released-YYYYMMDD-HHMM` and sets `test_prefix` when `test_build` is true.
- **build**: Matrixed per OS/version/platform; sets up QEMU + Buildx, logs into GHCR, and runs `docker/build-push-action@v6` with `push-by-digest=true` to emit per-platform digests and upload them as artifacts.
- **merge**: Matrixed per OS/version; downloads digest artifacts (`merge-multiple: true`) and uses `docker buildx imagetools create` to publish multi-arch manifests.
- **extract**: Pulls merged tags for each platform set and packages binaries into per-arch `*.tar.gz`.
- **release**: Optional (gated by `create_release`); publishes all tarballs under the generated version tag.

## Parallelization

- Matrix jobs set `strategy.fail-fast: false` so all OS/version/platform combinations run in parallel and continue even if one variant fails.
- Build jobs pick runners per platform (amd64 on tenki-standard-*, arm64 on ubicloud-*) to maximize concurrency; digest uploads keep later stages decoupled.
- Merge and extract jobs also run as matrices to parallelize per OS/version and platform group.

## Build and tag pattern

- Build pushes per-platform images by digest with canonical naming; digests are persisted as artifacts.
- Merge tags a manifest list combining those digests as `ghcr.io/${{ github.repository }}:${test_prefix}${ARTIFACT_NAME}-${os}-${version}-${prepare.version}`. The same logical tag is reused while attaching multiple platform digests (push by digest, then tag).
- This pattern yields a single multi-platform tag while retaining per-platform digests for inspection.

## ARTIFACT_NAME and naming

- Workflows must set a fixed `ARTIFACT_NAME` (e.g., `python-2.7.18-13ubuntu1.5-openssl1.1.1zd`). Node workflows compute it per job as `node-nvm-${NVM_VERSION}-node-${matrix.node_version}`.
- `ARTIFACT_NAME` feeds image tags, artifact names, tarball filenames, and release titles. Always prepend `test_prefix` when `test_build` is enabled.

## Test build flag

- `inputs.test_build` sets `test_prefix=testbuild_donotuseinproduction_`, marks releases as prerelease, and injects a warning body.
- All tags, tarballs, and artifact names must include the prefix for test runs to keep them isolated from production names.

## Merge and extraction strategy

- Build uploads a touched file per digest; merge downloads the matching set and creates manifest lists per OS/version with `docker buildx imagetools create`.
- Extract pulls the merged tag for each platform, copies `/opt/<artifact>` out, and writes `${test_prefix}${ARTIFACT_NAME}-${os}-${version}-${arch}.tar.gz`.

## Release artifacts

- Image tags follow `ghcr.io/${{ github.repository }}:${test_prefix}${ARTIFACT_NAME}-${os}-${version}-${prepare.version}` (with the test prefix when applicable).
- GitHub releases (when `create_release` is true) bundle all per-arch tarballs under the generated version tag.

## Related Tools

- **ripvex**: Secure download tool (replacement for curl with hash verification) - https://github.com/lucrnz/ripvex
- **apt-get-safe.sh**: Utility script that installs apt packages but skips unavailable ones

# Documentation Requirements

## Change Log (.llm/docs/)

- When requesting changes to the repo, create a new markdown entry in `.llm/docs/`.
- Name files `YYYYMMDD_title.md` (example: `20251206_change-description-here.md`).
- Describe what changed and the technical reasoning; focus on decisions and context.
- Skip creating an entry if the only rationale is "user requested."
- This directory is the knowledge base for the project’s development history.

## Creating new projects

- New projects should be in their own directory, and must have a README.md.
- Add a link to the project in the project root README.md
