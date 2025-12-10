# Containerized cleanup tool with uv

- Created `tools/cleanup-untagged-images/pyproject.toml` to define the cleanup script as a proper Python project managed by uv.
- Created `tools/cleanup-untagged-images/Dockerfile` using multi-stage builds: builder stage uses `ghcr.io/astral-sh/uv:0.9-python3.14-alpine` (combined uv+Python image) to create a virtual environment, final stage uses `python:3.14-alpine` and contains only the runtime environment and script.
- Updated `.github/workflows/cleanup-untagged-images.yml` to build the Docker image and run the cleanup script from the container instead of executing the Python script directly on the runner.
- Uses Python 3.14 Alpine base image.
- Uses uv 0.9.x via Astral's official combined image, which automatically receives patch updates while staying on the same minor version.
- Simplified Dockerfile by using the pre-built uv+Python image instead of manually copying the uv binary from a separate stage.
- Passes necessary environment variables (GH_TOKEN, GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY) to the container at runtime.

## Rationale
Aligns the cleanup tool with the repository's Docker-first approach and provides consistency, isolation, and reproducibility.

Using uv provides modern Python package management and faster dependency resolution.

The containerized approach ensures the tool runs in a controlled environment with pinned dependencies, matching the patterns used for other projects in the repository (python-2.7, node-nvm, clang-19).

Using Astral's official combined uv+Python image simplifies the build process and follows Docker best practices as recommended in the official uv documentation.

Python 3.14 is the latest stable release (3.14.2) and provides modern language features. The 0.9 tag for uv provides automatic security and bug fix updates within the 0.9.x series while maintaining version stability.
