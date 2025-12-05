#!/bin/bash

# This script installs packages from repositories, but skips those that are not available.

if [ $# -eq 0 ]; then
  echo "Usage: $0 <package1> [package2 ...]"
  exit 1
fi

APT_FLAGS="-y --no-install-recommends"

for PACKAGE in "$@"; do
  if apt-cache show "$PACKAGE" >/dev/null 2>&1; then
    apt-get install $APT_FLAGS "$PACKAGE"
  else
    echo "Package '$PACKAGE' is not available in the repositories."
  fi
done
