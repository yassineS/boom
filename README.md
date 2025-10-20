# boom

samtools bindings for Go

This package provides Go bindings for the htslib library, supporting:
- SAM/BAM files for sequence alignment data
- VCF/BCF files for variant call data

## Installation

This package requires htslib to be installed on your system.

### Installing htslib

**Ubuntu/Debian:**
```bash
sudo apt-get install libhts-dev
```

**Fedora/RHEL:**
```bash
sudo dnf install htslib-devel
```

**macOS:**
```bash
brew install htslib
```

### Installing boom

```bash
go get github.com/yassineS/boom
```

## Requirements

- Go 1.24 or later
- htslib 1.19 or later

## Documentation

See https://www.htslib.org/ for htslib documentation.
