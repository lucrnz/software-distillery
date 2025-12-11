# Enhanced Python 2.7 Workflow with Force Upload and Timestamp Override

## Summary

Enhanced the GitHub Actions workflow for Python 2.7 builds with comprehensive force upload functionality and timestamp override capabilities. The implementation includes robust error handling, conditional dependency installation, and improved API interaction for GitHub Container Registry management.

## Technical Changes

### 1. New Workflow Inputs
- **`force_upload`**: Boolean input (default: false) to enable forced re-upload of container images even when they already exist
- **`override_timestamp`**: String input (default: empty) to override the automatically generated timestamp for release tags

### 2. Timestamp Override Logic
- Modified the version generation in the `prepare` job to use provided timestamp when `override_timestamp` is set
- Enables hotfix re-uploads without changing existing tags that users may already be using
- Falls back to automatic timestamp generation when no override is provided

### 3. Comprehensive Force Upload Implementation
- **Package Existence Check**: Validates whether the target container package exists before attempting operations
- **Version Discovery**: Fetches all package versions and identifies the specific version matching the target tag
- **Safe Deletion**: Attempts to delete existing versions with proper error handling and rollback capabilities

### 4. Conditional Dependency Management
- **jq Installation**: Only installs `jq` JSON processor if not already available on the system
- **Feedback Messages**: Provides clear status messages about dependency availability and installation
- **Efficiency**: Avoids unnecessary package installations in optimized runner environments

### 5. Enhanced API Error Handling
- **HTTP Status Code Validation**: Comprehensive validation of all API responses (200, 204, 403, 404, etc.)
- **Rate Limiting Detection**: Automatic detection and user-friendly reporting of GitHub API rate limit issues
- **Authentication Validation**: Early detection of credential and permission problems with actionable guidance
- **Response Integrity Checks**: Validates API response format and content before processing

### 6. Robust Script Architecture
- **Separated Response Parsing**: Uses curl's `-w` flag to cleanly separate HTTP status codes from response bodies
- **Structured Error Handling**: Different error paths for various scenarios (package not found, insufficient permissions, etc.)
- **Proper Exit Semantics**: Uses `exit 0` for expected conditions and `exit 1` for actual errors
- **Input Validation**: Verifies tool availability before usage

### 7. Improved Operational Experience
- **Detailed Progress Logging**: Comprehensive status messages throughout the operation
- **Actionable Error Messages**: Clear explanations of failures with suggested remediation steps
- **Permission Guidance**: Specific token scope requirements when permission errors occur
- **Graceful Degradation**: Handles expected conditions (non-existent packages) without failure

## Context and Reasoning

The workflow enhancements address several operational challenges in container registry management:

- **Force Upload Need**: Accidental digest removal can break existing image pulls, requiring forced recreation of multi-arch manifests
- **Timestamp Override**: Enables hotfix deployments and fixes without breaking existing user deployments
- **Error Handling Gaps**: Original implementation lacked comprehensive error detection and user feedback
- **Dependency Efficiency**: Conditional installation prevents unnecessary operations in CI environments

The implementation prioritizes reliability and user experience while maintaining security best practices. By providing clear error feedback and graceful handling of edge cases, the workflow becomes more suitable for production CI/CD pipelines managing critical container infrastructure.

## Security Considerations

- Maintains secure API token handling with proper scope validation
- Provides clear guidance on required token permissions
- Implements rate limiting awareness to prevent API abuse
- Validates all external API responses before processing
