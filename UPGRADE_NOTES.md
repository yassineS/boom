# Upgrade Summary

## Changes Made

This PR successfully upgrades the boom package to use the latest htslib library (version 1.19+) instead of the embedded samtools 0.1.18.

### Go Version
- The repository now uses Go 1.24.9 (latest version)
- Created go.mod file for proper module management

### htslib Migration
Migrated from old samtools 0.1.18 API to modern htslib 1.19+ API:

**API Changes:**
- `samfile_t` → `samFile` (which is `htsFile`)
- `bam_header_t` → `sam_hdr_t` (with `bam_hdr_t` as compatibility alias)
- `samopen()` → `hts_open()`
- `samclose()` → `hts_close()`
- `samread()` → `sam_read1()`
- `samwrite()` → `sam_write1()`
- `bam_index_load()` → `hts_idx_load()` with HTS_FMT_BAI
- `bam_index_build()` → `sam_index_build()`
- `bam_get_tid()` → `sam_hdr_name2tid()`
- `bam_is_big_endian()` → `ed_is_big()`
- `bam_verbose` → `hts_set_log_level()/hts_get_log_level()`

**Structure Changes:**
- `bam1_core_t.pos`, `mpos`, `isize` are now `hts_pos_t` (int64_t) instead of `int32_t`
- `bam1_t.data_len` → `bam1_t.l_data`
- `bam1_t.l_aux` is now calculated using formula: `l_data - (n_cigar<<2) - l_qname - l_qseq - ((l_qseq + 1)>>1)`
- Removed header hash management functions (now handled internally by htslib)

**File Operations:**
- `samFdOpen()` now uses `/dev/fd/N` syntax for file descriptor operations
- Updated `samFile` struct to store both `htsFile*` and `sam_hdr_t*` separately

### Code Cleanup
- Removed all symlinks to embedded samtools-0.1.18 C files
- Updated cgo directives to link against `-lhts` instead of `-lz` and building samtools
- Updated import paths in examples from `github.com/biogo/boom` to `github.com/yassineS/boom`
- Added `.gitignore` to exclude compiled binaries

### Testing
- Verified samtest example works correctly (reading SAM files)
- Verified bamtest example works correctly (reading BAM files, writing SAM output, fetch operations)
- Both examples successfully process test data files

### Security
- CodeQL scan completed: No security vulnerabilities found
- All code changes reviewed and approved

## Breaking Changes

None for users of the Go API. The public Go interface remains the same. However:
- Users must now have htslib 1.19+ installed on their system
- The embedded samtools-0.1.18 is no longer included
- Import path changed to `github.com/yassineS/boom`

## Installation Requirements

Users need to install htslib development files:
```bash
# Ubuntu/Debian
sudo apt-get install libhts-dev

# Fedora/RHEL
sudo dnf install htslib-devel

# macOS
brew install htslib
```

## Compatibility

- Go 1.24.9 or later recommended
- htslib 1.19 or later required
- Tested on Linux (Ubuntu) with htslib 1.19+ds-1.1build3
