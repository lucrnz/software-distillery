# apt-get-safe.sh improvements

## What changed
- Collects requested packages into `available` and `unavailable` lists; warns to stderr about any unavailable packages instead of attempting them individually.
- Exits with status 1 when no requested packages are available to install (caller can treat this as a failure), instead of silently succeeding.
- Uses `exec apt-get install …` so the final process is `apt-get`, allowing signals/exit status to propagate directly without an extra shell.

## Why
- Batch installation reduces resolver overhead and keeps the log cleaner while maintaining consistent flags.
- Explicit stderr warnings make it clear which requested packages were skipped.
- Non-zero exit on “nothing to install” gives automation an actionable failure signal.
- Replacing the shell with `apt-get` simplifies signal handling and ensures the reported status is from the install itself.
