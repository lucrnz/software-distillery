## Summary
- Added a `clang++-19` symlink to the clang-19 Alpine build output so consumers relying on the versioned CXX name don’t fail during configure/compiler detection.

## Details
- The clang-19 install already sets `CXX=clang++-19` in `env.sh`, but the install may only provide `clang`/`clang++`. We now create a `clang++-19` symlink (relative, inside `/opt/clang19/bin`) pointing to `clang` after installation to guarantee the versioned binary exists for downstream builds on Alpine.

