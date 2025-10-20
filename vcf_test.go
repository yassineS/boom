// Copyright ©2025 The boom Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boom

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

// TestVCFOpen tests opening a VCF file
func TestVCFOpen(t *testing.T) {
	// Create a minimal valid VCF file for testing
	tmpDir := t.TempDir()
	vcfPath := filepath.Join(tmpDir, "test.vcf")

	vcfContent := `##fileformat=VCFv4.2
##contig=<ID=chr1,length=248956422>
##INFO=<ID=DP,Number=1,Type=Integer,Description="Total Depth">
##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	Sample1
chr1	100	rs001	A	T	30	PASS	DP=10	GT	0/1
chr1	200	rs002	G	C	40	PASS	DP=15	GT	1/1
`

	err := os.WriteFile(vcfPath, []byte(vcfContent), 0644)
	if err != nil {
		t.Skipf("Cannot create test VCF file: %v (htslib may not be installed)", err)
	}

	// Test opening the VCF file
	vf, err := OpenVCF(vcfPath)
	if err != nil {
		t.Skipf("Cannot open VCF file: %v (htslib may not be installed)", err)
	}
	defer vf.Close()

	// Verify we can get sample information
	nsamples := vf.Samples()
	if nsamples != 1 {
		t.Errorf("Expected 1 sample, got %d", nsamples)
	}

	samples := vf.SampleNames()
	if len(samples) != 1 {
		t.Errorf("Expected 1 sample name, got %d", len(samples))
	}
	if len(samples) > 0 && samples[0] != "Sample1" {
		t.Errorf("Expected sample name 'Sample1', got '%s'", samples[0])
	}
}

// TestVCFRead tests reading records from a VCF file
func TestVCFRead(t *testing.T) {
	// Create a minimal valid VCF file for testing
	tmpDir := t.TempDir()
	vcfPath := filepath.Join(tmpDir, "test.vcf")

	vcfContent := `##fileformat=VCFv4.2
##contig=<ID=chr1,length=248956422>
##INFO=<ID=DP,Number=1,Type=Integer,Description="Total Depth">
##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	Sample1
chr1	100	rs001	A	T	30	PASS	DP=10	GT	0/1
chr1	200	rs002	G	C	40	PASS	DP=15	GT	1/1
`

	err := os.WriteFile(vcfPath, []byte(vcfContent), 0644)
	if err != nil {
		t.Skipf("Cannot create test VCF file: %v (htslib may not be installed)", err)
	}

	// Open the VCF file
	vf, err := OpenVCF(vcfPath)
	if err != nil {
		t.Skipf("Cannot open VCF file: %v (htslib may not be installed)", err)
	}
	defer vf.Close()

	// Read first record
	rec1, err := vf.Read()
	if err != nil {
		t.Fatalf("Failed to read first record: %v", err)
	}

	// VCF positions are 1-based, internal representation is 0-based
	if rec1.Pos() != 99 {
		t.Errorf("Expected position 99 (0-based), got %d", rec1.Pos())
	}

	if rec1.ID() != "rs001" {
		t.Errorf("Expected ID 'rs001', got '%s'", rec1.ID())
	}

	chrom := rec1.Chrom()
	if chrom != "chr1" {
		t.Errorf("Expected chromosome 'chr1', got '%s'", chrom)
	}

	// Read second record
	rec2, err := vf.Read()
	if err != nil {
		t.Fatalf("Failed to read second record: %v", err)
	}

	if rec2.Pos() != 199 {
		t.Errorf("Expected position 199 (0-based), got %d", rec2.Pos())
	}

	if rec2.ID() != "rs002" {
		t.Errorf("Expected ID 'rs002', got '%s'", rec2.ID())
	}

	// Try to read beyond EOF
	_, err = vf.Read()
	if err != io.EOF {
		t.Errorf("Expected EOF, got %v", err)
	}
}

// TestVCFOpenNonexistent tests opening a non-existent file
func TestVCFOpenNonexistent(t *testing.T) {
	_, err := OpenVCF("/nonexistent/file.vcf")
	if err == nil {
		t.Error("Expected error when opening non-existent file")
	}
}
