package main

import (
	"fmt"
	"io"
	"log"

	"github.com/yassineS/boom"
)

func main() {
	// Set verbosity to minimal
	boom.Verbosity(0)

	// Test file path - use the test VCF file in testdata
	filename := "../../testdata/test.vcf"

	// Open the VCF file
	vf, err := boom.OpenVCF(filename)
	if err != nil {
		log.Fatalf("Error opening VCF file: %v", err)
	}
	defer vf.Close()

	// Print sample information
	samples := vf.SampleNames()
	fmt.Printf("Number of samples: %d\n", vf.Samples())
	if len(samples) > 0 {
		fmt.Printf("Sample names: %v\n", samples)
	}

	// Read and print variant records
	fmt.Println("\nVariant records:")
	count := 0
	for {
		rec, err := vf.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading record: %v", err)
		}

		// Print chromosome, position (1-based for display), and ID
		fmt.Printf("Variant %d: %s:%d ID=%s\n", count+1, rec.Chrom(), rec.Pos()+1, rec.ID())
		count++

		// Limit output for demonstration
		if count >= 10 {
			fmt.Println("... (showing first 10 variants)")
			break
		}
	}

	if count == 0 {
		fmt.Println("No variants found in file")
	} else {
		fmt.Printf("\nTotal variants read: %d\n", count)
	}
}
