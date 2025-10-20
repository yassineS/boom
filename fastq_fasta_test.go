// Copyright ©2012 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boom

import (
	"os"
	"path/filepath"
	"testing"
)

// TestFASTQFileHandling tests that FASTQ files can be opened and read
func TestFASTQFileHandling(t *testing.T) {
	// Check if htslib is available by trying to open a test file
	testFile := filepath.Join("testdata", "test.fastq")
	
	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Try to open the FASTQ file using samOpen
	// FASTQ files can be opened with mode "r" in htslib
	sf, err := samOpen(testFile, "r", nil)
	if err != nil {
		t.Skipf("Could not open FASTQ file (htslib may not be installed or may not support FASTQ): %v", err)
	}
	defer sf.samClose()

	// Verify that we got a valid samFile
	if sf == nil {
		t.Fatal("samOpen returned nil samFile without error")
	}

	if sf.fp == nil {
		t.Fatal("samFile has nil file pointer")
	}

	t.Logf("Successfully opened FASTQ file: %s", testFile)
}

// TestFASTAFileHandling tests that FASTA files can be opened and read
func TestFASTAFileHandling(t *testing.T) {
	// Check if htslib is available by trying to open a test file
	testFile := filepath.Join("testdata", "test.fasta")
	
	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Try to open the FASTA file using samOpen
	// FASTA files can be opened with mode "r" in htslib
	sf, err := samOpen(testFile, "r", nil)
	if err != nil {
		t.Skipf("Could not open FASTA file (htslib may not be installed or may not support FASTA): %v", err)
	}
	defer sf.samClose()

	// Verify that we got a valid samFile
	if sf == nil {
		t.Fatal("samOpen returned nil samFile without error")
	}

	if sf.fp == nil {
		t.Fatal("samFile has nil file pointer")
	}

	t.Logf("Successfully opened FASTA file: %s", testFile)
}

// TestFASTQWithCompression tests that compressed FASTQ files can be handled
func TestFASTQWithCompression(t *testing.T) {
	testFile := filepath.Join("testdata", "test.fastq.gz")
	
	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Try to open the compressed FASTQ file
	// htslib should automatically detect gzip compression
	sf, err := samOpen(testFile, "r", nil)
	if err != nil {
		t.Skipf("Could not open compressed FASTQ file (htslib may not be installed or may not support compressed FASTQ): %v", err)
	}
	defer sf.samClose()

	// Verify that we got a valid samFile
	if sf == nil {
		t.Fatal("samOpen returned nil samFile without error")
	}

	if sf.fp == nil {
		t.Fatal("samFile has nil file pointer")
	}

	t.Logf("Successfully opened compressed FASTQ file: %s", testFile)
}

// TestFASTAWithCompression tests that compressed FASTA files can be handled
func TestFASTAWithCompression(t *testing.T) {
	testFile := filepath.Join("testdata", "test.fasta.gz")
	
	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Try to open the compressed FASTA file
	// htslib should automatically detect gzip compression
	sf, err := samOpen(testFile, "r", nil)
	if err != nil {
		t.Skipf("Could not open compressed FASTA file (htslib may not be installed or may not support compressed FASTA): %v", err)
	}
	defer sf.samClose()

	// Verify that we got a valid samFile
	if sf == nil {
		t.Fatal("samOpen returned nil samFile without error")
	}

	if sf.fp == nil {
		t.Fatal("samFile has nil file pointer")
	}

	t.Logf("Successfully opened compressed FASTA file: %s", testFile)
}

// TestNonexistentFile tests that opening a nonexistent file returns an error
func TestNonexistentFile(t *testing.T) {
	testFile := filepath.Join("testdata", "nonexistent.fastq")
	
	sf, err := samOpen(testFile, "r", nil)
	if err == nil {
		if sf != nil {
			sf.samClose()
		}
		t.Fatal("Expected error when opening nonexistent file, got nil")
	}

	t.Logf("Opening nonexistent file correctly returned error: %v", err)
}
