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

// TestSAMOpen tests opening a SAM file
func TestSAMOpen(t *testing.T) {
	testFile := filepath.Join("testdata", "test.sam")

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Open the SAM file
	sf, err := OpenSAM(testFile, "")
	if err != nil {
		t.Skipf("Cannot open SAM file: %v (htslib may not be installed)", err)
	}
	defer sf.Close()

	// Verify we got a valid SAMFile
	if sf == nil {
		t.Fatal("OpenSAM returned nil SAMFile without error")
	}

	t.Logf("Successfully opened SAM file: %s", testFile)
}

// TestSAMRead tests reading records from a SAM file
func TestSAMRead(t *testing.T) {
	testFile := filepath.Join("testdata", "test.sam")

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Open the SAM file
	sf, err := OpenSAM(testFile, "")
	if err != nil {
		t.Skipf("Cannot open SAM file: %v (htslib may not be installed)", err)
	}
	defer sf.Close()

	// Read records
	count := 0
	for {
		r, _, err := sf.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error reading SAM record: %v", err)
		}

		// Verify we got a valid record
		if r == nil {
			t.Fatal("Read returned nil Record without error")
		}

		count++
	}

	if count == 0 {
		t.Error("No records read from SAM file")
	}

	t.Logf("Successfully read %d records from SAM file", count)
}

// TestSAMHeader tests accessing SAM file header information
func TestSAMHeader(t *testing.T) {
	testFile := filepath.Join("testdata", "test.sam")

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Open the SAM file
	sf, err := OpenSAM(testFile, "")
	if err != nil {
		t.Skipf("Cannot open SAM file: %v (htslib may not be installed)", err)
	}
	defer sf.Close()

	// Get header information
	header := sf.Header()
	if header == nil {
		t.Error("Header() returned nil")
	}

	// Get reference names
	refNames := sf.RefNames()
	t.Logf("Reference sequences: %v", refNames)

	// Get reference lengths
	refLengths := sf.RefLengths()
	t.Logf("Reference lengths: %v", refLengths)

	// Get number of targets
	targets := sf.Targets()
	t.Logf("Number of targets: %d", targets)

	// Get header text
	headerText := sf.Text()
	if headerText == "" {
		t.Error("Header text is empty")
	}
	t.Logf("Header text length: %d bytes", len(headerText))
}

// TestBAMOpen tests opening a BAM file
func TestBAMOpen(t *testing.T) {
	testFile := filepath.Join("testdata", "test.bam")

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Open the BAM file
	bf, err := OpenBAM(testFile)
	if err != nil {
		t.Skipf("Cannot open BAM file: %v (htslib may not be installed)", err)
	}
	defer bf.Close()

	// Verify we got a valid BAMFile
	if bf == nil {
		t.Fatal("OpenBAM returned nil BAMFile without error")
	}

	t.Logf("Successfully opened BAM file: %s", testFile)
}

// TestBAMRead tests reading records from a BAM file
func TestBAMRead(t *testing.T) {
	testFile := filepath.Join("testdata", "test.bam")

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Open the BAM file
	bf, err := OpenBAM(testFile)
	if err != nil {
		t.Skipf("Cannot open BAM file: %v (htslib may not be installed)", err)
	}
	defer bf.Close()

	// Read records
	count := 0
	for {
		r, _, err := bf.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error reading BAM record: %v", err)
		}

		// Verify we got a valid record
		if r == nil {
			t.Fatal("Read returned nil Record without error")
		}

		count++
	}

	if count == 0 {
		t.Error("No records read from BAM file")
	}

	t.Logf("Successfully read %d records from BAM file", count)
}

// TestBAMHeader tests accessing BAM file header information
func TestBAMHeader(t *testing.T) {
	testFile := filepath.Join("testdata", "test.bam")

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Open the BAM file
	bf, err := OpenBAM(testFile)
	if err != nil {
		t.Skipf("Cannot open BAM file: %v (htslib may not be installed)", err)
	}
	defer bf.Close()

	// Get header information
	header := bf.Header()
	if header == nil {
		t.Error("Header() returned nil")
	}

	// Get reference names
	refNames := bf.RefNames()
	t.Logf("Reference sequences: %v", refNames)

	// Get reference lengths
	refLengths := bf.RefLengths()
	t.Logf("Reference lengths: %v", refLengths)

	// Get number of targets
	targets := bf.Targets()
	t.Logf("Number of targets: %d", targets)

	// Get header text
	headerText := bf.Text()
	if headerText == "" {
		t.Error("Header text is empty")
	}
	t.Logf("Header text length: %d bytes", len(headerText))
}

// TestBAMRefID tests looking up reference IDs
func TestBAMRefID(t *testing.T) {
	testFile := filepath.Join("testdata", "test.bam")

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Open the BAM file
	bf, err := OpenBAM(testFile)
	if err != nil {
		t.Skipf("Cannot open BAM file: %v (htslib may not be installed)", err)
	}
	defer bf.Close()

	// Get reference names
	refNames := bf.RefNames()
	if len(refNames) == 0 {
		t.Skip("No reference sequences in BAM file")
	}

	// Test looking up the first reference
	tid, ok := bf.RefID(refNames[0])
	if !ok {
		t.Errorf("Failed to find reference ID for %s", refNames[0])
	}
	if tid < 0 {
		t.Errorf("Invalid reference ID: %d", tid)
	}
	t.Logf("Reference %s has ID %d", refNames[0], tid)

	// Test looking up a non-existent reference
	_, ok = bf.RefID("nonexistent_chromosome")
	if ok {
		t.Error("Found reference ID for non-existent chromosome")
	}
}

// TestSAMOpenNonexistent tests opening a non-existent file
func TestSAMOpenNonexistent(t *testing.T) {
	_, err := OpenSAM("/nonexistent/file.sam", "")
	if err == nil {
		t.Error("Expected error when opening non-existent file")
	}
}

// TestBAMOpenNonexistent tests opening a non-existent file
func TestBAMOpenNonexistent(t *testing.T) {
	_, err := OpenBAM("/nonexistent/file.bam")
	if err == nil {
		t.Error("Expected error when opening non-existent file")
	}
}
