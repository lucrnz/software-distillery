## Summary
- Fixed Alpine node-nvm builder failure caused by Python 3.13.11 `_bz2` import missing `libbz2.so.1` during Node 24.11.1 source build fallback.
- Added the required bzip2 and related runtime libraries to `node-nvm/Dockerfile.alpine` so Python can load `_bz2` when nvm compiles Node.

## Details
- The prior Alpine builder image installed Python from the software-distillery release but lacked the shared bzip2 library, which broke `configure.py` for Node source builds after the binary tarball 404 fallback.
- Updated the builder stage packages to include `libbz2` and other standard compression/IO libs (zlib, xz-libs, readline, sqlite-libs, gdbm, expat, libffi, libnsl) to match Python’s compiled module needs and prevent future missing-SONAME errors.
- Verified locally with `docker build --no-cache --build-arg ALPINE_VERSION=3.22 -t software-distillery-node-nvm_alpine322 . -f Dockerfile.alpine`.

