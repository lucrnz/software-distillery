# Nightly GHCR cleanup guard

- Scheduled `cleanup-untagged-images` to run at 03:00 UTC (00:00 Argentina) to align with local midnight timing.
- Added a pre-cleanup gate using GitHub CLI to poll in-progress workflow runs, excluding the current run by `GITHUB_RUN_ID`.
- The gate loops with 10-minute sleeps until no other workflows are running, then proceeds to delete untagged GHCR image versions.

Rationale: avoid concurrent pipeline interference while keeping nightly cleanup aligned to the requested timezone. Leveraged preinstalled `gh` on `ubuntu-latest` for lightweight run-state checks without extra dependencies.

