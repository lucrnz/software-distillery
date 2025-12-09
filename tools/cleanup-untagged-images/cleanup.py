#!/usr/bin/env python3
"""
Delete all untagged GHCR container versions for a repository.

Uses GitHub REST API with a token that has packages:write.
"""
from __future__ import annotations

import argparse
import json
import os
import sys
import urllib.error
import urllib.parse
import urllib.request
from typing import Iterable, List, Tuple


API_ROOT = "https://api.github.com"


class GitHubError(Exception):
    """Raised when the GitHub API returns an error."""


def _request(url: str, token: str, method: str = "GET") -> str:
    req = urllib.request.Request(
        url,
        method=method,
        headers={
            "Authorization": f"Bearer {token}",
            "Accept": "application/vnd.github+json",
            "User-Agent": "ghcr-cleanup-script",
        },
    )
    try:
        with urllib.request.urlopen(req) as resp:
            return resp.read().decode("utf-8")
    except urllib.error.HTTPError as exc:
        body = exc.read().decode("utf-8", errors="ignore")
        raise GitHubError(f"HTTP {exc.code} for {url}: {body}") from exc


def get_owner_type(owner: str, token: str) -> str:
    """Return 'Organization' or 'User'."""
    data = _request(f"{API_ROOT}/users/{owner}", token)
    payload = json.loads(data)
    owner_type = payload.get("type")
    if owner_type not in {"User", "Organization"}:
        raise GitHubError(f"Unknown owner type for {owner}: {owner_type}")
    return owner_type


def list_version_ids_without_tags(base_url: str, token: str) -> Iterable[str]:
    """Yield all version IDs with no tags, paging through results."""
    page = 1
    while True:
        page_url = f"{base_url}?state=active&per_page=100&page={page}"
        data = _request(page_url, token)
        versions = json.loads(data)
        if not versions:
            break
        for version in versions:
            tags = (
                version.get("metadata", {})
                .get("container", {})
                .get("tags", [])
            )
            if not tags:
                vid = version.get("id")
                if vid is not None:
                    yield str(vid)
        page += 1


def delete_version(base_url: str, version_id: str, token: str) -> None:
    _request(f"{base_url}/{version_id}", token, method="DELETE")


def cleanup(owner: str, repo: str, token: str) -> Tuple[int, List[str]]:
    owner_type = get_owner_type(owner, token)
    if owner_type == "Organization":
        base_url = f"{API_ROOT}/orgs/{owner}/packages/container/{repo}/versions"
    else:
        base_url = f"{API_ROOT}/users/{owner}/packages/container/{repo}/versions"

    deleted: List[str] = []
    while True:
        ids = list(list_version_ids_without_tags(base_url, token))
        if not ids:
            break
        for vid in ids:
            delete_version(base_url, vid, token)
            deleted.append(vid)
            print(f"Deleted untagged version id={vid}")
    return len(deleted), deleted


def parse_args(argv: List[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Delete all untagged GHCR container versions for a repo."
    )
    parser.add_argument(
        "--owner",
        default=os.environ.get("GITHUB_REPOSITORY_OWNER"),
        help="Repository owner (defaults to GITHUB_REPOSITORY_OWNER env).",
    )
    parser.add_argument(
        "--repo",
        default=os.environ.get("GITHUB_REPOSITORY", "").split("/")[-1],
        help="Repository name (defaults to last path of GITHUB_REPOSITORY env).",
    )
    parser.add_argument(
        "--token",
        default=os.environ.get("GH_TOKEN") or os.environ.get("GITHUB_TOKEN"),
        help="GitHub token with packages:write (defaults to GH_TOKEN/GITHUB_TOKEN).",
    )
    return parser.parse_args(argv)


def main(argv: List[str]) -> int:
    args = parse_args(argv)

    if not args.owner:
        print("Owner is required (use --owner or set GITHUB_REPOSITORY_OWNER).", file=sys.stderr)
        return 1
    if not args.repo:
        print("Repo is required (use --repo or set GITHUB_REPOSITORY).", file=sys.stderr)
        return 1
    if not args.token:
        print("Token is required (set GH_TOKEN or GITHUB_TOKEN).", file=sys.stderr)
        return 1

    try:
        count, _ = cleanup(args.owner, args.repo, args.token)
    except GitHubError as exc:
        print(f"Cleanup failed: {exc}", file=sys.stderr)
        return 1

    print(f"Deleted {count} untagged container image versions from ghcr.io/{args.owner}/{args.repo}")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))

