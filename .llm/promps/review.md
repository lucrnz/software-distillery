## Your Task

Perform a thorough analysis of this container build infrastructure project to identify issues, bugs, security vulnerabilities, and problems that could prevent it from being production-ready. Additionally, suggest optimizations for build performance, image size, workflow efficiency, and maintainability.

## Analysis Approach

### Step 1: Understand the Codebase

Before flagging issues, first understand:
- Project structure and build targets (check `AGENTS.md` for overview)
- Dockerfile patterns (multi-stage builds, ripvex downloads)
- GitHub Actions workflow conventions
- Utility scripts and their purposes

### Step 2: Systematic Review

Analyze in this order:
1. Dockerfiles (`*/Dockerfile.alpine`, `*/Dockerfile.debian`)
2. GitHub Actions workflows (`.github/workflows/*.yml`)
3. Shell scripts (`tools/scripts/`)
4. Python utilities (`tools/*/`)
5. Project configuration and documentation

## Scope

- Analyze all Dockerfiles in project directories
- Review all GitHub Actions workflow files
- Review shell scripts and Python utilities in `tools/`
- Check for consistency with project conventions documented in `AGENTS.md`
- Do not flag the use of legacy/EOL software (e.g., Python 2.7) as an issue—this is intentional
- Do not flag missing unit tests

## What to Look For

Note: this list is not exhaustive; you can find other kinds of issues.

### Security

- Missing hash verification on downloads (should use ripvex with `--hash`)
- Hardcoded secrets, API keys, or credentials
- Overly permissive workflow permissions
- Running containers as root when unnecessary
- Downloading from untrusted sources
- Missing signature verification for critical dependencies
- Secrets leaked in build logs or layer history
- Insecure TLS configurations
- Command injection via unsanitized variables in shell scripts

### Dockerfile-Specific Issues

- Missing `SHELL` directive with proper error handling (`set -euxo pipefail`)
- Inefficient layer ordering (frequently changing layers before stable ones)
- Missing cleanup of package manager caches (`rm -rf /var/lib/apt/lists/*`, `apk --no-cache`)
- Build dependencies leaking into final stage
- Missing health checks for service containers
- Improper use of `ARG` vs `ENV`
- Missing version pinning on base images or packages
- Unnecessary files copied into final image
- Missing `.dockerignore` patterns

### Workflow-Specific Issues

- Missing `fail-fast: false` where appropriate for matrix builds
- Incorrect artifact retention periods
- Missing or incorrect `needs` dependencies between jobs
- Improper handling of `test_prefix` for test builds
- Missing GHCR login before push operations
- Inconsistent `ARTIFACT_NAME` patterns
- Race conditions in parallel jobs
- Missing error handling in shell steps
- Workflows that could expose secrets via PR triggers

### Shell Script Issues

- Missing `set -euo pipefail` (or equivalent strict mode)
- Unquoted variables that could cause word splitting
- Missing input validation
- Silent failures that should produce errors
- Hardcoded paths that should be configurable
- Missing error messages for failure cases
- Race conditions in concurrent execution

### Python Script Issues

- Missing type hints in public APIs
- Bare `except` clauses that swallow errors
- Missing input validation on user-provided data
- Unclosed resources (files, connections)
- Missing error context in exceptions
- Security issues with `urllib` usage (e.g., missing timeout)

### Reliability

- Builds that depend on external resources without fallbacks
- Missing smoke tests in final container stages
- Flaky operations without retry logic
- Missing timeouts on network operations
- Unbounded resource consumption

## Issue Severity Categories

### Critical

Immediate security risks, data loss/corruption potential, or catastrophic failures:
- Credential exposure in images or logs
- Remote code execution vulnerabilities
- Supply chain attacks (unverified downloads)
- Builds that push malicious content

### High

Significant problems under certain conditions:
- Security weaknesses requiring specific conditions to exploit
- Bugs affecting core build functionality
- Race conditions in workflow jobs
- Missing input validation on external input

### Medium

Maintainability impacts or edge case problems:
- Error handling that swallows context
- Code patterns that are brittle or error-prone
- Missing strict mode in shell scripts
- Inconsistencies across similar Dockerfiles

### Low

Minor improvements with no immediate risk:
- Code style inconsistencies
- Opportunities for simplification
- Minor documentation gaps
- Non-critical deviations from conventions

## What NOT to Flag

- Use of legacy software (Python 2.7, older Alpine versions)—this is the project's purpose
- Use of `ripvex` for downloads—this is the established pattern
- Custom runner labels (tenki-standard-*, ubicloud-*)—these are intentional
- Multi-stage builds with separate downloader stages—this is the pattern
- Installation to `/opt/<software>` with `env.sh` files—this is convention
- Use of `test_prefix` for test builds—this is the pattern
- Matrix builds covering many OS/version combinations—this is expected

## Optimization Opportunities

In addition to issues, identify opportunities to improve:

### Build Time Optimizations

- Layer caching improvements (order dependencies before code)
- Parallel download stages
- BuildKit cache mounts for package managers
- Reducing redundant operations across stages

### Image Size Optimizations

- Multi-stage build improvements
- Unnecessary packages in final stage
- Large files that could be excluded
- Compression opportunities

### Workflow Optimizations

- Job parallelization opportunities
- Caching strategies (Docker layer cache, workflow cache)
- Runner selection efficiency
- Artifact size reduction
- Matrix consolidation opportunities

### Maintainability Optimizations

- DRY opportunities across Dockerfiles
- Shared base images or stages
- Script consolidation
- Configuration externalization

## Output

Present your findings directly to the user in the following format:

```markdown
# Codebase Review

## Summary

Brief overview of the current state of the codebase, written like a senior engineer that just finished reviewing the entire project as a PR.

| Severity | Count |
|----------|-------|
| 🚨 Critical | X     |
| 🔴 High     | X     |
| 🟠 Medium   | X     |
| 🟡 Low      | X     |

## 🚨 Critical Issues

### [C-1] <Short descriptive title>
**File:** `path/to/file:line`
**Category:** Security | Reliability | etc.

**Description:**
What the issue is and why it matters.

**Code:**
```dockerfile
# The problematic code snippet
```

**Impact:**
Specific consequences if not addressed.

**Suggested Fix:**
```dockerfile
# How to fix it
```

---

## 🔴 High Issues
[Same format - Omit this section if none are found]

## 🟠 Medium Issues
[Same format - Omit this section if none are found]

## 🟡 Low Issues
[Same format - Omit this section if none are found]

[If none are found in any category, write a brief congratulations message starting with this emoji: 🎉]

## 🔧 Optimization Suggestions

### Build Time

#### [O-BT-1] <Short descriptive title>
**Files:** `affected/files`

**Current State:**
Description of current implementation.

**Suggested Optimization:**
How to improve it and expected benefit.

---

### Image Size

#### [O-IS-1] <Short descriptive title>
[Same format]

---

### Workflow Efficiency

#### [O-WF-1] <Short descriptive title>
[Same format]

---

### Maintainability

#### [O-MT-1] <Short descriptive title>
[Same format]

---

[If no optimizations found in a category, omit that subsection]

## 📝 Notes

Any observations about the codebase that don't fit into issues or optimizations but are worth mentioning.
```

