# Migrate cleanup script from urllib to requests

## Why
- PyGithub does not support the GitHub Packages API (confirmed via extensive documentation research)
- The `requests` library is more Pythonic and cleaner than urllib for HTTP operations
- Better error handling with specific exception types (HTTPError, Timeout, ConnectionError)
- Built-in connection pooling via `requests.Session()` for better performance

## What Changed

### Dependencies
- Added `requests>=2.28.0,<3.0.0` to `tools/cleanup-untagged-images/pyproject.toml`

### Code Changes in cleanup.py
1. **Imports**: Removed `urllib.error`, `urllib.parse`, `urllib.request`; added `requests` and `Optional` from typing
2. **Session Management**: Added module-level `_session` and `_get_session()` function for connection pooling
3. **_request() Function**: Replaced urllib implementation with requests
   - Uses `session.request()` instead of `urllib.request.Request`
   - Uses `response.raise_for_status()` for error detection
   - Uses `response.text` instead of `resp.read().decode("utf-8")`
   - Added timeout support (30 seconds)
   - Enhanced error handling with specific exception types

### What Stayed the Same
- All function signatures unchanged
- CLI arguments and environment variables unchanged
- JSON parsing logic unchanged
- Pagination implementation unchanged
- Dockerfile unchanged (uv automatically installs new dependency)
- Workflow unchanged

## Benefits
1. **Cleaner Code**: `response.text` vs `resp.read().decode("utf-8")`
2. **Better Error Messages**: Specific exception types (Timeout, ConnectionError, etc.)
3. **Connection Pooling**: Session reuses TCP connections across multiple API calls
4. **Timeout Handling**: Built-in 30-second timeout prevents hanging
5. **Industry Standard**: requests is the de facto standard for HTTP in Python

## Technical Details
- **Connection Pooling**: The `_get_session()` function creates a module-level session that's reused across all API calls, improving performance for multiple requests
- **Type Hints**: Used `Optional[requests.Session]` for Python 3.8+ compatibility
- **Error Context**: Enhanced error messages include HTTP method, status code, and response body
- **Timeout**: 30-second timeout prevents the script from hanging on network issues

## Rationale
After researching PyGithub's capabilities, discovered it does not support the GitHub Packages API endpoints needed for this cleanup script (`/orgs/{owner}/packages/container/{repo}/versions`). The Packages API is a separate surface area from the core GitHub REST API that PyGithub focuses on.

Rather than using PyGithub for authentication only and mixing libraries, migrated to `requests` for a cleaner, more maintainable solution. The `requests` library provides better ergonomics than stdlib urllib while keeping the script focused on its single purpose.

## Risk Assessment
**Low Risk**: The `requests` library is extremely well-tested and stable. The migration maintains identical behavior while improving error handling. The script's functionality is straightforward (HTTP GET/DELETE operations), making the migration low-risk.
