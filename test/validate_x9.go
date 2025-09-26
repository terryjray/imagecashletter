package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/moov-io/imagecashletter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run validate_x9.go <filename>")
		fmt.Println("Example: go run validate_x9.go /path/to/your/file.x937")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", filename)
	}

	fmt.Printf("🔍 Validating X9 file: %s\n", filename)
	fmt.Println(strings.Repeat("=", 60))

	// Step 1: Detect encoding
	encoding, err := detectX9Encoding(filename)
	if err != nil {
		fmt.Printf("❌ Could not detect encoding: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📝 Detected encoding: %s\n", strings.ToUpper(encoding))

	// Step 2: Open file and create reader with appropriate options
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	// Configure reader options based on encoding
	readerOptions := []imagecashletter.ReaderOption{
		imagecashletter.BufferSizeOption(65536), // 64k buffer like sickle
		imagecashletter.ReadVariableLineLengthOption(),
	}

	if encoding == "ebcdic" {
		readerOptions = append(readerOptions, imagecashletter.ReadEbcdicEncodingOption())
		fmt.Println("🔧 EBCDIC encoding option enabled")
	}

	reader := imagecashletter.NewReader(file, readerOptions...)

	// Step 3: Read and validate
	fmt.Println("\n📖 Reading file...")
	iclFile, err := reader.Read()
	if err != nil {
		fmt.Printf("❌ Error reading file: %v\n", err)
		
		// Provide specific guidance based on error type
		if strings.Contains(err.Error(), "Bundle control without a current bundle") {
			fmt.Println("\n💡 This file has structural issues and should be moved to hold directory")
		} else if strings.Contains(err.Error(), "ArchiveTypeIndicator 0 is invalid") {
			fmt.Println("\n⚠️  Warning: Invalid ArchiveTypeIndicator (non-fatal)")
		} else if strings.Contains(err.Error(), "TestFileIndicator") {
			fmt.Println("\n💡 Missing TestFileIndicator - this is a common vendor variation")
		} else if strings.Contains(err.Error(), "unknown record type") {
			fmt.Println("\n💡 Unknown record type - check for null padding or encoding issues")
		}
		
		os.Exit(1)
	}

	// Step 4: Validate file structure
	fmt.Println("✅ File read successfully!")
	fmt.Println("\n🔍 Validating file structure...")
	if err = iclFile.Validate(); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
		os.Exit(1)
	}

	// Step 5: Display file summary
	fmt.Println("✅ File validation passed!")
	fmt.Println("\n📊 File Summary:")
	fmt.Printf("   • Cash Letters: %d\n", len(iclFile.CashLetters))
	fmt.Printf("   • Total Records: %d\n", iclFile.Control.TotalRecordCount)
	fmt.Printf("   • Total Items: %d\n", iclFile.Control.TotalItemCount)
	fmt.Printf("   • Total Amount: $%.2f\n", float64(iclFile.Control.FileTotalAmount)/100)

	if len(iclFile.CashLetters) > 0 {
		cl := iclFile.CashLetters[0]
		fmt.Println("\n📋 First Cash Letter Details:")
		fmt.Printf("   • ID: '%s'\n", cl.CashLetterHeader.CashLetterID)
		fmt.Printf("   • Business Date: %s\n", cl.CashLetterHeader.CashLetterBusinessDate.Format("2006-01-02"))
		fmt.Printf("   • Destination: %s\n", cl.CashLetterHeader.DestinationRoutingNumber)
		fmt.Printf("   • Originator: %s\n", cl.CashLetterHeader.OriginatorContactName)
	}

	fmt.Println("\n🎉 File is valid and ready for processing!")
}

// detectX9Encoding detects if the file is ASCII or EBCDIC encoded
// Based on the logic from sickle parser, with support for null padding
func detectX9Encoding(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Read enough bytes to check multiple positions
	header := make([]byte, 10)
	_, err = f.Read(header)
	if err != nil {
		return "", err
	}

	// Check for null padding first (bytes 0-1 should be 0x00 0x00)
	if header[0] == 0x00 && header[1] == 0x00 {
		// File has null padding, check bytes 4-5 for record type
		byte4, byte5 := header[4], header[5]
		switch {
		case byte4 == 0x30 && byte5 == 0x31: // ASCII '0' '1'
			return "ascii", nil
		case byte4 == 0xF0 && byte5 == 0xF1: // EBCDIC '0' '1'
			return "ebcdic", nil
		default:
			return "", fmt.Errorf("unrecognized encoding at bytes 5-6 (with null padding): %#x %#x", byte4, byte5)
		}
	} else {
		// No null padding, check bytes 4-5 for record type
		byte4, byte5 := header[4], header[5]
		switch {
		case byte4 == 0x30 && byte5 == 0x31: // ASCII '0' '1'
			return "ascii", nil
		case byte4 == 0xF0 && byte5 == 0xF1: // EBCDIC '0' '1'
			return "ebcdic", nil
		default:
			// Try bytes 0-1 in case file was already preprocessed
			byte0, byte1 := header[0], header[1]
			switch {
			case byte0 == 0x30 && byte1 == 0x31: // ASCII '0' '1'
				return "ascii", nil
			case byte0 == 0xF0 && byte1 == 0xF1: // EBCDIC '0' '1'
				return "ebcdic", nil
			default:
				return "", fmt.Errorf("unrecognized encoding at bytes 1-2: %#x %#x", byte0, byte1)
			}
		}
	}
}
