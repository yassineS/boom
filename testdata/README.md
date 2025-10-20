# Test Data Files

This directory contains test files for testing FASTQ and FASTA file handling in the boom package.

## Files

- `test.fastq` - Sample FASTQ file with 3 sequences
- `test.fastq.gz` - Gzip-compressed version of test.fastq
- `test.fasta` - Sample FASTA file with 3 sequences
- `test.fasta.gz` - Gzip-compressed version of test.fasta

## Format Details

### FASTQ Format
FASTQ files contain sequence data with quality scores. Each record has 4 lines:
1. Sequence identifier (starts with @)
2. Nucleotide sequence
3. Plus sign (+)
4. Quality scores (same length as sequence)

### FASTA Format
FASTA files contain sequence data without quality scores. Each record has:
1. Sequence identifier line (starts with >)
2. One or more lines of sequence data

## Usage in Tests

These test files are used by the tests in `fastq_fasta_test.go` to verify that:
- FASTQ files can be opened and handled
- FASTA files can be opened and handled
- Compressed versions (gzip) of both formats can be handled
- The htslib wrapper correctly interfaces with these file formats
