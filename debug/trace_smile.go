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
		fmt.Println("Usage: trace_smile <file.bgf>")
		os.Exit(1)
	}

	data, _ := os.ReadFile(os.Args[1])
	parts := bytes.SplitN(data, []byte("\n"), 2)

	var header map[string]interface{}
	json.Unmarshal(parts[0], &header)

	gr, _ := gzip.NewReader(bytes.NewReader(parts[1]))
	decompressed, _ := io.ReadAll(gr)
	gr.Close()

	fmt.Printf("Decompressed size: %d bytes\n\n", len(decompressed))

	// Trace decode
	offset := 4 // Skip :)\n and version
	keys := []string{}
	depth := 0

	for i := 0; i < 100 && offset < len(decompressed); i++ {
		b := decompressed[offset]
		indent := ""
		for j := 0; j < depth; j++ {
			indent += "  "
		}

		fmt.Printf("%s[%04x] 0x%02x: ", indent, offset, b)

		if b == 0xFA {
			fmt.Printf("START_OBJECT\n")
			offset++
			depth++
		} else if b == 0xFB {
			fmt.Printf("END_OBJECT\n")
			offset++
			depth--
		} else if b == 0xF8 {
			fmt.Printf("START_ARRAY\n")
			offset++
			depth++
		} else if b == 0xF9 {
			fmt.Printf("END_ARRAY\n")
			offset++
			depth--
		} else if b == 0x23 {
			fmt.Printf("true\n")
			offset++
		} else if b == 0x22 {
			fmt.Printf("false\n")
			offset++
		} else if b == 0x21 {
			fmt.Printf("null\n")
			offset++
		} else if b >= 0x40 && b < 0xC0 {
			// ASCII string: 0x40-0xBF encode length 0-127
			length := int(b - 0x40)
			offset++
			if offset+length <= len(decompressed) {
				str := string(decompressed[offset : offset+length])
				keys = append(keys, str)
				fmt.Printf("ASCII[%d] '%s' (key#%d)\n", length, str, len(keys)-1)
				offset += length
			}
		} else if b >= 0xC0 && b < 0xE0 {
			value := int(b) - 0xD0
			fmt.Printf("SMALL_INT = %d\n", value)
			offset++
		} else {
			fmt.Printf("UNKNOWN\n")
			offset++
		}
	}
}
