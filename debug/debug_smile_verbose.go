package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: debug_smile_verbose <file.bgf>")
		os.Exit(1)
	}

	// Read file
	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse header
	newlinePos := -1
	for i, b := range raw {
		if b == '\n' {
			newlinePos = i
			break
		}
	}

	if newlinePos < 0 {
		fmt.Println("No newline found in file")
		os.Exit(1)
	}

	var header map[string]interface{}
	if err := json.Unmarshal(raw[:newlinePos], &header); err != nil {
		fmt.Printf("Error parsing header: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Header: %v\n", header)

	// Decompress
	data, err := gzip.NewReader(bytes.NewReader(raw[newlinePos+1:]))
	if err != nil {
		fmt.Printf("Error creating gzip reader: %v\n", err)
		os.Exit(1)
	}
	defer data.Close()

	decompressed, err := io.ReadAll(data)
	if err != nil {
		fmt.Printf("Error decompressing: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nDecompressed size: %d bytes\n", len(decompressed))
	fmt.Printf("SMILE header: %02x %02x %02x %02x\n",
		decompressed[0], decompressed[1], decompressed[2], decompressed[3])

	// Manual decode to see what's happening
	offset := 4 // Skip SMILE header
	if decompressed[offset] != 0xFA {
		fmt.Printf("Expected START_OBJECT at offset %d, got 0x%02x\n", offset, decompressed[offset])
		os.Exit(1)
	}

	offset++ // Skip START_OBJECT
	keyCount := 0
	depth := 1

	fmt.Printf("\nDecoding top-level object keys:\n")

	for offset < len(decompressed) && keyCount < 30 {
		b := decompressed[offset]

		if b == 0xFB { // END_OBJECT
			depth--
			if depth == 0 {
				fmt.Printf("\n*** Top-level object ended at offset %d ***\n", offset)
				break
			}
			offset++
			continue
		}

		if b == 0xFA { // START_OBJECT
			depth++
			offset++
			continue
		}

		// If we're at depth 1, this should be a key
		if depth == 1 {
			keyStr, newOffset := readKey(decompressed, offset)
			if keyStr != "" {
				keyCount++
				fmt.Printf("%3d. Offset %04x: key='%s'\n", keyCount, offset, keyStr)

				// Skip the value
				offset = newOffset
				_, offset = skipValue(decompressed, offset)
			} else {
				offset++
			}
		} else {
			offset++
		}
	}
}

func readKey(data []byte, offset int) (string, int) {
	if offset >= len(data) {
		return "", offset
	}

	b := data[offset]

	// SHORT ASCII key
	if b >= 0x40 && b < 0x80 {
		length := int(b - 0x40)
		if length == 0 {
			return "", offset + 1
		}
		if offset+1+length > len(data) {
			return "", offset
		}
		return string(data[offset+1 : offset+1+length]), offset + 1 + length
	} else if b >= 0x80 && b < 0xC0 {
		length := int(b - 0x80 + 1)
		if offset+1+length > len(data) {
			return "", offset
		}
		return string(data[offset+1 : offset+1+length]), offset + 1 + length
	}

	return "", offset + 1
}

func skipValue(data []byte, offset int) (string, int) {
	if offset >= len(data) {
		return "EOF", offset
	}

	b := data[offset]

	// START_OBJECT
	if b == 0xFA {
		depth := 1
		offset++
		for offset < len(data) && depth > 0 {
			if data[offset] == 0xFA {
				depth++
			} else if data[offset] == 0xFB {
				depth--
			}
			offset++
		}
		return "object", offset
	}

	// START_ARRAY
	if b == 0xF8 {
		depth := 1
		offset++
		for offset < len(data) && depth > 0 {
			if data[offset] == 0xF8 {
				depth++
			} else if data[offset] == 0xF9 {
				depth--
			}
			offset++
		}
		return "array", offset
	}

	// Boolean/null
	if b >= 0x21 && b <= 0x23 {
		return "bool/null", offset + 1
	}

	// Small int
	if b >= 0xC0 && b < 0xE0 {
		return "smallint", offset + 1
	}

	// String
	if b >= 0x20 && b < 0xC0 {
		var length int
		if b >= 0x20 && b < 0x40 {
			length = int(b - 0x20 + 1)
		} else if b >= 0x40 && b < 0x80 {
			length = int(b - 0x40 + 1)
		} else if b >= 0x80 && b < 0xC0 {
			length = int(b - 0x80 + 1)
		}
		return "string", offset + 1 + length
	}

	// Unknown - skip 1 byte
	return fmt.Sprintf("unknown:0x%02x", b), offset + 1
}
