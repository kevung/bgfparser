package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run examine_offset.go <bgf_file> [offset]")
		return
	}

	filename := os.Args[1]
	targetOffset := 1142
	if len(os.Args) >= 3 {
		fmt.Sscanf(os.Args[2], "%d", &targetOffset)
	}

	// Read file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Read header line
	headerBuf := make([]byte, 1024)
	n, _ := file.Read(headerBuf)

	// Find first newline
	headerEnd := 0
	for i := 0; i < n; i++ {
		if headerBuf[i] == '\n' {
			headerEnd = i + 1
			break
		}
	}

	fmt.Printf("Header: %s\n", string(headerBuf[:headerEnd]))

	// Seek back to after header
	file.Seek(int64(headerEnd), 0)

	// Read rest with gzip
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		fmt.Printf("Error creating gzip reader: %v\n", err)
		return
	}
	defer gzReader.Close()

	data, err := io.ReadAll(gzReader)
	if err != nil {
		fmt.Printf("Error reading gzip data: %v\n", err)
		return
	}

	fmt.Printf("Decompressed data size: %d bytes\n\n", len(data))

	// Show context around the target offset
	start := targetOffset - 50
	if start < 0 {
		start = 0
	}
	end := targetOffset + 50
	if end > len(data) {
		end = len(data)
	}

	fmt.Printf("Context around offset %d:\n", targetOffset)
	fmt.Printf("Offset | Hex  | Dec | Char | Type\n")
	fmt.Printf("-------|------|-----|------|-----\n")

	for i := start; i < end; i++ {
		b := data[i]
		marker := " "
		if i == targetOffset {
			marker = ">>>"
		}

		char := "."
		if b >= 32 && b <= 126 {
			char = string(b)
		}

		typeStr := describeSmileByte(b)

		fmt.Printf("%s%4d | 0x%02X | %3d | %-4s | %s\n",
			marker, i, b, b, char, typeStr)
	}

	// Try to trace what keys we're in
	fmt.Printf("\n\nKey context analysis:\n")
	analyzeKeyContext(data, targetOffset)
}

func describeSmileByte(b byte) string {
	switch {
	case b == 0xFA:
		return "START_OBJECT"
	case b == 0xFB:
		return "END_OBJECT"
	case b == 0xF8:
		return "START_ARRAY"
	case b == 0xF9:
		return "END_ARRAY"
	case b == 0x23:
		return "TRUE"
	case b == 0x22:
		return "FALSE"
	case b == 0x21:
		return "NULL"
	case b >= 0x20 && b < 0x40:
		return fmt.Sprintf("TINY_ASCII (len=%d)", b-0x20)
	case b >= 0x40 && b < 0x80:
		return fmt.Sprintf("SHORT_ASCII (len=%d)", b-0x40)
	case b >= 0x80 && b < 0xC0:
		return fmt.Sprintf("SHORT_ASCII_SHARED (len=%d)", b-0x80+1)
	case b >= 0xC0 && b < 0xE0:
		return fmt.Sprintf("SMALL_INT (val=%d)", int(b)-0xD0)
	case b >= 0xE0 && b < 0xE4:
		return "LONG_ASCII"
	case b >= 0xE4 && b < 0xE8:
		return "LONG_UNICODE"
	case b == 0xE8:
		return "BIG_INTEGER"
	case b == 0xE9:
		return "FLOAT32"
	case b == 0xEA:
		return "FLOAT64"
	case b == 0xEB:
		return "BIG_DECIMAL"
	case b < 0x20:
		return fmt.Sprintf("SHARED_REF (#%d)", b)
	default:
		return "UNKNOWN"
	}
}

func analyzeKeyContext(data []byte, targetOffset int) {
	// Walk backward from target to find what object/array context we're in
	depth := 0
	arrayDepth := 0
	objectDepth := 0

	fmt.Printf("Walking backward from offset %d:\n", targetOffset)

	for i := targetOffset; i >= 0 && i > targetOffset-200; i-- {
		b := data[i]

		switch b {
		case 0xFB: // END_OBJECT
			objectDepth++
			depth++
		case 0xFA: // START_OBJECT
			objectDepth--
			depth--
			if objectDepth < 0 {
				fmt.Printf("  Offset %d: Found START_OBJECT (we're inside this object)\n", i)
				// Try to find the key for this object
				if i > 0 {
					findPrecedingKey(data, i)
				}
				return
			}
		case 0xF9: // END_ARRAY
			arrayDepth++
			depth++
		case 0xF8: // START_ARRAY
			arrayDepth--
			depth--
			if arrayDepth < 0 {
				fmt.Printf("  Offset %d: Found START_ARRAY (we're inside this array)\n", i)
				// Try to find the key for this array
				if i > 0 {
					findPrecedingKey(data, i)
				}
				return
			}
		}

		if depth < 0 {
			fmt.Printf("  Offset %d: Found containing structure\n", i)
			break
		}
	}
}

func findPrecedingKey(data []byte, valueOffset int) {
	// Walk backward from value to find the key
	for i := valueOffset - 1; i >= 0 && i > valueOffset-50; i-- {
		b := data[i]

		// Check if this looks like a key (string)
		if b >= 0x20 && b < 0x40 {
			// Tiny ASCII
			length := int(b - 0x20)
			if i+1+length <= valueOffset {
				key := string(data[i+1 : i+1+length])
				fmt.Printf("    Preceding key might be: '%s' (tiny ASCII at offset %d)\n", key, i)
				return
			}
		} else if b >= 0x40 && b < 0xC0 {
			// Short ASCII
			var length int
			if b >= 0x40 && b < 0x80 {
				length = int(b - 0x40)
			} else {
				length = int(b - 0x80 + 1)
			}
			if i+1+length <= valueOffset {
				key := string(data[i+1 : i+1+length])
				fmt.Printf("    Preceding key might be: '%s' (short ASCII at offset %d)\n", key, i)
				return
			}
		} else if b < 0x20 {
			// Shared reference
			fmt.Printf("    Preceding key is shared reference #%d at offset %d\n", b, i)
			return
		}
	}
}
