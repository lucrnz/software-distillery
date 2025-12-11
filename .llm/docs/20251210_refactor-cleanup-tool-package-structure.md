# Refactor cleanup tool to proper Python package structure

## Why

The cleanup tool was originally a single-file script (`cleanup.py`) that was difficult to maintain and lacked proper Python packaging. The refactoring addresses several architectural and maintainability issues:

1. **Package Structure**: Convert from single-file script to proper Python package with modules
2. **Configuration Management**: Replace environment variable handling with structured configuration using Pydantic
3. **Dependencies**: Migrate from urllib to requests (more robust HTTP library)
4. **Build System**: Adopt modern Python packaging standards with pyproject.toml
5. **Container Optimization**: Use combined uv+Python image for faster builds and smaller final image

## What Changed

### Project Structure
- **Removed**: `tools/cleanup-untagged-images/cleanup.py` (single-file script)
- **Added**: `tools/cleanup-untagged-images/ghcr_cleanup/` package with:
  - `__init__.py`: Package initialization
  - `main.py`: Main application logic (migrated from cleanup.py)

### Dependencies & Packaging
- **pyproject.toml**: Complete rewrite to use modern packaging standards
  - Added `pydantic-settings>=2.0.0,<3.0.0` for configuration management
  - Maintained `requests>=2.28.0,<3.0.0` for HTTP operations
  - Added console script entry point: `ghcr-cleanup = "ghcr_cleanup.main:main"`
  - Configured proper package discovery and build system

### Code Architecture
- **Configuration**: Replaced manual environment variable parsing with `pydantic_settings.BaseSettings`
  - Structured validation and type hints
  - Environment variable mapping with `Field(alias=...)`
  - Better error messages for missing/invalid configuration

- **HTTP Operations**: Migrated from urllib to requests library
  - Session-based connection pooling for performance
  - Better error handling with specific exception types
  - Timeout support (30 seconds) to prevent hanging
  - Cleaner API with `response.raise_for_status()`

- **Error Handling**: Enhanced with custom `GitHubError` exception and detailed error messages

### Container & Build System
- **Dockerfile**: Updated to use combined uv+Python image (`ghcr.io/astral-sh/uv:0.9-python3.14-alpine`)
  - Faster builds due to combined image
  - Updated Python version to 3.14
  - Proper package copying for new structure

- **Added Files**:
  - `.dockerignore`: Exclude unnecessary files from build context
  - `.gitignore`: Python-specific ignore patterns
  - `uv.lock`: Dependency lock file for reproducible builds

### Workflow Integration
- **Updated**: `.github/workflows/cleanup-untagged-images.yml`
  - Maintained same functionality but updated for containerized tool
  - Added proper permissions and environment variables
  - Kept nightly schedule and manual dispatch capability

## Benefits

1. **Maintainability**: Proper package structure with clear separation of concerns
2. **Configuration**: Type-safe configuration with validation and helpful error messages
3. **Performance**: HTTP connection pooling reduces latency for multiple API calls
4. **Reliability**: Better error handling and timeout protection
5. **Standards Compliance**: Modern Python packaging practices
6. **Build Speed**: Combined uv+Python image reduces build time and image size

## Technical Details

- **Pydantic Settings**: Uses `BaseSettings` with field aliases to map environment variables to configuration object
- **Session Management**: Module-level `requests.Session` for connection reuse across API calls
- **Error Context**: Enhanced error messages include HTTP method, status code, and response body for debugging
- **Timeout Protection**: 30-second timeout prevents network hangs
- **Package Entry Point**: Console script allows direct execution as `ghcr-cleanup` command

## Migration Approach

The refactoring maintains 100% backward compatibility:
- Same CLI interface and environment variables
- Same workflow integration
- Same GitHub API interactions
- Same cleanup logic and safety checks

## Risk Assessment

**Low Risk**: The refactoring maintains identical external behavior while improving internal architecture. The requests library is extremely stable, and pydantic-settings provides robust configuration validation. All changes were tested to ensure the cleanup functionality works identically to the previous version.
