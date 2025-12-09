## Summary
- Added missing Python runtime shared libraries to the `node-nvm` Alpine target image so embedded pyenv/Python runs without missing `libbz2` and related deps.

## Context
- The Node build dogfoods `python-3-pyenv`; the Alpine target stage previously only installed minimal runtime packages, causing import failures (`_bz2` missing `libbz2.so.1`).
- Mapped build-time dev packages from `python-3-pyenv` to their runtime equivalents and added them to the target image.

## Changes
- Updated `node-nvm/Dockerfile.alpine` target stage to install `curl libbz2 zlib xz-libs ncurses-libs readline sqlite-libs gdbm expat libnsl` alongside existing basics.

