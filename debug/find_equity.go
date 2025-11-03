package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run find_equity.go <bgf_file>")
		return
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer file.Close()

	// Read header
	headerBuf := make([]byte, 1024)
	n, _ := file.Read(headerBuf)
	headerEnd := 0
	for i := 0; i < n; i++ {
		if headerBuf[i] == '\n' {
			headerEnd = i + 1
			break
		}
	}

	file.Seek(int64(headerEnd), 0)
	gzReader, _ := gzip.NewReader(file)
	data, _ := io.ReadAll(gzReader)

	// Find "equity"
	idx := strings.Index(string(data), "equity")
	if idx >= 0 {
		fmt.Printf("Found \"equity\" at offset %d\n\n", idx)
		start := idx - 30
		if start < 0 {
			start = 0
		}
		end := idx + 50
		if end > len(data) {
			end = len(data)
		}

		fmt.Println("Context:")
		for i := start; i < end; i++ {
			marker := "  "
			if i == idx {
				marker = ">>"
			}
			b := data[i]
			char := "."
			if b >= 32 && b <= 126 {
				char = string(b)
			}
			fmt.Printf("%s%4d: 0x%02X (%3d) %s  %s\n", marker, i, b, b, char, describeSmileByte(b))
		}
	} else {
		fmt.Println("\"equity\" not found in data")
	}
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
		return fmt.Sprintf("TINY_ASCII(len=%d)", b-0x20)
	case b >= 0x40 && b < 0x80:
		return fmt.Sprintf("SHORT_ASCII(len=%d)", b-0x40)
	case b >= 0x80 && b < 0xC0:
		return fmt.Sprintf("SHORT_ASCII_SHARED(len=%d)", b-0x80+1)
	case b >= 0xC0 && b < 0xE0:
		return fmt.Sprintf("SMALL_INT(%d)", int(b)-0xD0)
	default:
		return ""
	}
}
