## Summary
- updated node-nvm Alpine and Debian Dockerfiles to source ripvex from `ghcr.io/lucrnz/ripvex:dev-20251209-c34d36e`
- replaced simple-downloader usage with ripvex for NVM installer, LLVM key, and Clang tarball downloads
- copied ripvex (and CA certificates) into build stages to keep secure download tooling available during builds

## Rationale
- align with repository guidance to favor ripvex for verified, pinned downloads
- pin the ripvex image tag for reproducible builds across build variants

