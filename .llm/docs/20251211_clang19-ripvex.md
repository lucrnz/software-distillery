## Summary
- Switched `clang-19/Dockerfile.alpine` downloader stage to use `ghcr.io/lucrnz/ripvex:dev-20251211-c34d36e` instead of curl and added the ripvex binary into an Alpine-based stage.
- Kept the download hashless because GitHub auto-generated archives can change checksums, while still keeping TLS and CA verification in place.

## Rationale
- Aligns with the repo convention to use ripvex for secure downloads and avoids maintaining a curl-specific stage.
- Using the pinned ripvex image ensures a consistent downloader tool across builds.
- The hash is omitted intentionally for GitHub archives that may reflow, preventing spurious build breaks while still using HTTPS + CA validation.

