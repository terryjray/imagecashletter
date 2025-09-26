package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run simple_strip.go <input_file> <output_file>")
		fmt.Println("This utility detects and strips null padding from X9 files")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Open input file
	in, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	// Create output file
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	// Read first 4 bytes to detect null padding
	prefix := make([]byte, 4)
	_, err = io.ReadFull(in, prefix)
	if err != nil {
		fmt.Printf("Error reading file prefix: %v\n", err)
		os.Exit(1)
	}

	// Check if this looks like null padding (starts with 0x00 0x00)
	if prefix[0] == 0x00 && prefix[1] == 0x00 {
		fmt.Printf("Detected null padding: %02x %02x %02x %02x\n", prefix[0], prefix[1], prefix[2], prefix[3])
		fmt.Println("Stripping null padding...")
		
		// Copy the rest of the file
		_, err = io.Copy(out, in)
		if err != nil {
			fmt.Printf("Error copying file: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✅ Successfully stripped null padding\n")
	} else {
		fmt.Printf("No null padding detected: %02x %02x %02x %02x\n", prefix[0], prefix[1], prefix[2], prefix[3])
		fmt.Println("Copying file as-is...")
		
		// Write the prefix back and copy the rest
		_, err = out.Write(prefix)
		if err != nil {
			fmt.Printf("Error writing prefix: %v\n", err)
			os.Exit(1)
		}
		
		_, err = io.Copy(out, in)
		if err != nil {
			fmt.Printf("Error copying file: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✅ File copied without modification\n")
	}

	fmt.Printf("Output written to: %s\n", outputFile)
}
