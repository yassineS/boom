# boom

samtools bindings for Go

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

## Testing

The package includes tests for handling various file formats supported by htslib:

- BAM/SAM files (alignment formats)
- FASTQ files (sequence format with quality scores)
- FASTA files (sequence format)
- Compressed versions of the above formats (.gz)

To run tests:

```bash
go test -v
```

**Note:** Tests require htslib to be installed. If htslib is not available, tests will be skipped automatically.

Test data files are located in the `testdata/` directory and include sample FASTQ and FASTA files for validation.
