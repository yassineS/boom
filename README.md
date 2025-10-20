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

## Testing

The package includes comprehensive tests for handling various file formats supported by htslib:

- **SAM/BAM files** (sam_bam_test.go) - Tests for reading and writing sequence alignment data
- **VCF files** (vcf_test.go) - Tests for reading variant call data
- **FASTQ/FASTA files** (fastq_fasta_test.go) - Tests for reading sequence data
- **BAM indexing** (index_test.go) - Tests for loading and using BAM index files for random access
- **Compressed formats** - Tests for handling gzip-compressed versions (.gz)

To run tests:

```bash
go test -v
```

**Note:** Tests require htslib to be installed. If htslib is not available, tests will be skipped automatically.

Test data files are located in the `testdata/` directory and include:
- Sample FASTA/FASTQ files for sequence format validation
- Sample SAM/BAM files for alignment format validation
- Sample VCF file for variant format validation
- BAM index files for testing random access functionality
