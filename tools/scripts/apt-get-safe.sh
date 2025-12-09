#!/bin/bash
set -euo pipefail

# Install only available packages; warn about the rest.

if [ $# -eq 0 ]; then
  echo "Usage: $0 <package1> [package2 ...]"
  exit 1
fi

APT_FLAGS="-y --no-install-recommends"

available=()
unavailable=()

for PACKAGE in "$@"; do
  if apt-cache show "$PACKAGE" >/dev/null 2>&1; then
    available+=("$PACKAGE")
  else
    unavailable+=("$PACKAGE")
  fi
done

if [ ${#unavailable[@]} -gt 0 ]; then
  printf "Warning: skipping unavailable packages: %s\n" "${unavailable[*]}" >&2
fi

if [ ${#available[@]} -eq 0 ]; then
  echo "No available packages to install; exiting." >&2
  exit 1
fi

exec apt-get install $APT_FLAGS "${available[@]}"
