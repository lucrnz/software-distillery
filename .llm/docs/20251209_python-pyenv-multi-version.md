## 20251209 - Python pyenv multi-version workflow

- Added multi-version matrix for Python 3.14.2, 3.13.11, 3.12.12, 3.11.14, 3.10.19 across Alpine, Debian, and Ubuntu variants with amd64/arm64 builds.
- Switched workflow artifact naming to `python-${python_version}-pyenv-${PYENV_VERSION}` and fixed Dockerfile path to `python-3-pyenv/`.
- Updated Dockerfile defaults to Python 3.14.2 while keeping pyenv 2.6.15; release and artifact patterns now handle all built versions.

