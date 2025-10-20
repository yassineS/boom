# Test Data Files

This directory contains test files for testing various file formats supported by the boom package.

## Files

### Sequence Files
- `test.fastq` - Sample FASTQ file with 3 sequences
- `test.fastq.gz` - Gzip-compressed version of test.fastq
- `test.fasta` - Sample FASTA file with 3 sequences
- `test.fasta.gz` - Gzip-compressed version of test.fasta

### Alignment Files
- `test.sam` - Sample SAM alignment file
- `test.bam` - Sample BAM alignment file
- `test-sort.bam` - Sample sorted BAM alignment file (used for testing indexed access)
- `test-sort.bam.bai` - BAM index file for test-sort.bam

### Variant Files
- `test.vcf` - Sample VCF (Variant Call Format) file with 5 variants across 2 chromosomes

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

### SAM/BAM Format
SAM (Sequence Alignment/Map) and BAM (binary SAM) files contain aligned sequence data:
- SAM is a text-based format
- BAM is a compressed binary format
- BAI files are index files for BAM files enabling rapid random access

### VCF Format
VCF (Variant Call Format) files contain variant call data:
- VCF is a text-based format for storing gene sequence variations
- Each record represents a genomic variant with position, reference, and alternate alleles
- Can include sample genotype information and quality metrics

## Usage in Tests

These test files are used by the tests to verify that:
- FASTQ files can be opened and handled (fastq_fasta_test.go)
- FASTA files can be opened and handled (fastq_fasta_test.go)
- SAM/BAM alignment files can be opened and handled (sam_bam_test.go)
- VCF variant files can be opened and handled (vcf_test.go)
- BAM index files can be loaded and used for random access (index_test.go)
- Compressed versions (gzip) of various formats can be handled
- The htslib wrapper correctly interfaces with these file formats

Example programs in the `examples/` directory also use these test files.
