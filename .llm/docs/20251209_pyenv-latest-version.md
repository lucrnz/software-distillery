## Summary

- Workflow now fetches the latest pyenv release from GitHub at runtime and exposes it to all jobs.
- Build and artifact naming now carry the discovered pyenv version, ensuring tags/tarballs reflect the actual build input.

## Rationale

- Avoids drift from hardcoded pyenv versions and keeps images aligned with upstream releases.
- Ensures consumers can see the exact pyenv version in image tags and release assets for traceability.

