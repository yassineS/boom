// Copyright ©2025 The boom Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boom

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadIndex tests loading a BAM index file
func TestLoadIndex(t *testing.T) {
	testFile := filepath.Join("testdata", "test-sort.bam")

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist", testFile)
	}

	// Load the index
	idx, err := LoadIndex(testFile)
	if err != nil {
		t.Skipf("Cannot load BAM index: %v (htslib may not be installed or index file may not exist)", err)
	}

	// Verify we got a valid index
	if idx == nil {
		t.Fatal("LoadIndex returned nil Index without error")
	}

	t.Logf("Successfully loaded index for: %s", testFile)
}

// TestLoadIndexNonexistent tests loading an index for a non-existent file
func TestLoadIndexNonexistent(t *testing.T) {
	_, err := LoadIndex("/nonexistent/file.bam")
	if err == nil {
		t.Error("Expected error when loading index for non-existent file")
	}
}

// TestBAMFetch tests fetching records from a specific region
func TestBAMFetch(t *testing.T) {
	testFile := filepath.Join("testdata", "test-sort.bam")

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

	// Load the index
	idx, err := LoadIndex(testFile)
	if err != nil {
		t.Skipf("Cannot load BAM index: %v (index file may not exist)", err)
	}

	// Get reference names
	refNames := bf.RefNames()
	if len(refNames) == 0 {
		t.Skip("No reference sequences in BAM file")
	}

	// Get reference ID for the first reference
	tid, ok := bf.RefID(refNames[0])
	if !ok {
		t.Fatalf("Failed to find reference ID for %s", refNames[0])
	}

	// Fetch records from a region
	count := 0
	fn := func(r *Record) (done bool) {
		if r == nil {
			t.Error("Fetch callback received nil Record")
			return true
		}
		count++
		t.Logf("Fetched record: %s at position %d-%d", r.Name(), r.Start(), r.End())
		// Continue fetching all records in the region
		return false
	}

	// Fetch records from the entire reference sequence
	refLen := bf.RefLengths()[0]
	ret, err := bf.Fetch(idx, tid, 0, int(refLen), fn)
	if err != nil {
		t.Fatalf("Error fetching records: %v", err)
	}

	t.Logf("Fetch returned %d, found %d records", ret, count)
}

// TestBAMFetchSpecificRegion tests fetching records from a specific genomic region
func TestBAMFetchSpecificRegion(t *testing.T) {
	testFile := filepath.Join("testdata", "test-sort.bam")

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

	// Load the index
	idx, err := LoadIndex(testFile)
	if err != nil {
		t.Skipf("Cannot load BAM index: %v (index file may not exist)", err)
	}

	// Get reference names
	refNames := bf.RefNames()
	if len(refNames) == 0 {
		t.Skip("No reference sequences in BAM file")
	}

	// Get reference ID for the first reference
	tid, ok := bf.RefID(refNames[0])
	if !ok {
		t.Fatalf("Failed to find reference ID for %s", refNames[0])
	}

	// Fetch records from a specific region
	start := 1000000
	end := 2000000
	count := 0
	fn := func(r *Record) (done bool) {
		if r == nil {
			t.Error("Fetch callback received nil Record")
			return true
		}
		count++
		// Verify the record is within the requested region
		if r.End() < start {
			t.Errorf("Record end position %d is before region start %d", r.End(), start)
		}
		return false
	}

	ret, err := bf.Fetch(idx, tid, start, end, fn)
	if err != nil {
		t.Fatalf("Error fetching records: %v", err)
	}

	t.Logf("Fetch returned %d, found %d records in region %s:%d-%d", ret, count, refNames[0], start, end)
}
