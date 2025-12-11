# Update Python 2.7 Workflow to Use Platform-Specific Runners

## Summary

Updated `.github/workflows/build-python-2.7.yml` to use platform-specific runners following the same pattern as the Node-NVM workflow. This change improves build performance and resource utilization by matching runners to target architectures.

## Changes Made

### Runner Selection Logic

- **AMD64 builds**: Now use `tenki-standard-large-plus-16c-32g` (16-core, 32GB RAM)
- **ARM64 builds**: Now use `ubicloud-standard-16-arm` (ARM-based runner)

### Implementation Details

1. Changed `runs-on` from hardcoded `tenki-standard-large-plus-16c-32g` to dynamic `${{ matrix.runner }}`
2. Added `runner` field to all matrix entries (30 total entries across Alpine, Debian, and Ubuntu variants)
3. Each platform now selects the appropriate runner:
   - `linux/amd64` → `tenki-standard-large-plus-16c-32g`
   - `linux/arm64` → `ubicloud-standard-16-arm`

## Technical Reasoning

### Performance Optimization

- ARM64 builds now run on native ARM hardware instead of emulated x86, eliminating QEMU overhead
- AMD64 builds continue to use high-performance x86 runners with more cores and memory
- Better resource utilization across different runner types

### Cost Efficiency

- ARM workloads use ARM-specific runners which may be more cost-effective
- Prevents over-provisioning resources for ARM builds when running on x86 hardware

### Build Speed

- Native ARM execution should significantly reduce build times for ARM64 targets
- Maintains fast build times for AMD64 targets on optimized x86 hardware

## Consistency

This change brings the Python 2.7 workflow in line with the Node-NVM workflow pattern, ensuring consistent CI/CD practices across all build projects in the repository.
