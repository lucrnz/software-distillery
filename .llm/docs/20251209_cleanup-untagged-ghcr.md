# Nightly GHCR untagged cleanup

- Added `.github/workflows/cleanup-untagged-images.yml` to run nightly at 00:00 UTC (and via manual dispatch) with concurrency guarding to avoid overlapping runs.
- Cleanup logic lives on a Python script in `tools/cleanup-untagged-images/cleanup.py`, called from the workflow after checkout.
- The script uses `GITHUB_TOKEN` with `packages: write`, detects user vs org to build the API path, paginates through all versions, and deletes every digest-only (untagged) container version for `ghcr.io/${{ github.repository }}` until none remain, then reports totals.
