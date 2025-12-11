#!/usr/bin/env python3
"""
Delete all untagged GHCR container versions for a repository.

Uses GitHub REST API with a token that has packages:write.
"""
from __future__ import annotations

import json
import sys
import time
from typing import Any, Dict, Iterable, List, Optional, Tuple

import requests
from pydantic import AliasChoices, Field
from pydantic_settings import BaseSettings, SettingsConfigDict


API_ROOT = "https://api.github.com"


class GitHubError(Exception):
    """Raised when the GitHub API returns an error."""


# Module-level session for connection pooling
_session: Optional[requests.Session] = None


def _get_session(token: str) -> requests.Session:
    """Get or create a configured requests session."""
    global _session
    if _session is None:
        _session = requests.Session()
        _session.headers.update({
            "Authorization": f"Bearer {token}",
            "Accept": "application/vnd.github+json",
            "User-Agent": "ghcr-cleanup-script",
        })
    return _session


def _request(url: str, token: str, method: str = "GET") -> str:
    """Make an HTTP request to GitHub API and return response text."""
    session = _get_session(token)
    try:
        response = session.request(method, url, timeout=30)
        response.raise_for_status()
        return response.text
    except requests.HTTPError as exc:
        status_code = exc.response.status_code if exc.response else "unknown"
        error_body = exc.response.text if exc.response else "No response"
        raise GitHubError(
            f"HTTP {status_code} for {method} {url}: {error_body}"
        ) from exc
    except requests.Timeout as exc:
        raise GitHubError(f"Request timeout for {method} {url} (30s)") from exc
    except requests.ConnectionError as exc:
        raise GitHubError(f"Connection failed for {method} {url}: {str(exc)}") from exc
    except requests.RequestException as exc:
        raise GitHubError(f"Request failed for {method} {url}: {str(exc)}") from exc


def get_owner_type(owner: str, token: str) -> str:
    """Return 'Organization' or 'User'."""
    data = _request(f"{API_ROOT}/users/{owner}", token)
    payload = json.loads(data)
    owner_type = payload.get("type")
    if owner_type not in {"User", "Organization"}:
        raise GitHubError(f"Unknown owner type for {owner}: {owner_type}")
    return owner_type


# Wait intervals for exponential backoff: 30s, 60s, then 5min (max)
WAIT_INTERVALS = [30, 60, 300]
MAX_WAIT = 300  # 5 minutes


def get_active_workflow_runs(
    owner: str,
    repo: str,
    token: str,
    exclude_run_id: Optional[str] = None,
) -> List[Dict[str, Any]]:
    """Get list of in_progress or queued workflow runs."""
    active_runs: List[Dict[str, Any]] = []
    for status in ["in_progress", "queued"]:
        url = f"{API_ROOT}/repos/{owner}/{repo}/actions/runs?status={status}&per_page=100"
        data = _request(url, token)
        payload = json.loads(data)
        for run in payload.get("workflow_runs", []):
            run_id = str(run.get("id", ""))
            if exclude_run_id and run_id == exclude_run_id:
                continue
            active_runs.append(run)
    return active_runs


def wait_for_workflows_to_complete(
    owner: str,
    repo: str,
    token: str,
    exclude_run_id: Optional[str] = None,
) -> None:
    """Wait for all workflows to complete using exponential backoff."""
    interval_index = 0

    while True:
        active_runs = get_active_workflow_runs(owner, repo, token, exclude_run_id)

        if not active_runs:
            print("No active workflows found. Proceeding with cleanup.")
            return

        # Show up to 5 workflow names
        run_names = [
            f"{r.get('name', 'Unknown')} (#{r.get('run_number', '?')})"
            for r in active_runs[:5]
        ]
        suffix = f" and {len(active_runs) - 5} more..." if len(active_runs) > 5 else ""
        print(f"Found {len(active_runs)} active workflow(s): {', '.join(run_names)}{suffix}")

        # Determine wait time using exponential backoff
        if interval_index < len(WAIT_INTERVALS):
            wait_time = WAIT_INTERVALS[interval_index]
            interval_index += 1
        else:
            wait_time = MAX_WAIT

        print(f"Waiting {wait_time} seconds before checking again...")
        time.sleep(wait_time)


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


class Settings(BaseSettings):
    """Application settings for GHCR cleanup."""

    model_config = SettingsConfigDict(
        env_prefix="",
        case_sensitive=False,
        cli_parse_args=True,
        cli_prog_name="ghcr-cleanup",
    )

    owner: Optional[str] = Field(
        default=None,
        alias="owner",
        description="Repository owner",
        validation_alias=AliasChoices("owner", "GITHUB_REPOSITORY_OWNER"),
    )
    repo: Optional[str] = Field(
        default=None,
        alias="repo",
        description="Repository name",
        validation_alias=AliasChoices("repo", "GITHUB_REPOSITORY"),
    )
    token: Optional[str] = Field(
        default=None,
        alias="token",
        description="GitHub token with packages:write permission",
        validation_alias=AliasChoices("token", "GH_TOKEN", "GITHUB_TOKEN"),
    )
    wait_for_runners: bool = Field(
        default=False,
        alias="w",
        description="Wait for running workflows to complete before cleanup",
        validation_alias=AliasChoices("wait_for_runners", "w"),
    )
    current_runner_id: Optional[str] = Field(
        default=None,
        description="Current workflow run ID to exclude from wait check",
        validation_alias=AliasChoices("current_runner_id", "GITHUB_RUN_ID"),
    )

    def validate_settings(self) -> None:
        """Validate that required settings are present."""
        if not self.owner:
            raise ValueError("Owner is required (use --owner or set GITHUB_REPOSITORY_OWNER)")
        if not self.repo:
            raise ValueError("Repo is required (use --repo or set GITHUB_REPOSITORY)")
        if not self.token:
            raise ValueError("Token is required (set GH_TOKEN or GITHUB_TOKEN)")
        if self.current_runner_id and not self.wait_for_runners:
            raise ValueError("--current-runner-id requires --wait-for-runners (-w) to be enabled")

    @property
    def repo_name(self) -> str:
        """Extract repository name from GITHUB_REPOSITORY if needed."""
        if self.repo and "/" in self.repo:
            return self.repo.split("/")[-1]
        return self.repo or ""


def main(argv: Optional[List[str]] = None) -> int:
    try:
        # Parse settings from CLI args and environment variables
        settings = Settings(_cli_parse_args=argv or sys.argv[1:])

        # Validate required settings
        settings.validate_settings()

        repo_name = settings.repo_name or settings.repo

        # Wait for workflows to complete if requested
        if settings.wait_for_runners:
            wait_for_workflows_to_complete(
                settings.owner,  # type: ignore
                repo_name,  # type: ignore
                settings.token,  # type: ignore
                settings.current_runner_id,
            )

        # Run cleanup
        count, _ = cleanup(settings.owner, repo_name, settings.token)  # type: ignore

        print(f"Deleted {count} untagged container image versions from ghcr.io/{settings.owner}/{repo_name}")
        return 0

    except ValueError as exc:
        print(f"Configuration error: {exc}", file=sys.stderr)
        return 1
    except GitHubError as exc:
        print(f"Cleanup failed: {exc}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))

