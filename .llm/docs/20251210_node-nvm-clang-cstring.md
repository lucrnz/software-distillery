## Summary
- Fixed Node-NVM Alpine source builds failing with multiple clang-19 toolchain issues: missing `<cstring>` headers, ARM NEON type conflicts, and linker unable to find crtbeginS.o.

## Details
- **Missing system headers**: Added `musl-dev` to builder dependencies and configured include paths with `CPLUS_INCLUDE_PATH` and `C_INCLUDE_PATH` pointing to GCC's include directories so clang can find standard C/C++ headers.
- **ARM NEON conflicts**: Added `--target=aarch64-alpine-linux-musl` flags to force clang to target Alpine's musl libc instead of glibc, preventing ARM NEON type definition conflicts.
- **Linker issues**: Added `--gcc-toolchain=/usr`, `-B/usr/lib/gcc/aarch64-alpine-linux-musl/13.2.1`, and `LIBRARY_PATH`/`LDFLAGS` to ensure clang's linker can find GCC toolchain libraries and crtbeginS.o startup files.
- **Clang detection**: Modified logic to prefer system clang19 if available (Alpine repos) over downloaded software-distillery clang, which proved more compatible.

