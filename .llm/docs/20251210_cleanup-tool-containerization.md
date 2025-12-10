# Containerized cleanup tool with uv

- Created `tools/cleanup-untagged-images/pyproject.toml` to define the cleanup script as a proper Python project managed by uv.
- Created `tools/cleanup-untagged-images/Dockerfile` using multi-stage builds: builder stage installs uv 0.9.17 and creates a virtual environment, final stage contains only the runtime environment and script.
- Updated `.github/workflows/cleanup-untagged-images.yml` to build the Docker image and run the cleanup script from the container instead of executing the Python script directly on the runner.
- Uses Python 3.13 Alpine base image for minimal footprint.
- Passes necessary environment variables (GH_TOKEN, GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY) to the container at runtime.

Rationale: aligns the cleanup tool with the repository's Docker-first approach and provides consistency, isolation, and reproducibility. Using uv provides modern Python package management and faster dependency resolution. The containerized approach ensures the tool runs in a controlled environment with pinned dependencies, matching the patterns used for other projects in the repository (python-2.7, node-nvm, clang-19). Updated to uv 0.9.17 (latest as of Dec 9, 2025) for most recent features and security updates.
