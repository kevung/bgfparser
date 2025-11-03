package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Debug tool to examine SMILE encoded BGF files
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: test_smile <file.bgf>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Split header and data
	parts := bytes.SplitN(data, []byte("\n"), 2)
	if len(parts) != 2 {
		fmt.Println("Invalid BGF format")
		os.Exit(1)
	}

	// Parse header
	var header map[string]interface{}
	if err := json.Unmarshal(parts[0], &header); err != nil {
		fmt.Printf("Error parsing header: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Header: %+v\n\n", header)

	// Decompress
	gr, err := gzip.NewReader(bytes.NewReader(parts[1]))
	if err != nil {
		fmt.Printf("Error creating gzip reader: %v\n", err)
		os.Exit(1)
	}
	defer gr.Close()

	decompressed, err := io.ReadAll(gr)
	if err != nil {
		fmt.Printf("Error decompressing: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Decompressed size: %d bytes\n", len(decompressed))
	fmt.Printf("First 200 bytes (hex):\n")
	for i := 0; i < 200 && i < len(decompressed); i++ {
		if i%16 == 0 {
			fmt.Printf("\n%04x: ", i)
		}
		fmt.Printf("%02x ", decompressed[i])
	}
	fmt.Printf("\n\nFirst 200 bytes (chars):\n")
	for i := 0; i < 200 && i < len(decompressed); i++ {
		if i%64 == 0 {
			fmt.Printf("\n%04x: ", i)
		}
		if decompressed[i] >= 32 && decompressed[i] <= 126 {
			fmt.Printf("%c", decompressed[i])
		} else {
			fmt.Printf(".")
		}
	}
	fmt.Println()

	// Check for SMILE header
	if len(decompressed) >= 4 {
		fmt.Printf("\nFirst 4 bytes: %c%c%c%c (0x%02x 0x%02x 0x%02x 0x%02x)\n",
			printable(decompressed[0]), printable(decompressed[1]),
			printable(decompressed[2]), printable(decompressed[3]),
			decompressed[0], decompressed[1], decompressed[2], decompressed[3])
	}

	// Extract strings
	fmt.Println("\nExtracted strings (min 5 chars):")
	strings := extractStrings(decompressed, 5)
	for i, s := range strings {
		if i >= 50 { // Limit output
			break
		}
		if len(s) > 60 {
			fmt.Printf("%3d: %s...\n", i+1, s[:60])
		} else {
			fmt.Printf("%3d: %s\n", i+1, s)
		}
	}
}

func printable(b byte) byte {
	if b >= 32 && b <= 126 {
		return b
	}
	return '.'
}

func extractStrings(data []byte, minLen int) []string {
	var strings []string
	var current []byte

	for _, b := range data {
		if b >= 32 && b <= 126 {
			current = append(current, b)
		} else {
			if len(current) >= minLen {
				strings = append(strings, string(current))
			}
			current = nil
		}
	}

	if len(current) >= minLen {
		strings = append(strings, string(current))
	}

	return strings
}
