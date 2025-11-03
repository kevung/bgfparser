package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run trace_1835.go <bgf_file>")
		return
	}

	filename := os.Args[1]

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

	// Manual trace of what should happen at offset 1832-1840
	fmt.Println("Manual trace of offsets 1832-1840:")
	fmt.Println("Expected flow:")
	fmt.Println("  1832: 0x40 -> readKey() -> empty string \"\" -> offset becomes 1833")
	fmt.Println("  1833: 0xF8 -> decode() sees START_ARRAY -> calls readArray() -> offset becomes 1834")
	fmt.Println("  1834: 0xDA -> decode() sees SMALL_INT -> reads value 10 -> offset becomes 1835")
	fmt.Println("  1835: 0xDA -> decode() sees SMALL_INT -> reads value 10 -> offset becomes 1836")
	fmt.Println("  1836: 0xC1 -> decode() sees SMALL_INT -> reads value -15 -> offset becomes 1837")
	fmt.Println("  1837: 0xC1 -> decode() sees SMALL_INT -> reads value -15 -> offset becomes 1838")
	fmt.Println("  1838: 0xF9 -> readArray() sees END_ARRAY -> consumes it -> offset becomes 1839")
	fmt.Println("  1839: 0x30 -> readObject() reads next key")
	fmt.Println()

	// Now let's see what actually happens
	fmt.Println("Actual bytes:")
	for i := 1832; i <= 1840; i++ {
		b := data[i]
		fmt.Printf("  %d: 0x%02X (%d) = %s\n", i, b, b, describeSmileByte(b))
	}
	fmt.Println()

	// Simulate array reading
	fmt.Println("Simulating readArray() from offset 1833:")
	offset := 1833
	if data[offset] != 0xF8 {
		fmt.Printf("ERROR: Expected START_ARRAY at %d, got 0x%02X\n", offset, data[offset])
		return
	}
	offset++ // Skip START_ARRAY, now at 1834

	arrayElements := []interface{}{}
	for offset < len(data) {
		b := data[offset]
		fmt.Printf("  At offset %d: byte 0x%02X (%s)\n", offset, b, describeSmileByte(b))

		if b == 0xF9 { // END_ARRAY
			fmt.Printf("    -> END_ARRAY found, consuming and returning\n")
			offset++
			break
		}

		// Try to decode the value
		if b >= 0xC0 && b < 0xE0 {
			value := int(b) - 0xD0
			fmt.Printf("    -> Decoded SMALL_INT: %d\n", value)
			offset++
			arrayElements = append(arrayElements, value)
		} else {
			fmt.Printf("    -> ERROR: Unexpected type in array\n")
			break
		}
	}

	fmt.Printf("\nArray elements: %v\n", arrayElements)
	fmt.Printf("Final offset after readArray(): %d\n", offset)
	fmt.Printf("Next byte at %d: 0x%02X (%s)\n", offset, data[offset], describeSmileByte(data[offset]))
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
	case b >= 0x20 && b < 0x40:
		return fmt.Sprintf("TINY_ASCII (len=%d)", b-0x20)
	case b >= 0x40 && b < 0x80:
		return fmt.Sprintf("SHORT_ASCII (len=%d)", b-0x40)
	case b >= 0x80 && b < 0xC0:
		return fmt.Sprintf("SHORT_ASCII_SHARED (len=%d)", b-0x80+1)
	case b >= 0xC0 && b < 0xE0:
		return fmt.Sprintf("SMALL_INT (val=%d)", int(b)-0xD0)
	default:
		return "OTHER"
	}
}
