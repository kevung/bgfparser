package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kevung/bgfparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: trace_decode <bgf_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	match, err := bgfparser.ParseBGF(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	if match != nil {
		// Pretty print the decoded data
		printData("", match.Data, 0, 5)

		fmt.Printf("\n=== Decoding Summary ===\n")
		if match.DecodingWarning != "" {
			fmt.Printf("Warning: %s\n", match.DecodingWarning)
		}
		if offset, ok := match.Data["_finalOffset"].(int); ok {
			fmt.Printf("Final offset: %d\n", offset)
		}
		if total, ok := match.Data["_totalBytes"].(int); ok {
			fmt.Printf("Total bytes: %d\n", total)
			if offset, ok := match.Data["_finalOffset"].(int); ok {
				pct := float64(offset) * 100.0 / float64(total)
				fmt.Printf("Decoded: %.2f%%\n", pct)
			}
		}
	}
}

func printData(key string, value interface{}, depth int, maxDepth int) {
	indent := strings.Repeat("  ", depth)

	if depth >= maxDepth {
		fmt.Printf("%s%s: <max depth reached>\n", indent, key)
		return
	}

	switch v := value.(type) {
	case map[string]interface{}:
		if key != "" {
			fmt.Printf("%s%s:\n", indent, key)
		}

		// Print metadata fields first
		metaKeys := []string{"_keyCount", "_errorOffset", "_lastError", "_finalOffset", "_totalBytes"}
		for _, mk := range metaKeys {
			if val, ok := v[mk]; ok {
				printData(mk, val, depth+1, maxDepth)
			}
		}

		// Print regular fields
		for k, val := range v {
			// Skip metadata we already printed
			isMeta := false
			for _, mk := range metaKeys {
				if k == mk {
					isMeta = true
					break
				}
			}
			if !isMeta {
				printData(k, val, depth+1, maxDepth)
			}
		}

	case []interface{}:
		if key != "" {
			fmt.Printf("%s%s: [\n", indent, key)
		}
		for i, item := range v {
			if i > 10 {
				fmt.Printf("%s  ... (%d more items)\n", indent, len(v)-i)
				break
			}
			printData(fmt.Sprintf("[%d]", i), item, depth+1, maxDepth)
		}
		if key != "" {
			fmt.Printf("%s]\n", indent)
		}

	case string:
		// Truncate long strings
		if len(v) > 100 {
			fmt.Printf("%s%s: %q...\n", indent, key, v[:100])
		} else {
			fmt.Printf("%s%s: %q\n", indent, key, v)
		}

	case float64, float32, int, int32, int64, bool:
		fmt.Printf("%s%s: %v\n", indent, key, v)

	case nil:
		fmt.Printf("%s%s: null\n", indent, key)

	default:
		fmt.Printf("%s%s: %v (type: %T)\n", indent, key, v, v)
	}
}
